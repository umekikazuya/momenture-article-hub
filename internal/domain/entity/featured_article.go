package entity

import "time"

// FeaturedArticle represents a featured article.
type FeaturedArticle struct {
	ID        int64
	ArticleID int64
	SortOrder int
	CreatedAt time.Time
}

// NewFeaturedArticle creates a new FeaturedArticle instance.
func NewFeaturedArticle(articleID int64, sortOrder int) *FeaturedArticle {
	return &FeaturedArticle{
		ArticleID: articleID,
		SortOrder: sortOrder,
		CreatedAt: time.Now(),
	}
}
