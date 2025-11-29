package featured_article

import (
	"context"

	"momenture-article-hub/internal/domain/repository"
)

// FeaturedArticleUsecase defines the use case for featured articles.
type FeaturedArticleUsecase interface {
	GetFeaturedArticles(ctx context.Context) ([]*FeaturedArticleOutput, error)
	AddFeaturedArticle(ctx context.Context, input AddFeaturedArticleInput) error
	RemoveFeaturedArticle(ctx context.Context, articleID int64) error
}

type featuredArticleUsecase struct {
	featuredArticleRepo repository.FeaturedArticleRepository
	articleRepo         repository.ArticleRepository
}

// NewFeaturedArticleUsecase creates a new FeaturedArticleUsecase.
func NewFeaturedArticleUsecase(featuredArticleRepo repository.FeaturedArticleRepository, articleRepo repository.ArticleRepository) FeaturedArticleUsecase {
	return &featuredArticleUsecase{
		featuredArticleRepo: featuredArticleRepo,
		articleRepo:         articleRepo,
	}
}

func (uc *featuredArticleUsecase) GetFeaturedArticles(ctx context.Context) ([]*FeaturedArticleOutput, error) {
	// Implementation goes here
	return nil, nil
}

func (uc *featuredArticleUsecase) AddFeaturedArticle(ctx context.Context, input AddFeaturedArticleInput) error {
	// Implementation goes here
	return nil
}

func (uc *featuredArticleUsecase) RemoveFeaturedArticle(ctx context.Context, articleID int64) error {
	// Implementation goes here
	return nil
}
