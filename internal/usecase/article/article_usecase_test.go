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

func (m *MockArticleRepository) Create(ctx context.Context, article *entity.Article) (uint64, error) {
	args := m.Called(ctx, article)
	return args.Get(0).(uint64), args.Error(1)
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

		mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.Article")).Return(uint64(1), nil).Run(func(args mock.Arguments) {
			article := args.Get(1).(*entity.Article)
			assert.Equal(t, input.Title, article.Title.String())
			assert.Equal(t, input.Status, article.Status.String())
			assert.WithinDuration(t, time.Now(), article.CreatedAt, 2*time.Second)
			assert.WithinDuration(t, time.Now(), article.UpdatedAt, 2*time.Second)
		})

		output, err := uc.CreateArticle(ctx, input)

		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, uint64(1), output.ID)
		assert.NotEqual(t, uint64(0), output.ID)
		assert.Equal(t, input.Title, output.Title)
		assert.WithinDuration(t, time.Now(), output.CreatedAt, 2*time.Second)
		assert.WithinDuration(t, time.Now(), output.UpdatedAt, 2*time.Second)

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

		mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.Article")).Return(uint64(2), nil).Run(func(args mock.Arguments) {
			article := args.Get(1).(*entity.Article)
			assert.Equal(t, input.Title, article.Title.String())
			assert.Equal(t, input.Status, article.Status.String())
			require.NotNil(t, article.Body)
			assert.Equal(t, *input.Body, article.Body.String())
			require.NotNil(t, article.ProviderType)
			assert.Equal(t, vo.ProviderType(*input.ProviderType), *article.ProviderType)
			require.NotNil(t, article.Link)
			assert.Equal(t, *input.Link, article.Link.String())
		})

		output, err := uc.CreateArticle(ctx, input)

		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, uint64(2), output.ID)

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

	t.Run("本文が空文字列の場合はBodyはnil", func(t *testing.T) {
		mockRepo := new(MockArticleRepository)
		uc := article.NewArticleUsecase(mockRepo)

		input := article.CreateArticleInput{
			Title:  "Valid Title",
			Body:   ptr(""),
			Status: "draft",
		}

		mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.Article")).Return(uint64(3), nil).Run(func(args mock.Arguments) {
			article := args.Get(1).(*entity.Article)
			assert.Equal(t, input.Title, article.Title.String())
			assert.Equal(t, input.Status, article.Status.String())
			assert.Nil(t, article.Body)
		})

		// When
		output, err := uc.CreateArticle(ctx, input)

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, uint64(3), output.ID)

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
		mockRepo.On("Create", ctx, mock.AnythingOfType("*entity.Article")).Return(uint64(0), dbError)

		output, err := uc.CreateArticle(ctx, input)

		assert.Error(t, err)
		assert.Nil(t, output)
		assert.Contains(t, err.Error(), "db error")

		mockRepo.AssertExpectations(t)
	})
}

func TestArticleUsecase_GetArticles(t *testing.T) {
	t.Run("全ての記事をページネーションなしで取得", func(t *testing.T) {
		t.Skip("GetArticles method not implemented yet - RED phase")
	})

	t.Run("ステータスでフィルタリング", func(t *testing.T) {
		t.Skip("GetArticles method not implemented yet - RED phase")
	})

	t.Run("プロバイダタイプでフィルタリング", func(t *testing.T) {
		t.Skip("GetArticles method not implemented yet - RED phase")
	})

	t.Run("ソート順を指定して取得", func(t *testing.T) {
		t.Skip("GetArticles method not implemented yet - RED phase")
	})

	t.Run("ページネーション適用", func(t *testing.T) {
		t.Skip("GetArticles method not implemented yet - RED phase")
	})

	t.Run("記事が見つからない場合は空の結果", func(t *testing.T) {
		t.Skip("GetArticles method not implemented yet - RED phase")
	})

	t.Run("リポジトリエラーは適切に処理される", func(t *testing.T) {
		t.Skip("GetArticles method not implemented yet - RED phase")
	})
}

func TestArticleUsecase_GetArticleByID(t *testing.T) {
	t.Run("IDで記事を取得", func(t *testing.T) {
		t.Skip("GetArticleByID method not implemented yet - RED phase")
	})

	t.Run("記事が見つからない場合はエラー", func(t *testing.T) {
		t.Skip("GetArticleByID method not implemented yet - RED phase")
	})

	t.Run("リポジトリエラーは適切に処理される", func(t *testing.T) {
		t.Skip("GetArticleByID method not implemented yet - RED phase")
	})
}

func TestArticleUsecase_UpdateArticle(t *testing.T) {
	t.Run("タイトルのみ更新", func(t *testing.T) {
		t.Skip("UpdateArticle method not implemented yet - RED phase")
	})

	t.Run("ステータスを下書きから公開済みに変更", func(t *testing.T) {
		t.Skip("UpdateArticle method not implemented yet - RED phase")
	})

	t.Run("オプションフィールドをクリア", func(t *testing.T) {
		t.Skip("UpdateArticle method not implemented yet - RED phase")
	})

	t.Run("記事が見つからない場合はエラー", func(t *testing.T) {
		t.Skip("UpdateArticle method not implemented yet - RED phase")
	})

	t.Run("更新内容のバリデーションエラー", func(t *testing.T) {
		t.Skip("UpdateArticle method not implemented yet - RED phase")
	})

	t.Run("ドメインルール違反の場合はエラー", func(t *testing.T) {
		t.Skip("UpdateArticle method not implemented yet - RED phase")
	})

	t.Run("リポジトリエラーは適切に処理される", func(t *testing.T) {
		t.Skip("UpdateArticle method not implemented yet - RED phase")
	})
}

func TestArticleUsecase_DeleteArticle(t *testing.T) {
	t.Run("記事が見つからない場合はエラー", func(t *testing.T) {
		t.Skip("DeleteArticle method not implemented yet - RED phase")
	})
}
