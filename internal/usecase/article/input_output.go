package article

import "time"

// GetArticlesInput is the input for retrieving articles.
type FindByCriteriaInput struct {
	Status       *string `json:"status" validate:"omitempty,oneof=draft published"`
	ProviderType *string `json:"provider_type" validate:"omitempty"`
	SortBy       *string `json:"sort_by" validate:"omitempty,oneof=created_at updated_at title"`
	SortOrder    *string `json:"sort_order" validate:"omitempty,oneof=asc desc"`
	Page         int     `json:"page" validate:"gte=1"`
	Limit        int     `json:"limit" validate:"gte=1,lte=100"`
}

// FindByCriteriaOutput is the output for retrieving articles by criteria.
type FindByCriteriaOutput struct {
	Articles   []FindArticleByIDOutput `json:"articles"`
	Total      int64                   `json:"total"`
	Page       int                     `json:"page"`
	Limit      int                     `json:"limit"`
	TotalPages int                     `json:"total_pages"`
}

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

// FindArticleByIDOutput is the output for finding an article by ID.
type FindArticleByIDOutput struct {
	ID           uint64    `json:"id"`
	Title        string    `json:"title"`
	Body         string    `json:"body"`
	Status       string    `json:"status"`
	ProviderType string    `json:"provider_type"`
	Link         string    `json:"link"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// UpdateArticleInput is the input for updating an article.
type UpdateArticleInput struct {
	Title        *string `json:"title,omitempty"`
	Body         *string `json:"body,omitempty"`
	Status       *string `json:"status,omitempty"`
	ProviderType *string `json:"provider_type,omitempty"`
	Link         *string `json:"link,omitempty"`
}

// UpdateArticleOutput is the output for updating an article.
type UpdateArticleOutput struct {
	ID           uint64    `json:"id"`
	Title        string    `json:"title"`
	Body         string    `json:"body"`
	Status       string    `json:"status"`
	ProviderType string    `json:"provider_type"`
	Link         string    `json:"link"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
