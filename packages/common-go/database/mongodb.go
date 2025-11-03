package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoDB wraps MongoDB client
type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

// Connect connects to MongoDB
func Connect(uri, databaseName string, maxPoolSize, minPoolSize uint64, timeout time.Duration) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	opts := options.Client().
		ApplyURI(uri).
		SetMaxPoolSize(maxPoolSize).
		SetMinPoolSize(minPoolSize)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	// Ping to verify connection
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	return &MongoDB{
		Client:   client,
		Database: client.Database(databaseName),
	}, nil
}

// Disconnect disconnects from MongoDB
func (db *MongoDB) Disconnect(ctx context.Context) error {
	return db.Client.Disconnect(ctx)
}

// Collection returns a collection
func (db *MongoDB) Collection(name string) *mongo.Collection {
	return db.Database.Collection(name)
}

// Transaction executes a transaction
func (db *MongoDB) Transaction(ctx context.Context, fn func(sessCtx mongo.SessionContext) error) error {
	session, err := db.Client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	return mongo.WithSession(ctx, session, func(sessCtx mongo.SessionContext) error {
		if err := session.StartTransaction(); err != nil {
			return err
		}

		if err := fn(sessCtx); err != nil {
			session.AbortTransaction(sessCtx)
			return err
		}

		return session.CommitTransaction(sessCtx)
	})
}

