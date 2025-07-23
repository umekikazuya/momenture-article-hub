package repository

import (
	"context"

	"github.com/umekikazuya/momenture-article-hub/internal/domain/entity"
)

// ArticleRepository は記事の永続化を担うリポジトリインターフェース
type ArticleRepository interface {
	FindAll(ctx context.Context) ([]*entity.Article, error)
	FindByID(ctx context.Context, id uint64) (*entity.Article, error)
	FindByCriteria(ctx context.Context, criteria ArticleQueryCriteria) ([]*entity.Article, int, error)
	Create(ctx context.Context, article *entity.Article) (*entity.Article, error)
	Update(ctx context.Context, article *entity.Article) error
	Delete(ctx context.Context, id uint64) error
}

// ArticleQueryCriteria は記事検索の条件を表す
type ArticleQueryCriteria struct {
	Status         *string
	ProviderType   *string
	SortBy         *string
	SortOrder      *string
	Page           int
	Limit          int
	IncludeDeleted bool
}
