package main

import (
	"log"

	"github.com/umekikazuya/momenture-article-hub/internal/config"
	"github.com/umekikazuya/momenture-article-hub/internal/infrastructure/persistence"
	"github.com/umekikazuya/momenture-article-hub/internal/infrastructure/persistence/postgres"
	"github.com/umekikazuya/momenture-article-hub/internal/interfaces/handler"
	"github.com/umekikazuya/momenture-article-hub/internal/interfaces/router"
	"github.com/umekikazuya/momenture-article-hub/internal/usecase/article"
)

func main() {
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	db, err := persistence.NewPostgreSQLDB(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get underlying sql.DB: %v", err)
	}
	defer sqlDB.Close()

	articleRepository := postgres.NewPostgresArticleRepository(db)

	articleUsecase := article.NewArticleUsecase(articleRepository)

	articleHandler := handler.NewArticleHandler(articleUsecase)

	routerConfig := &router.RouterConfig{
		ArticleHandler: articleHandler,
	}

	r := router.NewRouter(routerConfig)

	log.Fatal(r.Run(":8080"))
}
