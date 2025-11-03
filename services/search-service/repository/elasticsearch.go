package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

const (
	contentIndex = "content"
)

// ElasticsearchClient wraps Elasticsearch client
type ElasticsearchClient struct {
	client *elasticsearch.Client
}

// NewElasticsearchClient creates a new Elasticsearch client
func NewElasticsearchClient() (*ElasticsearchClient, error) {
	addresses := os.Getenv("ELASTICSEARCH_ADDRESSES")
	if addresses == "" {
		addresses = "http://localhost:9200"
	}

	cfg := elasticsearch.Config{
		Addresses: strings.Split(addresses, ","),
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	es := &ElasticsearchClient{client: client}
	
	// Create index if not exists
	if err := es.CreateIndex(); err != nil {
		return nil, fmt.Errorf("failed to create index: %w", err)
	}

	return es, nil
}

// CreateIndex creates the content index with mappings
func (c *ElasticsearchClient) CreateIndex() error {
	mapping := map[string]interface{}{
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"id":          {"type": "keyword"},
				"title":       {"type": "text", "analyzer": "standard"},
				"description": {"type": "text", "analyzer": "standard"},
				"genre":       {"type": "keyword"},
				"category":    {"type": "keyword"},
				"cast":        {"type": "keyword"},
				"directors":   {"type": "keyword"},
				"tags":        {"type": "keyword"},
				"release_year": {"type": "integer"},
				"rating":      {"type": "float"},
			},
		},
	}

	body, _ := json.Marshal(mapping)
	req := esapi.IndicesCreateRequest{
		Index: contentIndex,
		Body:  bytes.NewReader(body),
	}

	res, err := req.Do(context.Background(), c.client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() && res.StatusCode != 400 { // 400 = index already exists
		return fmt.Errorf("error creating index: %s", res.String())
	}

	return nil
}

// IndexContent indexes a content document
func (c *ElasticsearchClient) IndexContent(ctx context.Context, id string, doc interface{}) error {
	body, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	req := esapi.IndexRequest{
		Index:      contentIndex,
		DocumentID: id,
		Body:       bytes.NewReader(body),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, c.client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error indexing document: %s", res.String())
	}

	return nil
}

// Search searches for content
func (c *ElasticsearchClient) Search(ctx context.Context, query string, filters map[string]interface{}, page, pageSize int) ([]map[string]interface{}, int64, error) {
	must := []map[string]interface{}{
		{
			"multi_match": map[string]interface{}{
				"query":  query,
				"fields": []string{"title^2", "description", "cast", "directors"},
			},
		},
	}

	// Add filters
	if filters != nil {
		for key, value := range filters {
			must = append(must, map[string]interface{}{
				"term": map[string]interface{}{key: value},
			})
		}
	}

	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": must,
			},
		},
		"from": (page - 1) * pageSize,
		"size": pageSize,
	}

	body, _ := json.Marshal(searchQuery)
	req := esapi.SearchRequest{
		Index: []string{contentIndex},
		Body:  bytes.NewReader(body),
	}

	res, err := req.Do(ctx, c.client)
	if err != nil {
		return nil, 0, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, 0, fmt.Errorf("search error: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, 0, err
	}

	hits := result["hits"].(map[string]interface{})
	total := int64(hits["total"].(map[string]interface{})["value"].(float64))
	
	var documents []map[string]interface{}
	for _, hit := range hits["hits"].([]interface{}) {
		doc := hit.(map[string]interface{})["_source"].(map[string]interface{})
		documents = append(documents, doc)
	}

	return documents, total, nil
}

// Autocomplete provides autocomplete suggestions
func (c *ElasticsearchClient) Autocomplete(ctx context.Context, query string) ([]string, error) {
	searchQuery := map[string]interface{}{
		"suggest": map[string]interface{}{
			"title_suggest": map[string]interface{}{
				"prefix": query,
				"completion": map[string]interface{}{
					"field": "title_suggest",
				},
			},
		},
	}

	body, _ := json.Marshal(searchQuery)
	req := esapi.SearchRequest{
		Index: []string{contentIndex},
		Body:  bytes.NewReader(body),
	}

	res, err := req.Do(ctx, c.client)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("autocomplete error: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	suggestions := []string{}
	if suggest, ok := result["suggest"].(map[string]interface{}); ok {
		if titleSuggest, ok := suggest["title_suggest"].([]interface{}); ok {
			for _, item := range titleSuggest {
				if options, ok := item.(map[string]interface{})["options"].([]interface{}); ok {
					for _, option := range options {
						if text, ok := option.(map[string]interface{})["text"].(string); ok {
							suggestions = append(suggestions, text)
						}
					}
				}
			}
		}
	}

	return suggestions, nil
}

// DeleteContent removes content from index
func (c *ElasticsearchClient) DeleteContent(ctx context.Context, id string) error {
	req := esapi.DeleteRequest{
		Index:      contentIndex,
		DocumentID: id,
	}

	res, err := req.Do(ctx, c.client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() && res.StatusCode != 404 {
		return fmt.Errorf("error deleting document: %s", res.String())
	}

	return nil
}

// GetClient returns the underlying Elasticsearch client
func (c *ElasticsearchClient) GetClient() *elasticsearch.Client {
	return c.client
}

