package postgres

import (
	"context"
	"errors"

	"github.com/umekikazuya/momenture-article-hub/internal/domain/entity"
	"github.com/umekikazuya/momenture-article-hub/internal/domain/repository"
)

type PostgresArticleRepository struct {
}

func NewPostgresArticleRepository() repository.ArticleRepository {
	return &PostgresArticleRepository{}
}

// FindAll は全ての記事を取得する
func (r *PostgresArticleRepository) FindAll(ctx context.Context) ([]*entity.Article, error) {
	return nil, errors.New("not implemented")
}

// FindByID は指定されたIDの記事を取得する
func (r *PostgresArticleRepository) FindByID(ctx context.Context, id uint64) (*entity.Article, error) {
	return nil, errors.New("not implemented")
}

// FindByCriteria は指定された条件に一致する記事を取得する
func (r *PostgresArticleRepository) FindByCriteria(ctx context.Context, criteria repository.ArticleQueryCriteria) ([]*entity.Article, int, error) {
	return nil, 0, errors.New("not implemented")
}

// Create は新しい記事を作成する
func (r *PostgresArticleRepository) Create(ctx context.Context, article *entity.Article) (*entity.Article, error) {
	return nil, errors.New("not implemented")
}

// Update は既存の記事を更新する
func (r *PostgresArticleRepository) Update(ctx context.Context, article *entity.Article) error {
	return errors.New("not implemented")
}

// Delete は指定されたIDの記事を論理削除する
func (r *PostgresArticleRepository) Delete(ctx context.Context, id uint64) error {
	return errors.New("not implemented")
}
