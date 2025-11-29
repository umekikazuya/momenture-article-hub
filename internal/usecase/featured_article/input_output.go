package featured_article

import "time"

// AddFeaturedArticleInput represents the input for adding a featured article.
type AddFeaturedArticleInput struct {
	ArticleID int64 `json:"article_id"`
	SortOrder int   `json:"sort_order"`
}

// FeaturedArticleOutput represents a single featured article for output.
type FeaturedArticleOutput struct {
	ID        int64     `json:"id"`
	ArticleID int64     `json:"article_id"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
}
