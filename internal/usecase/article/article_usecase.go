package article

import (
	"context"
	"fmt"

	"github.com/umekikazuya/momenture-article-hub/internal/domain/entity"
	"github.com/umekikazuya/momenture-article-hub/internal/domain/repository"
)

// ArticleUsecase defines the interface for article use cases.
type ArticleUsecase struct {
	repo repository.ArticleRepository
}

// NewArticleUsecase creates a new ArticleUsecase.
func NewArticleUsecase(repo repository.ArticleRepository) *ArticleUsecase {
	return &ArticleUsecase{repo: repo}
}

// CreateArticle creates a new article.
func (uc *ArticleUsecase) CreateArticle(ctx context.Context, input CreateArticleInput) (*CreateArticleOutput, error) {
	articleEntity, err := entity.NewArticle(
		input.Title,
		input.Status,
		entity.WithBody(input.Body),
		entity.WithLink(input.Link),
		entity.WithProviderType(input.ProviderType),
	)
	if err != nil {
		return nil, err
	}

	id, err := uc.repo.Create(ctx, articleEntity)
	if err != nil {
		return nil, err
	}

	return &CreateArticleOutput{
		ID:           id,
		Title:        articleEntity.Title.String(),
		Body:         articleEntity.Body.String(),
		Status:       articleEntity.Status.String(),
		ProviderType: articleEntity.ProviderType.String(),
		Link:         articleEntity.Link.String(),
		CreatedAt:    articleEntity.CreatedAt,
		UpdatedAt:    articleEntity.UpdatedAt,
	}, nil
}

// FindArticleByID retrieves an article by its ID.
func (uc *ArticleUsecase) FindArticleByID(ctx context.Context, id uint64) (*FindArticleByIDOutput, error) {
	article, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &FindArticleByIDOutput{
		ID:           article.ID,
		Title:        article.Title.String(),
		Body:         article.Body.String(),
		Status:       article.Status.String(),
		ProviderType: article.ProviderType.String(),
		Link:         article.Link.String(),
		CreatedAt:    article.CreatedAt,
		UpdatedAt:    article.UpdatedAt,
	}, nil
}

// UpdateArticle updates an existing article.
func (uc *ArticleUsecase) UpdateArticle(ctx context.Context, id uint64, input UpdateArticleInput) (*UpdateArticleOutput, error) {
	article, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	err = article.Update(
		input.Title,
		input.Body,
		input.Status,
		input.ProviderType,
		input.Link,
	)

	if err != nil {
		return nil, err
	}

	err = uc.repo.Update(ctx, article)
	if err != nil {
		return nil, err
	}

	return &UpdateArticleOutput{
		ID:           article.ID,
		Title:        article.Title.String(),
		Body:         article.Body.String(),
		Status:       article.Status.String(),
		ProviderType: article.ProviderType.String(),
		Link:         article.Link.String(),
		CreatedAt:    article.CreatedAt,
		UpdatedAt:    article.UpdatedAt,
	}, nil
}

// DeleteArticle deletes an article by its ID.
func (uc *ArticleUsecase) DeleteArticle(ctx context.Context, id uint64) error {
	entity, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if entity == nil {
		return fmt.Errorf("a: %d", id)
	}
	return uc.repo.Delete(ctx, entity.ID)
}
