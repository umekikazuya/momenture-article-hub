package article

import "time"

// CreateArticleInput is the input for creating an article.
type CreateArticleInput struct {
	Title        string  `json:"title"`
	Body         *string `json:"body,omitempty"`
	Status       string  `json:"status,omitempty"`
	ProviderType *string `json:"provider_type,omitempty"`
	Link         *string `json:"link,omitempty"`
}

// CreateArticleOutput is the output for creating an article.
type CreateArticleOutput struct {
	ID           uint64    `json:"id"`
	Title        string    `json:"title"`
	Body         string    `json:"body"`
	Status       string    `json:"status"`
	ProviderType string    `json:"provider_type"`
	Link         string    `json:"link"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
