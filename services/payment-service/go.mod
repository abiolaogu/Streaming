module github.com/streamverse/payment-service

go 1.21

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/streamverse/common-go v0.0.0
	go.mongodb.org/mongo-driver v1.13.1
	go.uber.org/zap v1.26.0
)

replace github.com/streamverse/common-go => ../../packages/common-go
