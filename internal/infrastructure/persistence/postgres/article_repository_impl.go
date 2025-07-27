package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/umekikazuya/momenture-article-hub/internal/domain/entity"
	"github.com/umekikazuya/momenture-article-hub/internal/domain/repository"
	"github.com/umekikazuya/momenture-article-hub/internal/domain/vo"
)

// ArticleModel はデータベースのarticlesテーブルに対応するモデル
type ArticleModel struct {
	ID           uint64     `gorm:"column:id;primaryKey;autoIncrement"`
	Title        string     `gorm:"column:title;not null;size:255"`
	Body         *string    `gorm:"column:body;type:text"`
	Status       string     `gorm:"column:status;not null;size:20"`
	ProviderType *string    `gorm:"column:provider_type;size:50"`
	Link         *string    `gorm:"column:link;type:text"`
	CreatedAt    time.Time  `gorm:"column:created_at;not null"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;not null"`
	DeletedAt    *time.Time `gorm:"column:deleted_at"`
}

// TableName はテーブル名を指定
func (ArticleModel) TableName() string {
	return "articles"
}

type PostgresArticleRepository struct {
	db *gorm.DB
}

func NewPostgresArticleRepository(db *gorm.DB) repository.ArticleRepository {
	return &PostgresArticleRepository{db: db}
}

// toEntity はArticleModelからentity.Articleに変換する
func (r *PostgresArticleRepository) toEntity(model *ArticleModel) (*entity.Article, error) {
	if model == nil {
		return nil, nil
	}

	// 必須フィールドのバリデーション
	title, err := vo.NewArticleTitle(model.Title)
	if err != nil {
		return nil, err
	}

	status := vo.ArticleStatus(model.Status)
	if !status.IsValid() {
		return nil, fmt.Errorf("invalid article status: %s", model.Status)
	}

	article := &entity.Article{
		ID:        model.ID,
		Title:     title,
		Status:    status,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
		DeletedAt: model.DeletedAt,
	}

	// オプショナルフィールドの設定
	if model.Body != nil {
		body, err := vo.NewArticleBody(model.Body)
		if err != nil {
			return nil, err
		}
		article.Body = body
	}

	if model.ProviderType != nil {
		providerType, err := vo.NewProviderType(model.ProviderType)
		if err != nil {
			return nil, err
		}
		article.ProviderType = providerType
	}

	if model.Link != nil {
		link, err := vo.NewLink(model.Link)
		if err != nil {
			return nil, err
		}
		article.Link = link
	}

	return article, nil
}

// toModel はentity.ArticleからArticleModelに変換する
func (r *PostgresArticleRepository) toModel(article *entity.Article) *ArticleModel {
	if article == nil {
		return nil
	}

	model := &ArticleModel{
		ID:        article.ID,
		Title:     article.Title.String(),
		Status:    article.Status.String(),
		CreatedAt: article.CreatedAt,
		UpdatedAt: article.UpdatedAt,
		DeletedAt: article.DeletedAt,
	}

	if article.Body != nil {
		body := article.Body.String()
		model.Body = &body
	}

	if article.ProviderType != nil {
		providerType := article.ProviderType.String()
		model.ProviderType = &providerType
	}

	if article.Link != nil {
		link := article.Link.String()
		model.Link = &link
	}

	return model
}

// FindAll は全ての記事を取得する
func (r *PostgresArticleRepository) FindAll(ctx context.Context) ([]*entity.Article, error) {
	var models []ArticleModel
	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to find all articles: %w", err)
	}

	articles := make([]*entity.Article, 0, len(models))
	for _, model := range models {
		article, err := r.toEntity(&model)
		if err != nil {
			return nil, fmt.Errorf("failed to convert model to entity: %w", err)
		}
		articles = append(articles, article)
	}

	return articles, nil
}

// FindByID は指定されたIDの記事を取得する
func (r *PostgresArticleRepository) FindByID(ctx context.Context, id uint64) (*entity.Article, error) {
	var model ArticleModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find article by id %d: %w", id, err)
	}

	article, err := r.toEntity(&model)
	if err != nil {
		return nil, fmt.Errorf("failed to convert model to entity: %w", err)
	}

	return article, nil
}

// FindByCriteria は指定された条件に一致する記事を取得する
func (r *PostgresArticleRepository) FindByCriteria(ctx context.Context, criteria repository.ArticleQueryCriteria) ([]*entity.Article, int, error) {
	query := r.db.WithContext(ctx).Model(&ArticleModel{})

	// 削除されたレコードの扱い
	if !criteria.IncludeDeleted {
		query = query.Where("deleted_at IS NULL")
	}

	// ステータスフィルタ
	if criteria.Status != nil {
		query = query.Where("status = ?", *criteria.Status)
	}

	// プロバイダータイプフィルタ
	if criteria.ProviderType != nil {
		query = query.Where("provider_type = ?", *criteria.ProviderType)
	}

	// 総件数を取得
	var totalCount int64
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count articles: %w", err)
	}

	// ソート順設定
	if criteria.SortBy != nil {
		sortOrder := "ASC"
		if criteria.SortOrder != nil && *criteria.SortOrder == "DESC" {
			sortOrder = "DESC"
		}
		query = query.Order(fmt.Sprintf("%s %s", *criteria.SortBy, sortOrder))
	} else {
		// デフォルトはcreated_atの降順
		query = query.Order("created_at DESC")
	}

	// ページネーション
	if criteria.Limit > 0 {
		query = query.Limit(criteria.Limit)
	}
	if criteria.Page > 0 {
		offset := (criteria.Page - 1) * criteria.Limit
		query = query.Offset(offset)
	}

	var models []ArticleModel
	if err := query.Find(&models).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to find articles by criteria: %w", err)
	}

	articles := make([]*entity.Article, 0, len(models))
	for _, model := range models {
		article, err := r.toEntity(&model)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to convert model to entity: %w", err)
		}
		articles = append(articles, article)
	}

	return articles, int(totalCount), nil
}

// Create は新しい記事を作成する
func (r *PostgresArticleRepository) Create(ctx context.Context, article *entity.Article) (*entity.Article, error) {
	if article == nil {
		return nil, fmt.Errorf("article cannot be nil")
	}

	model := r.toModel(article)
	model.CreatedAt = time.Now()
	model.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, fmt.Errorf("failed to create article: %w", err)
	}

	// 作成された記事を返すために再度エンティティに変換
	createdArticle, err := r.toEntity(model)
	if err != nil {
		return nil, fmt.Errorf("failed to convert created model to entity: %w", err)
	}

	return createdArticle, nil
}

// Update は既存の記事を更新する
func (r *PostgresArticleRepository) Update(ctx context.Context, article *entity.Article) error {
	if article == nil {
		return fmt.Errorf("article cannot be nil")
	}

	model := r.toModel(article)
	model.UpdatedAt = time.Now()

	// IDで検索して更新
	result := r.db.WithContext(ctx).Model(&ArticleModel{}).Where("id = ?", article.ID).Updates(model)
	if result.Error != nil {
		return fmt.Errorf("failed to update article: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("article with id %d not found", article.ID)
	}

	return nil
}

// Delete は指定されたIDの記事を論理削除する
func (r *PostgresArticleRepository) Delete(ctx context.Context, id uint64) error {
	now := time.Now()
	result := r.db.WithContext(ctx).Model(&ArticleModel{}).Where("id = ? AND deleted_at IS NULL", id).Update("deleted_at", now)

	if result.Error != nil {
		return fmt.Errorf("failed to delete article: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("article with id %d not found or already deleted", id)
	}

	return nil
}
