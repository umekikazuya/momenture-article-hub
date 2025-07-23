package article_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/umekikazuya/momenture-article-hub/internal/domain/entity"
	"github.com/umekikazuya/momenture-article-hub/internal/domain/repository"
	"github.com/umekikazuya/momenture-article-hub/internal/domain/vo"
	"github.com/umekikazuya/momenture-article-hub/internal/usecase/article"
)

type MockArticleRepository struct {
	mock.Mock
}

func (m *MockArticleRepository) FindAll(ctx context.Context) ([]*entity.Article, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*entity.Article), args.Error(1)
}

func (m *MockArticleRepository) FindByID(ctx context.Context, id uint64) (*entity.Article, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Article), args.Error(1)
}

func (m *MockArticleRepository) Create(ctx context.Context, article *entity.Article) (*entity.Article, error) {
	args := m.Called(ctx, article)
	return args.Get(0).(*entity.Article), args.Error(1)
}

func (m *MockArticleRepository) Update(ctx context.Context, article *entity.Article) error {
	args := m.Called(ctx, article)
	return args.Error(0)
}

func (m *MockArticleRepository) Delete(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockArticleRepository) FindByCriteria(ctx context.Context, criteria repository.ArticleQueryCriteria) ([]*entity.Article, int, error) {
	args := m.Called(ctx, criteria)
	return args.Get(0).([]*entity.Article), args.Get(1).(int), args.Error(2)
}

func ptr[T any](v T) *T {
	return &v
}

func TestArticleUsecase_CreateArticle(t *testing.T) {
	ctx := context.Background()

	t.Run("必須フィールドのみで記事を作成", func(t *testing.T) {
		mockRepo := new(MockArticleRepository)
		uc := article.NewArticleUsecase(mockRepo)

		input := article.CreateArticleInput{
			Title:  "テスト記事タイトル",
			Status: "draft",
		}

		createdArticle := &entity.Article{
			ID:        1,
			Title:     vo.ArticleTitle("テスト記事タイトル"),
			Status:    vo.ArticleStatus("draft"),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.Article")).Return(createdArticle, nil)

		output, err := uc.CreateArticle(ctx, input)

		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, uint64(1), output.ID)
		assert.Equal(t, input.Title, output.Title)
		assert.Equal(t, input.Status, output.Status)

		mockRepo.AssertExpectations(t)
	})

	t.Run("全てのフィールドを指定して記事を作成", func(t *testing.T) {
		mockRepo := new(MockArticleRepository)
		uc := article.NewArticleUsecase(mockRepo)

		input := article.CreateArticleInput{
			Title:        "タイトル",
			Body:         ptr("本文"),
			Status:       "published",
			ProviderType: ptr("qiita"),
			Link:         ptr("https://example.com"),
		}

		body, err := vo.NewArticleBody(ptr("本文"))
		require.NoError(t, err)
		providerType, err := vo.NewProviderType(ptr("qiita"))
		require.NoError(t, err)
		link, err := vo.NewLink(ptr("https://example.com"))
		require.NoError(t, err)

		createdArticle := &entity.Article{
			ID:           2,
			Title:        vo.ArticleTitle(input.Title),
			Body:         body,
			Status:       vo.ArticleStatus(input.Status),
			ProviderType: providerType,
			Link:         link,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.Article")).Return(createdArticle, nil)

		output, err := uc.CreateArticle(ctx, input)

		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, uint64(2), output.ID)
		assert.Equal(t, input.Title, output.Title)
		assert.Equal(t, *input.Body, output.Body)
		assert.Equal(t, input.Status, output.Status)
		assert.Equal(t, *input.ProviderType, output.ProviderType)
		assert.Equal(t, *input.Link, output.Link)
		assert.WithinDuration(t, createdArticle.CreatedAt, output.CreatedAt, 2*time.Second)

		mockRepo.AssertExpectations(t)
	})

	t.Run("タイトルが空の場合はエラー", func(t *testing.T) {
		mockRepo := new(MockArticleRepository)
		uc := article.NewArticleUsecase(mockRepo)

		input := article.CreateArticleInput{
			Title:  "",
			Status: "draft",
		}

		output, err := uc.CreateArticle(ctx, input)

		assert.Error(t, err)
		assert.Nil(t, output)
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("タイトルが文字数制限を超える場合はエラー", func(t *testing.T) {
		mockRepo := new(MockArticleRepository)
		uc := article.NewArticleUsecase(mockRepo)

		longTitle := strings.Repeat("a", vo.MaxArticleTitleLength+1)
		input := article.CreateArticleInput{
			Title:  longTitle,
			Status: "draft",
		}

		output, err := uc.CreateArticle(ctx, input)

		assert.Error(t, err)
		assert.Nil(t, output)
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("本文が空文字列で入力された場合の本文の返り値は空文字", func(t *testing.T) {
		mockRepo := new(MockArticleRepository)
		uc := article.NewArticleUsecase(mockRepo)

		input := article.CreateArticleInput{
			Title:        "Valid Title",
			Body:         ptr(""),
			Status:       "draft",
			ProviderType: ptr("qiita"),
			Link:         ptr("https://example.com"),
		}

		body, err := vo.NewArticleBody(ptr(""))
		require.NoError(t, err)
		providerType, err := vo.NewProviderType(ptr("qiita"))
		require.NoError(t, err)
		link, err := vo.NewLink(ptr("https://example.com"))
		require.NoError(t, err)

		createdArticle := &entity.Article{
			ID:           3,
			Title:        vo.ArticleTitle(input.Title),
			Body:         body,
			Status:       vo.ArticleStatus(input.Status),
			ProviderType: providerType,
			Link:         link,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.Article")).Return(createdArticle, nil)

		// When
		output, err := uc.CreateArticle(ctx, input)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, uint64(3), output.ID)
		assert.Equal(t, input.Title, output.Title)
		assert.Equal(t, *input.Body, output.Body)
		assert.Equal(t, input.Status, output.Status)
		assert.Equal(t, *input.ProviderType, output.ProviderType)
		assert.Equal(t, *input.Link, output.Link)
		assert.WithinDuration(t, createdArticle.CreatedAt, output.CreatedAt, 2*time.Second)

		mockRepo.AssertExpectations(t)
	})

	t.Run("本文がnilで入力された場合の本文の返り値は空文字", func(t *testing.T) {
		mockRepo := new(MockArticleRepository)
		uc := article.NewArticleUsecase(mockRepo)

		input := article.CreateArticleInput{
			Title:        "Valid Title",
			Body:         nil,
			Status:       "draft",
			ProviderType: ptr("qiita"),
			Link:         ptr("https://example.com"),
		}

		body, err := vo.NewArticleBody(ptr(""))
		require.NoError(t, err)
		providerType, err := vo.NewProviderType(ptr("qiita"))
		require.NoError(t, err)
		link, err := vo.NewLink(ptr("https://example.com"))
		require.NoError(t, err)

		createdArticle := &entity.Article{
			ID:           3,
			Title:        vo.ArticleTitle(input.Title),
			Body:         body,
			Status:       vo.ArticleStatus(input.Status),
			ProviderType: providerType,
			Link:         link,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.Article")).Return(createdArticle, nil)

		// When
		output, err := uc.CreateArticle(ctx, input)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, uint64(3), output.ID)
		assert.Equal(t, input.Title, output.Title)
		assert.Equal(t, "", output.Body)
		assert.Equal(t, input.Status, output.Status)
		assert.Equal(t, *input.ProviderType, output.ProviderType)
		assert.Equal(t, *input.Link, output.Link)
		assert.WithinDuration(t, createdArticle.CreatedAt, output.CreatedAt, 2*time.Second)

		mockRepo.AssertExpectations(t)
	})

	t.Run("無効なステータスの場合はエラー", func(t *testing.T) {
		mockRepo := new(MockArticleRepository)
		uc := article.NewArticleUsecase(mockRepo)

		input := article.CreateArticleInput{
			Title:  "Valid Title",
			Status: "invalid_status",
		}

		output, err := uc.CreateArticle(ctx, input)

		assert.Error(t, err)
		assert.Nil(t, output)
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("無効なプロバイダタイプの場合はエラー", func(t *testing.T) {
		mockRepo := new(MockArticleRepository)
		uc := article.NewArticleUsecase(mockRepo)

		input := article.CreateArticleInput{
			Title:        "Valid Title",
			Status:       "draft",
			ProviderType: ptr("unknown_provider"),
		}

		output, err := uc.CreateArticle(ctx, input)

		assert.Error(t, err)
		assert.Nil(t, output)
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("リポジトリでエラーが発生した場合は適切に処理される", func(t *testing.T) {
		mockRepo := new(MockArticleRepository)
		uc := article.NewArticleUsecase(mockRepo)

		input := article.CreateArticleInput{
			Title:  "Valid Title",
			Status: "draft",
		}

		dbError := fmt.Errorf("db error")
		mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.Article")).Return((*entity.Article)(nil), dbError)

		output, err := uc.CreateArticle(ctx, input)

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.Contains(t, err.Error(), "db error")

		mockRepo.AssertExpectations(t)
	})
}

func TestArticleUsecase_GetArticles(t *testing.T) {
	t.Run("全ての記事をページネーションなしで取得", func(t *testing.T) {
		mockRepo := new(MockArticleRepository)
		uc := article.NewArticleUsecase(mockRepo)

		expectedArticles := []*entity.Article{
			{ID: 1, Title: vo.ArticleTitle("Article 1"), Status: vo.ArticleStatus("draft")},
			{ID: 2, Title: vo.ArticleTitle("Article 2"), Status: vo.ArticleStatus("published")},
		}

		mockRepo.On("FindAll", mock.Anything).Return(expectedArticles, nil)

		output, err := uc.FindAllArticles(context.Background())

		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.Len(t, output, len(expectedArticles))
		assert.Equal(t, int64(len(expectedArticles)), int64(len(output)))
	})

	t.Run("ステータスでフィルタリング", func(t *testing.T) {
		mockRepo := new(MockArticleRepository)
		uc := article.NewArticleUsecase(mockRepo)

		status := "draft"
		expectedArticles := []*entity.Article{
			{ID: 1, Title: vo.ArticleTitle("Draft Article 1"), Status: vo.ArticleStatus(status)},
			{ID: 2, Title: vo.ArticleTitle("Draft Article 2"), Status: vo.ArticleStatus(status)},
		}

		mockRepo.On("FindByCriteria", mock.Anything, repository.ArticleQueryCriteria{
			Status: &status,
			Page:   1,
			Limit:  10,
		}).Return(expectedArticles, len(expectedArticles), nil)

		input := article.FindByCriteriaInput{
			Status: ptr(status),
			Page:   1,
			Limit:  10,
		}

		output, err := uc.FindByCriteria(context.Background(), input)

		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.Len(t, output.Articles, len(expectedArticles))
		assert.Equal(t, int64(len(expectedArticles)), output.Total)
	})

	t.Run("プロバイダタイプでフィルタリング", func(t *testing.T) {
		mockRepo := new(MockArticleRepository)
		uc := article.NewArticleUsecase(mockRepo)

		providerType := "qiita"
		providerTypeVal := vo.ProviderType(providerType)
		expectedArticles := []*entity.Article{
			{ID: 1, Title: vo.ArticleTitle("Qiita Article 1"), ProviderType: &providerTypeVal},
			{ID: 2, Title: vo.ArticleTitle("Qiita Article 2"), ProviderType: &providerTypeVal},
		}

		mockRepo.On("FindByCriteria", mock.Anything, repository.ArticleQueryCriteria{
			ProviderType: &providerType,
			Page:         1,
			Limit:        10,
		}).Return(expectedArticles, len(expectedArticles), nil)

		input := article.FindByCriteriaInput{
			ProviderType: ptr(providerType),
			Page:         1,
			Limit:        10,
		}

		output, err := uc.FindByCriteria(context.Background(), input)

		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.Len(t, output.Articles, len(expectedArticles))
		assert.Equal(t, int64(len(expectedArticles)), output.Total)
	})

	t.Run("ソート順を指定して取得", func(t *testing.T) {
		mockRepo := new(MockArticleRepository)
		uc := article.NewArticleUsecase(mockRepo)

		expectedArticles := []*entity.Article{
			{ID: 1, Title: vo.ArticleTitle("Article 1"), CreatedAt: time.Now().Add(-time.Hour)},
			{ID: 2, Title: vo.ArticleTitle("Article 2"), CreatedAt: time.Now()},
		}

		mockRepo.On("FindByCriteria", mock.Anything, repository.ArticleQueryCriteria{
			SortBy:    ptr("created_at"),
			SortOrder: ptr("desc"),
			Page:      1,
			Limit:     10,
		}).Return(expectedArticles, len(expectedArticles), nil)

		input := article.FindByCriteriaInput{
			SortBy:    ptr("created_at"),
			SortOrder: ptr("desc"),
			Page:      1,
			Limit:     10,
		}

		output, err := uc.FindByCriteria(context.Background(), input)

		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.Len(t, output.Articles, len(expectedArticles))
		assert.Equal(t, int64(len(expectedArticles)), output.Total)
	})

	t.Run("ページネーション適用", func(t *testing.T) {
		mockRepo := new(MockArticleRepository)
		uc := article.NewArticleUsecase(mockRepo)

		expectedArticles := []*entity.Article{
			{ID: 1, Title: vo.ArticleTitle("Article 1"), CreatedAt: time.Now().Add(-time.Hour)},
			{ID: 2, Title: vo.ArticleTitle("Article 2"), CreatedAt: time.Now()},
		}

		mockRepo.On("FindByCriteria", mock.Anything, repository.ArticleQueryCriteria{
			Page:  1,
			Limit: 1,
		}).Return(expectedArticles[:1], 2, nil)

		input := article.FindByCriteriaInput{
			Page:  1,
			Limit: 1,
		}

		output, err := uc.FindByCriteria(context.Background(), input)

		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.Len(t, output.Articles, 1)
		assert.Equal(t, int64(2), output.Total)
	})

	t.Run("記事が見つからない場合は空の結果", func(t *testing.T) {
		mockRepo := new(MockArticleRepository)
		uc := article.NewArticleUsecase(mockRepo)

		mockRepo.On("FindByCriteria", mock.Anything, repository.ArticleQueryCriteria{
			Page:  1,
			Limit: 10,
		}).Return([]*entity.Article{}, 0, nil)

		input := article.FindByCriteriaInput{
			Page:  1,
			Limit: 10,
		}

		output, err := uc.FindByCriteria(context.Background(), input)

		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.Len(t, output.Articles, 0)
		assert.Equal(t, int64(0), output.Total)

		mockRepo.AssertExpectations(t)
	})

	t.Run("リポジトリエラーは適切に処理される", func(t *testing.T) {
		mockRepo := new(MockArticleRepository)
		uc := article.NewArticleUsecase(mockRepo)

		mockRepo.On("FindByCriteria", mock.Anything, repository.ArticleQueryCriteria{
			Page:  1,
			Limit: 10,
		}).Return(([]*entity.Article)(nil), 0, fmt.Errorf("db error"))

		input := article.FindByCriteriaInput{
			Page:  1,
			Limit: 10,
		}

		output, err := uc.FindByCriteria(context.Background(), input)

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.Contains(t, err.Error(), "db error")

		mockRepo.AssertExpectations(t)
	})
}

func TestArticleUsecase_GetArticleByID(t *testing.T) {
	ctx := context.Background()
	t.Run("IDで記事を取得", func(t *testing.T) {
		mockRepo := new(MockArticleRepository)
		uc := article.NewArticleUsecase(mockRepo)

		articleID := uint64(1)
		title, err := vo.NewArticleTitle("Test Article")
		require.NoError(t, err)
		body, err := vo.NewArticleBody(ptr("This is a test article body."))
		require.NoError(t, err)
		status := vo.ArticleStatus("draft")
		providerType, err := vo.NewProviderType(ptr("qiita"))
		require.NoError(t, err)
		link, err := vo.NewLink(ptr("https://example.com"))
		require.NoError(t, err)

		expectedArticle := &entity.Article{
			ID:           articleID,
			Title:        title,
			Body:         body,
			Status:       status,
			ProviderType: providerType,
			Link:         link,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		mockRepo.On("FindByID", ctx, articleID).Return(expectedArticle, nil)

		output, err := uc.FindArticleByID(ctx, articleID)

		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, expectedArticle.ID, output.ID)
		assert.Equal(t, expectedArticle.Title.String(), output.Title)
		assert.Equal(t, expectedArticle.Body.String(), output.Body)
		assert.Equal(t, expectedArticle.Status.String(), output.Status)
		assert.Equal(t, expectedArticle.ProviderType.String(), output.ProviderType)
		assert.Equal(t, expectedArticle.Link.String(), output.Link)
		assert.WithinDuration(t, expectedArticle.CreatedAt, output.CreatedAt, 2*time.Second)
		assert.WithinDuration(t, expectedArticle.UpdatedAt, output.UpdatedAt, 2*time.Second)

		mockRepo.AssertExpectations(t)
	})

	t.Run("記事が見つからない場合はエラー", func(t *testing.T) {
		mockRepo := new(MockArticleRepository)
		uc := article.NewArticleUsecase(mockRepo)

		articleID := uint64(999) // 存在しないID
		mockRepo.On("FindByID", ctx, articleID).Return(nil, fmt.Errorf("article not found"))

		output, err := uc.FindArticleByID(ctx, articleID)

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.Contains(t, err.Error(), "article not found")

		mockRepo.AssertExpectations(t)
	})

	t.Run("リポジトリエラーは適切に処理される", func(t *testing.T) {
		mockRepo := new(MockArticleRepository)
		uc := article.NewArticleUsecase(mockRepo)

		articleID := uint64(1)
		mockRepo.On("FindByID", ctx, articleID).Return(nil, fmt.Errorf("db error"))

		output, err := uc.FindArticleByID(ctx, articleID)

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.Contains(t, err.Error(), "db error")

		mockRepo.AssertExpectations(t)
	})
}

func TestArticleUsecase_UpdateArticle(t *testing.T) {
	ctx := context.Background()

	t.Run("タイトルのみ更新", func(t *testing.T) {
		mockRepo := new(MockArticleRepository)
		uc := article.NewArticleUsecase(mockRepo)

		input := article.UpdateArticleInput{
			Title: ptr("Updated Title"),
		}

		// 既存の記事を作成
		existingArticle, err := entity.NewArticle("Original Title", "draft")
		require.NoError(t, err)
		existingArticle.ID = 1

		mockRepo.On("FindByID", ctx, uint64(1)).Return(existingArticle, nil)
		mockRepo.On("Update", ctx, mock.AnythingOfType("*entity.Article")).Return(nil)

		output, err := uc.UpdateArticle(ctx, 1, input)

		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, "Updated Title", output.Title)
		mockRepo.AssertExpectations(t)
	})

	t.Run("ステータスを下書きから公開済みに変更", func(t *testing.T) {
		mockRepo := new(MockArticleRepository)
		uc := article.NewArticleUsecase(mockRepo)

		input := article.UpdateArticleInput{
			Status: ptr("published"),
		}

		// 既存の記事を作成
		existingArticle, err := entity.NewArticle("Original Title", "draft")
		require.NoError(t, err)
		existingArticle.ID = 1

		mockRepo.On("FindByID", ctx, uint64(1)).Return(existingArticle, nil)
		mockRepo.On("Update", ctx, mock.AnythingOfType("*entity.Article")).Return(nil)

		output, err := uc.UpdateArticle(ctx, 1, input)

		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, "published", output.Status)
		mockRepo.AssertExpectations(t)
	})

	t.Run("オプションフィールドをクリア", func(t *testing.T) {
		mockRepo := new(MockArticleRepository)
		uc := article.NewArticleUsecase(mockRepo)

		input := article.UpdateArticleInput{
			Body:         ptr(""),
			ProviderType: nil, // 空文字列ではなくnilでクリア
			Link:         nil, // 空文字列ではなくnilでクリア
		}

		// 既存の記事を作成（オプションフィールド付き）
		existingArticle, err := entity.NewArticle("Original Title", "draft",
			entity.WithBody(ptr("Original Body")),
			entity.WithProviderType(ptr("qiita")),
			entity.WithLink(ptr("https://original.com")),
		)
		require.NoError(t, err)
		existingArticle.ID = 1

		mockRepo.On("FindByID", ctx, uint64(1)).Return(existingArticle, nil)
		mockRepo.On("Update", ctx, mock.AnythingOfType("*entity.Article")).Return(nil)

		output, err := uc.UpdateArticle(ctx, 1, input)

		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.Empty(t, output.Body)
		assert.Empty(t, output.ProviderType)
		assert.Empty(t, output.Link)
		mockRepo.AssertExpectations(t)
	})

	t.Run("記事が見つからない場合はエラー", func(t *testing.T) {
		mockRepo := new(MockArticleRepository)
		uc := article.NewArticleUsecase(mockRepo)

		input := article.UpdateArticleInput{
			Body:         ptr("Updated Body"),
			ProviderType: ptr("Updated Provider"),
			Link:         ptr("Updated Link"),
		}

		mockRepo.On("FindByID", ctx, uint64(1)).Return(nil, fmt.Errorf("article not found"))

		output, err := uc.UpdateArticle(ctx, 1, input)

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.Contains(t, err.Error(), "article not found")

		mockRepo.AssertExpectations(t)
	})

	t.Run("更新内容のバリデーションエラー", func(t *testing.T) {
		mockRepo := new(MockArticleRepository)
		uc := article.NewArticleUsecase(mockRepo)

		input := article.UpdateArticleInput{
			// タイトルが長すぎる
			Title: ptr(strings.Repeat("a", vo.MaxArticleTitleLength+1)),
		}

		// 既存の記事を作成
		existingArticle, err := entity.NewArticle("Original Title", "draft")
		require.NoError(t, err)
		existingArticle.ID = 1

		mockRepo.On("FindByID", ctx, uint64(1)).Return(existingArticle, nil)

		output, err := uc.UpdateArticle(ctx, 1, input)

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.Contains(t, err.Error(), "title exceeds maximum length")

		mockRepo.AssertExpectations(t)
	})

	t.Run("ドメインルール違反の場合はエラー", func(t *testing.T) {
		mockRepo := new(MockArticleRepository)
		uc := article.NewArticleUsecase(mockRepo)

		input := article.UpdateArticleInput{
			Status: ptr("invalid_status"), // 無効なステータス
		}

		// 既存の記事を作成
		existingArticle, err := entity.NewArticle("Original Title", "draft")
		require.NoError(t, err)
		existingArticle.ID = 1

		mockRepo.On("FindByID", ctx, uint64(1)).Return(existingArticle, nil)

		output, err := uc.UpdateArticle(ctx, 1, input)

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.Contains(t, err.Error(), "invalid status")

		mockRepo.AssertExpectations(t)
	})

	t.Run("リポジトリエラーは適切に処理される", func(t *testing.T) {
		mockRepo := new(MockArticleRepository)
		uc := article.NewArticleUsecase(mockRepo)

		input := article.UpdateArticleInput{
			Title: ptr("Updated Title"),
		}

		// 既存の記事を作成
		existingArticle, err := entity.NewArticle("Original Title", "draft")
		require.NoError(t, err)
		existingArticle.ID = 1

		mockRepo.On("FindByID", ctx, uint64(1)).Return(existingArticle, nil)
		mockRepo.On("Update", ctx, mock.AnythingOfType("*entity.Article")).Return(fmt.Errorf("repository error"))

		output, err := uc.UpdateArticle(ctx, 1, input)

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.Contains(t, err.Error(), "repository error")

		mockRepo.AssertExpectations(t)
	})
}

func TestArticleUsecase_DeleteArticle(t *testing.T) {
	t.Run("記事が削除される", func(t *testing.T) {
		mockRepo := new(MockArticleRepository)
		uc := article.NewArticleUsecase(mockRepo)

		id := uint64(1)

		// 既存の記事を作成
		existingArticle, err := entity.NewArticle("Test Article", "draft")
		require.NoError(t, err)
		existingArticle.ID = id

		mockRepo.On("FindByID", mock.Anything, id).Return(existingArticle, nil)
		mockRepo.On("Delete", mock.Anything, id).Return(nil)

		err = uc.DeleteArticle(context.Background(), id)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
	t.Run("記事が見つからない場合はエラー", func(t *testing.T) {
		mockRepo := new(MockArticleRepository)
		uc := article.NewArticleUsecase(mockRepo)

		id := uint64(1)
		mockRepo.On("FindByID", mock.Anything, id).Return(nil, fmt.Errorf("article not found"))

		err := uc.DeleteArticle(context.Background(), id)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "article not found")
		mockRepo.AssertExpectations(t)
	})
}
