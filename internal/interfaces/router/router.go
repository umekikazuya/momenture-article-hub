package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/umekikazuya/momenture-article-hub/internal/interfaces/handler"
)

// RouterConfig はルーターを設定するための依存関係を保持。
type RouterConfig struct {
	ArticleHandler *handler.ArticleHandler
}

// NewRouter は新しいGinルーターインスタンスを生成し、ルートを設定。
func NewRouter(cfg *RouterConfig) *gin.Engine {
	r := gin.Default()

	r.GET("/up", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	v1 := r.Group("/v1")
	{
		v1.POST("/articles", cfg.ArticleHandler.CreateArticle)
		v1.GET("/articles", cfg.ArticleHandler.GetArticles)
		v1.GET("/articles/:id", cfg.ArticleHandler.FindArticleByID)
		v1.PUT("/articles/:id", cfg.ArticleHandler.UpdateArticle)
		v1.DELETE("/articles/:id", cfg.ArticleHandler.DeleteArticle)
	}

	return r
}
