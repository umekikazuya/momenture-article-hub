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

// FindAllArticles retrieves all articles.
func (uc *ArticleUsecase) FindAllArticles(ctx context.Context) (*FindByCriteriaOutput, error) {
	articles, err := uc.repo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to find all articles: %w", err)
	}

	var articleOutputs []FindArticleByIDOutput
	for _, article := range articles {
		articleOutputs = append(articleOutputs, FindArticleByIDOutput{
			ID:           article.ID,
			Title:        article.Title.String(),
			Body:         article.Body.String(),
			Status:       article.Status.String(),
			ProviderType: article.ProviderType.String(),
			Link:         article.Link.String(),
			CreatedAt:    article.CreatedAt,
			UpdatedAt:    article.UpdatedAt,
		})
	}

	return &FindByCriteriaOutput{
		Articles:   articleOutputs,
		Total:      int64(len(articleOutputs)),
		Page:       1,
		Limit:      len(articleOutputs),
		TotalPages: 1,
	}, nil
}

// FindByCriteria retrieves articles based on the given criteria.
func (uc *ArticleUsecase) FindByCriteria(ctx context.Context, criteria FindByCriteriaInput) (*FindByCriteriaOutput, error) {
	// Convert input criteria to repository criteria
	repoCriteria := repository.ArticleQueryCriteria{
		Status:         criteria.Status,
		ProviderType:   criteria.ProviderType,
		SortBy:         criteria.SortBy,
		SortOrder:      criteria.SortOrder,
		Page:           criteria.Page,
		Limit:          criteria.Limit,
		IncludeDeleted: false, // Assuming we don't want to include deleted articles by default
	}

	articles, totalCount, err := uc.repo.FindByCriteria(ctx, repoCriteria)
	if err != nil {
		return nil, fmt.Errorf("failed to find articles by criteria: %w", err)
	}

	// Convert entities to output format
	var articleOutputs []FindArticleByIDOutput
	for _, article := range articles {
		articleOutputs = append(articleOutputs, FindArticleByIDOutput{
			ID:           article.ID,
			Title:        article.Title.String(),
			Body:         article.Body.String(),
			Status:       article.Status.String(),
			ProviderType: article.ProviderType.String(),
			Link:         article.Link.String(),
			CreatedAt:    article.CreatedAt,
			UpdatedAt:    article.UpdatedAt,
		})
	}

	// Calculate total pages
	totalPages := 0
	if criteria.Limit > 0 {
		totalPages = (totalCount + criteria.Limit - 1) / criteria.Limit // ceiling division
	}

	return &FindByCriteriaOutput{
		Articles:   articleOutputs,
		Total:      int64(totalCount),
		Page:       criteria.Page,
		Limit:      criteria.Limit,
		TotalPages: totalPages,
	}, nil
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

	newArticle, err := uc.repo.Create(ctx, articleEntity)
	if err != nil {
		return nil, err
	}

	return &CreateArticleOutput{
		ID:           newArticle.ID,
		Title:        newArticle.Title.String(),
		Body:         newArticle.Body.String(),
		Status:       newArticle.Status.String(),
		ProviderType: newArticle.ProviderType.String(),
		Link:         newArticle.Link.String(),
		CreatedAt:    newArticle.CreatedAt,
		UpdatedAt:    newArticle.UpdatedAt,
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
	return uc.repo.Delete(ctx, entity.ID)
}
