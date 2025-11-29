package repository

import (
	"context"

	"github.com/umekikazuya/momenture-article-hub/internal/domain/entity"
)

// FeaturedArticleRepository は注目記事の永続化を担うリポジトリインターフェース
type FeaturedArticleRepository interface {
	FindAll(ctx context.Context) ([]*entity.FeaturedArticle, error)
	Create(ctx context.Context, featuredArticle *entity.FeaturedArticle) error
	Delete(ctx context.Context, articleID int64) error
}
