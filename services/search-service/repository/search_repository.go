package repository

// SearchRepository wraps Elasticsearch operations
type SearchRepository struct {
	client *ElasticsearchClient
}

// NewSearchRepository creates a new search repository
func NewSearchRepository(client *ElasticsearchClient) *SearchRepository {
	return &SearchRepository{
		client: client,
	}
}

// IndexContent indexes content
func (r *SearchRepository) IndexContent(ctx context.Context, id string, doc interface{}) error {
	return r.client.IndexContent(ctx, id, doc)
}

// Search searches for content
func (r *SearchRepository) Search(ctx context.Context, query string, filters map[string]interface{}, page, pageSize int) ([]map[string]interface{}, int64, error) {
	return r.client.Search(ctx, query, filters, page, pageSize)
}

// Autocomplete provides autocomplete suggestions
func (r *SearchRepository) Autocomplete(ctx context.Context, query string) ([]string, error) {
	return r.client.Autocomplete(ctx, query)
}

// DeleteContent deletes content from index
func (r *SearchRepository) DeleteContent(ctx context.Context, id string) error {
	return r.client.DeleteContent(ctx, id)
}

