package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/umekikazuya/momenture-article-hub/internal/interfaces/handler"
	"github.com/umekikazuya/momenture-article-hub/internal/usecase/article"
)

// MockArticleUsecase is a mock for ArticleUsecase.
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

func (m *MockArticleUsecase) FindArticleByID(ctx context.Context, id uint64) (*article.FindArticleByIDOutput, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*article.FindArticleByIDOutput), args.Error(1)
}

func (m *MockArticleUsecase) FindByCriteria(ctx context.Context, criteria article.FindByCriteriaInput) (*article.FindByCriteriaOutput, error) {
	args := m.Called(ctx, criteria)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*article.FindByCriteriaOutput), args.Error(1)
}

func (m *MockArticleUsecase) UpdateArticle(ctx context.Context, id uint64, input article.UpdateArticleInput) (*article.UpdateArticleOutput, error) {
	args := m.Called(ctx, id, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*article.UpdateArticleOutput), args.Error(1)
}

func (m *MockArticleUsecase) DeleteArticle(ctx context.Context, id uint64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func ptr[T any](v T) *T {
	return &v
}

func newTestServer(mockUsecase *MockArticleUsecase) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	h := handler.NewArticleHandler(mockUsecase)

	r.POST("/articles", h.CreateArticle)
	r.GET("/articles/:id", h.GetArticleByID)
	r.GET("/articles", h.GetArticles)
	r.PUT("/articles/:id", h.UpdateArticle)
	r.DELETE("/articles/:id", h.DeleteArticle)

	return r
}

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

func TestArticleHandler_CreateArticle(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockUsecase := new(MockArticleUsecase)
		router := newTestServer(mockUsecase)

		input := article.CreateArticleInput{
			Title:  "Test Title",
			Status: "draft",
		}
		expectedArticle := &article.CreateArticleOutput{
			ID:     1,
			Title:  "Test Title",
			Status: "draft",
		}

		mockUsecase.On("CreateArticle", mock.Anything, input).Return(expectedArticle, nil)

		w := sendRequest(router, "POST", "/articles", input)

		assert.Equal(t, http.StatusCreated, w.Code)
		var responseBody article.CreateArticleOutput
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		require.NoError(t, err)
		assert.Equal(t, expectedArticle.ID, responseBody.ID)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("Validation Error", func(t *testing.T) {
		mockUsecase := new(MockArticleUsecase)
		router := newTestServer(mockUsecase)

		input := gin.H{"title": ""} // Invalid: title is required

		w := sendRequest(router, "POST", "/articles", input)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestArticleHandler_GetArticleByID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockUsecase := new(MockArticleUsecase)
		router := newTestServer(mockUsecase)

		expectedArticle := &article.FindArticleByIDOutput{
			ID:    123,
			Title: "Test Title",
		}

		mockUsecase.On("FindArticleByID", mock.Anything, uint64(123)).Return(expectedArticle, nil)

		w := sendRequest(router, "GET", "/articles/123", nil)

		assert.Equal(t, http.StatusOK, w.Code)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		mockUsecase := new(MockArticleUsecase)
		router := newTestServer(mockUsecase)

		mockUsecase.On("FindArticleByID", mock.Anything, uint64(999)).Return(nil, fmt.Errorf("article not found"))

		w := sendRequest(router, "GET", "/articles/999", nil)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockUsecase.AssertExpectations(t)
	})
}

func TestArticleHandler_GetArticles(t *testing.T) {
	t.Run("Success with default parameters", func(t *testing.T) {
		mockUsecase := new(MockArticleUsecase)
		router := newTestServer(mockUsecase)

		expectedInput := article.FindByCriteriaInput{Page: 1, Limit: 10}
		expectedOutput := &article.FindByCriteriaOutput{Articles: []article.FindArticleByIDOutput{}}
		mockUsecase.On("FindByCriteria", mock.Anything, expectedInput).Return(expectedOutput, nil)

		w := sendRequest(router, "GET", "/articles", nil)

		assert.Equal(t, http.StatusOK, w.Code)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("Success with query parameters", func(t *testing.T) {
		mockUsecase := new(MockArticleUsecase)
		router := newTestServer(mockUsecase)

		expectedInput := article.FindByCriteriaInput{Status: ptr("published"), Page: 2, Limit: 20}
		expectedOutput := &article.FindByCriteriaOutput{Articles: []article.FindArticleByIDOutput{}}
		mockUsecase.On("FindByCriteria", mock.Anything, expectedInput).Return(expectedOutput, nil)

		w := sendRequest(router, "GET", "/articles?status=published&page=2&limit=20", nil)

		assert.Equal(t, http.StatusOK, w.Code)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("Validation Error on limit", func(t *testing.T) {
		mockUsecase := new(MockArticleUsecase)
		router := newTestServer(mockUsecase)

		w := sendRequest(router, "GET", "/articles?limit=200", nil)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Validation Error on page", func(t *testing.T) {
		mockUsecase := new(MockArticleUsecase)
		router := newTestServer(mockUsecase)

		w := sendRequest(router, "GET", "/articles?page=-1", nil)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Success with page=0 defaults to page=1", func(t *testing.T) {
		mockUsecase := new(MockArticleUsecase)
		router := newTestServer(mockUsecase)

		expectedInput := article.FindByCriteriaInput{Page: 1, Limit: 10}
		expectedOutput := &article.FindByCriteriaOutput{Articles: []article.FindArticleByIDOutput{}}
		mockUsecase.On("FindByCriteria", mock.Anything, expectedInput).Return(expectedOutput, nil)

		w := sendRequest(router, "GET", "/articles?page=0", nil)

		assert.Equal(t, http.StatusOK, w.Code)
		mockUsecase.AssertExpectations(t)
	})
}

func TestArticleHandler_UpdateArticle(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockUsecase := new(MockArticleUsecase)
		router := newTestServer(mockUsecase)

		input := article.UpdateArticleInput{Title: ptr("Updated Title")}
		expectedArticle := &article.UpdateArticleOutput{ID: 123, Title: "Updated Title"}

		mockUsecase.On("UpdateArticle", mock.Anything, uint64(123), input).Return(expectedArticle, nil)

		w := sendRequest(router, "PUT", "/articles/123", input)

		assert.Equal(t, http.StatusOK, w.Code)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("Validation Error", func(t *testing.T) {
		mockUsecase := new(MockArticleUsecase)
		router := newTestServer(mockUsecase)

		input := gin.H{"status": "invalid_status"}

		w := sendRequest(router, "PUT", "/articles/123", input)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestArticleHandler_DeleteArticle(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockUsecase := new(MockArticleUsecase)
		router := newTestServer(mockUsecase)

		mockUsecase.On("DeleteArticle", mock.Anything, uint64(123)).Return(nil)

		w := sendRequest(router, "DELETE", "/articles/123", nil)

		assert.Equal(t, http.StatusNoContent, w.Code)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		mockUsecase := new(MockArticleUsecase)
		router := newTestServer(mockUsecase)

		mockUsecase.On("DeleteArticle", mock.Anything, uint64(999)).Return(fmt.Errorf("article not found"))

		w := sendRequest(router, "DELETE", "/articles/999", nil)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockUsecase.AssertExpectations(t)
	})
}
