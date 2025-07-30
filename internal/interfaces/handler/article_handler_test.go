package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/umekikazuya/momenture-article-hub/internal/interfaces/handler"
	"github.com/umekikazuya/momenture-article-hub/internal/usecase/article"
)

type MockArticleUsecase struct {
	mock.Mock
}

func (m *MockArticleUsecase) CreateArticle(ctx context.Context, input article.CreateArticleInput) (*article.CreateArticleOutput, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*article.CreateArticleOutput), args.Error(1)
}

func (m *MockArticleUsecase) FindAllArticles(ctx context.Context) ([]article.FindArticleByIDOutput, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]article.FindArticleByIDOutput), args.Error(1)
}

func (m *MockArticleUsecase) FindByCriteria(ctx context.Context, criteria article.FindByCriteriaInput) (*article.FindByCriteriaOutput, error) {
	args := m.Called(ctx, criteria)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*article.FindByCriteriaOutput), args.Error(1)
}

// FindArticleByID はモックされたユースケースメソッドです。
func (m *MockArticleUsecase) FindArticleByID(ctx context.Context, id uint64) (*article.FindArticleByIDOutput, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*article.FindArticleByIDOutput), args.Error(1)
}

// UpdateArticle はモックされたユースケースメソッドです。
func (m *MockArticleUsecase) UpdateArticle(ctx context.Context, id uint64, input article.UpdateArticleInput) (*article.UpdateArticleOutput, error) {
	args := m.Called(ctx, id, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*article.UpdateArticleOutput), args.Error(1)
}

// DeleteArticle はモックされたユースケースメソッドです。
func (m *MockArticleUsecase) DeleteArticle(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// -----------------------------------------------------------------------------
// テストヘルパー関数
// -----------------------------------------------------------------------------

// ptr は任意の型のポインタを生成するヘルパー関数です。
func ptr[T any](v T) *T {
	return &v
}

// newTestServer はテスト用のGinルーターとハンドラーを設定します。
func newTestServer(mockUsecase *MockArticleUsecase) *gin.Engine {
	r := gin.Default()
	h := handler.NewArticleHandler(mockUsecase) // ⬅️ handler.NewArticleHandler はまだ存在しないため、コメントアウトまたは仮実装が必要

	// APIエンドポイントの登録
	r.POST("/articles", h.CreateArticle)
	r.GET("/articles", h.GetArticles)
	r.GET("/articles/:id", h.GetArticleByID)
	r.PUT("/articles/:id", h.UpdateArticle)
	r.DELETE("/articles/:id", h.DeleteArticle)

	return r
}

// sendRequest はHTTPリクエストを送信し、レスポンスを記録します。
func sendRequest(router *gin.Engine, method, path string, body interface{}) *httptest.ResponseRecorder {
	var req *http.Request
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		req, _ = http.NewRequest(method, path, bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, _ = http.NewRequest(method, path, nil)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

// -----------------------------------------------------------------------------
// 1.1. 記事の作成 (Create Article) ハンドラー テスト
// -----------------------------------------------------------------------------

func TestArticleHandler_CreateArticle(t *testing.T) {
	// ... (既存のテストケースは省略)
}

// -----------------------------------------------------------------------------
// 1.2. 特定の記事取得 (Get Article by ID) ハンドラー テスト
// -----------------------------------------------------------------------------
func TestArticleHandler_GetArticleByID(t *testing.T) {
	// シナリオ 1.2.1: 成功
	t.Run("成功", func(t *testing.T) {
		mockUsecase := new(MockArticleUsecase)
		router := newTestServer(mockUsecase)

		expectedArticle := &article.FindArticleByIDOutput{
			ID:        123,
			Title:     "Test Title",
			Body:      "Test Body",
			Status:    "published",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockUsecase.On("FindArticleByID", mock.Anything, uint64(123)).Return(expectedArticle, nil)

		w := sendRequest(router, "GET", "/articles/123", nil)

		assert.Equal(t, http.StatusOK, w.Code)
		var responseBody map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		require.NoError(t, err)
		assert.Equal(t, float64(123), responseBody["ID"])
		mockUsecase.AssertExpectations(t)
	})

	// シナリオ 1.2.2: 失敗 - 記事が見つからない
	t.Run("失敗 - 記事が見つからない", func(t *testing.T) {
		mockUsecase := new(MockArticleUsecase)
		router := newTestServer(mockUsecase)

		mockUsecase.On("FindArticleByID", mock.Anything, uint64(999)).Return(nil, fmt.Errorf("article not found"))

		w := sendRequest(router, "GET", "/articles/999", nil)

		assert.Equal(t, http.StatusNotFound, w.Code)
		var responseBody map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		require.NoError(t, err)
		assert.Equal(t, "article not found", responseBody["error"])
		mockUsecase.AssertExpectations(t)
	})

	// シナリオ 1.2.3: 失敗 - パスパラメータが無効
	t.Run("失敗 - パスパラメータが無効", func(t *testing.T) {
		mockUsecase := new(MockArticleUsecase)
		router := newTestServer(mockUsecase)

		w := sendRequest(router, "GET", "/articles/abc", nil)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var responseBody map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		require.NoError(t, err)
		assert.Equal(t, "Invalid article ID", responseBody["error"])
	})
}

// -----------------------------------------------------------------------------
// 1.3. 記事の一覧取得 (Get Articles) ハンドラー テスト
// -----------------------------------------------------------------------------
func TestArticleHandler_GetArticles(t *testing.T) {
	// シナリオ 1.3.1: 成功 - デフォルト条件
	t.Run("成功 - デフォルト条件", func(t *testing.T) {
		mockUsecase := new(MockArticleUsecase)
		router := newTestServer(mockUsecase)

		expectedOutput := &article.FindByCriteriaOutput{
			Articles: []article.FindArticleByIDOutput{
				{ID: 1, Title: "Article 1"},
			},
			Total: 1,
			Page:  1,
			Limit: 10,
		}
		expectedInput := article.FindByCriteriaInput{
			Page:  1,
			Limit: 10,
		}

		mockUsecase.On("FindByCriteria", mock.Anything, expectedInput).Return(expectedOutput, nil)

		w := sendRequest(router, "GET", "/articles", nil)

		assert.Equal(t, http.StatusOK, w.Code)
		// ... (レスポンスボディの検証)
		mockUsecase.AssertExpectations(t)
	})

	// シナリオ 1.3.2: 成功 - クエリパラメータ指定
	t.Run("成功 - クエリパラメータ指定", func(t *testing.T) {
		mockUsecase := new(MockArticleUsecase)
		router := newTestServer(mockUsecase)

		expectedInput := article.FindByCriteriaInput{
			Status: ptr("published"),
			Page:   2,
			Limit:  10,
		}
		mockUsecase.On("FindByCriteria", mock.Anything, expectedInput).Return(&article.FindByCriteriaOutput{}, nil)

		w := sendRequest(router, "GET", "/articles?status=published&page=2&limit=10", nil)

		assert.Equal(t, http.StatusOK, w.Code)
		mockUsecase.AssertExpectations(t)
	})

	// シナリオ 1.3.3: 失敗 - 不正なクエリパラメータ
	t.Run("失敗 - 不正なクエリパラメータ", func(t *testing.T) {
		mockUsecase := new(MockArticleUsecase)
		router := newTestServer(mockUsecase)

		w := sendRequest(router, "GET", "/articles?page=0", nil)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// -----------------------------------------------------------------------------
// 1.4. 記事の更新 (Update Article) ハンドラー テスト
// -----------------------------------------------------------------------------
func TestArticleHandler_UpdateArticle(t *testing.T) {
	// シナリオ 1.4.1: 成功 - タイトルのみ更新
	t.Run("成功 - タイトルのみ更新", func(t *testing.T) {
		mockUsecase := new(MockArticleUsecase)
		router := newTestServer(mockUsecase)

		updateInput := article.UpdateArticleInput{Title: ptr("新しいタイトル")}
		mockUsecase.On("UpdateArticle", mock.Anything, uint64(123), updateInput).Return(&article.UpdateArticleOutput{ID: 123, Title: "新しいタイトル"}, nil)

		w := sendRequest(router, "PUT", "/articles/123", gin.H{"title": "新しいタイトル"})

		assert.Equal(t, http.StatusOK, w.Code)
		mockUsecase.AssertExpectations(t)
	})

	// ... (他の更新シナリオ)

	// シナリオ 1.4.3: 失敗 - 記事が見つからない
	t.Run("失敗 - 記事が見つからない", func(t *testing.T) {
		mockUsecase := new(MockArticleUsecase)
		router := newTestServer(mockUsecase)

		updateInput := article.UpdateArticleInput{Title: ptr("新しいタイトル")}
		mockUsecase.On("UpdateArticle", mock.Anything, uint64(999), updateInput).Return(nil, fmt.Errorf("article not found"))

		w := sendRequest(router, "PUT", "/articles/999", gin.H{"title": "新しいタイトル"})

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockUsecase.AssertExpectations(t)
	})

	// シナリオ 1.4.4: 失敗 - 入力バリデーションエラー
	t.Run("失敗 - 入力バリデーションエラー", func(t *testing.T) {
		mockUsecase := new(MockArticleUsecase)
		router := newTestServer(mockUsecase)

		w := sendRequest(router, "PUT", "/articles/123", gin.H{"status": "invalid"})

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// -----------------------------------------------------------------------------
// 1.5. 記事の削除 (Soft Delete Article) ハンドラー テスト
// -----------------------------------------------------------------------------
func TestArticleHandler_DeleteArticle(t *testing.T) {
	// シナリオ 1.5.1: 成功
	t.Run("成功", func(t *testing.T) {
		mockUsecase := new(MockArticleUsecase)
		router := newTestServer(mockUsecase)

		mockUsecase.On("DeleteArticle", mock.Anything, uint64(123)).Return(nil)

		w := sendRequest(router, "DELETE", "/articles/123", nil)

		assert.Equal(t, http.StatusNoContent, w.Code)
		mockUsecase.AssertExpectations(t)
	})

	// シナリオ 1.5.2: 失敗 - 記事が見つからない
	t.Run("失敗 - 記事が見つからない", func(t *testing.T) {
		mockUsecase := new(MockArticleUsecase)
		router := newTestServer(mockUsecase)

		mockUsecase.On("DeleteArticle", mock.Anything, uint64(999)).Return(fmt.Errorf("article not found"))

		w := sendRequest(router, "DELETE", "/articles/999", nil)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockUsecase.AssertExpectations(t)
	})
}
