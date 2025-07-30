package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/umekikazuya/momenture-article-hub/internal/usecase/article"
)

// ArticleUsecase defines the interface for article-related use cases.
// This is used to allow mocking in tests.
type ArticleUsecase interface {
	CreateArticle(ctx context.Context, input article.CreateArticleInput) (*article.CreateArticleOutput, error)
	FindArticleByID(ctx context.Context, id uint64) (*article.FindArticleByIDOutput, error)
	FindByCriteria(ctx context.Context, criteria article.FindByCriteriaInput) (*article.FindByCriteriaOutput, error)
	UpdateArticle(ctx context.Context, id uint64, input article.UpdateArticleInput) (*article.UpdateArticleOutput, error)
	DeleteArticle(ctx context.Context, id uint64) error
}

// ArticleHandler handles HTTP requests related to articles.
type ArticleHandler struct {
	usecase  ArticleUsecase
	validate *validator.Validate
}

// NewArticleHandler creates a new ArticleHandler.
func NewArticleHandler(usecase ArticleUsecase) *ArticleHandler {
	return &ArticleHandler{
		usecase:  usecase,
		validate: validator.New(),
	}
}

// CreateArticle handles the creation of a new article.
func (h *ArticleHandler) CreateArticle(c *gin.Context) {
	var input article.CreateArticleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	articleData, err := h.usecase.CreateArticle(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, articleData)
}

// GetArticleByID handles fetching a single article by its ID.
func (h *ArticleHandler) GetArticleByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	articleData, err := h.usecase.FindArticleByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "article not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "article not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// Explicitly construct the response to match the test expectations.
	c.JSON(http.StatusOK, gin.H{
		"ID":        articleData.ID,
		"Title":     articleData.Title,
		"Body":      articleData.Body,
		"Status":    articleData.Status,
		"CreatedAt": articleData.CreatedAt,
		"UpdatedAt": articleData.UpdatedAt,
	})
}

// GetArticles handles fetching a list of articles based on query parameters.
func (h *ArticleHandler) GetArticles(c *gin.Context) {
	var input article.FindByCriteriaInput
	if err := c.ShouldBindQuery(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
		return
	}

	// Set defaults for pagination
	if input.Page == 0 {
		input.Page = 1
	}
	if input.Limit == 0 {
		input.Limit = 10
	}

	if err := h.validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	articles, err := h.usecase.FindByCriteria(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, articles)
}

// UpdateArticle handles updating an existing article.
func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	var input article.UpdateArticleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedArticle, err := h.usecase.UpdateArticle(c.Request.Context(), id, input)
	if err != nil {
		if err.Error() == "article not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "article not found"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedArticle)
}

// DeleteArticle handles soft-deleting an article.
func (h *ArticleHandler) DeleteArticle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	err = h.usecase.DeleteArticle(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "article not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "article not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.Status(http.StatusNoContent)
}
