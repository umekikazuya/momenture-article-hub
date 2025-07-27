package postgres

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/umekikazuya/momenture-article-hub/internal/config"
	"github.com/umekikazuya/momenture-article-hub/internal/domain/entity"
	"github.com/umekikazuya/momenture-article-hub/internal/domain/repository"
	"github.com/umekikazuya/momenture-article-hub/internal/infrastructure/persistence"
	"gorm.io/gorm"
)

var testDB *gorm.DB

func TestMain(m *testing.M) {
	cfg := &config.DatabaseConfig{
		Host:     os.Getenv("POSTGRES_HOST_TEST"),
		Port:     "5432",
		User:     os.Getenv("POSTGRES_USER_TEST"),
		Password: os.Getenv("POSTGRES_PASSWORD_TEST"),
		Name:     os.Getenv("POSTGRES_DB_TEST"),
	}
	var err error
	testDB, err = persistence.NewPostgreSQLDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}

	// テスト実行
	code := m.Run()

	// テストスイート終了後にDB接続を閉じる
	sqlDB, _ := testDB.DB()
	sqlDB.Close()
	os.Exit(code)
}

// テスト用ヘルパー関数
func stringPtr(s string) *string {
	return &s
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func setupTestRepository(t *testing.T) *PostgresArticleRepository {
	// 各テストの前にトランザクションを開始
	tx := testDB.Begin()

	// テスト終了時に必ずロールバックするように設定
	t.Cleanup(func() {
		tx.Rollback()
	})
	return &PostgresArticleRepository{db: tx}
}

func TestPostgresArticleRepository_Create_Success_MinimalFields(t *testing.T) {
	// Given: リポジトリインスタンスと必須フィールドのみの記事
	repo := setupTestRepository(t)
	require.NotNil(t, repo, "Repository should be initialized")

	// 必須フィールドのみで記事を作成
	article, err := entity.NewArticle(
		"Test Article Title",
		"draft",
	)
	require.NoError(t, err)
	require.NotNil(t, article)

	// When: Create メソッドを呼び出す
	ctx := context.Background()
	createdArticle, err := repo.Create(ctx, article)

	// Then: エラーがなく、正しく作成されること
	require.NoError(t, err)
	require.NotNil(t, createdArticle)

	// IDが割り当てられていること
	assert.NotZero(t, createdArticle.ID)

	// CreatedAtとUpdatedAtが設定されていること
	assert.False(t, createdArticle.CreatedAt.IsZero())
	assert.False(t, createdArticle.UpdatedAt.IsZero())

	// 必須フィールドが正しく設定されていること
	assert.Equal(t, "Test Article Title", createdArticle.Title.String())
	assert.True(t, createdArticle.Status.IsDraft())

	// オプションフィールドがnilであること
	assert.Nil(t, createdArticle.Body)
	assert.Nil(t, createdArticle.ProviderType)
	assert.Nil(t, createdArticle.Link)
	assert.Nil(t, createdArticle.DeletedAt)
}

func TestPostgresArticleRepository_Create_Success_AllFields(t *testing.T) {
	// Given: リポジトリインスタンスと全てのフィールドが設定された記事
	repo := setupTestRepository(t)
	require.NotNil(t, repo, "Repository should be initialized")

	// 全てのオプションフィールドを含む記事を作成
	bodyMarkdown := "# Test Article\nThis is test content."
	providerType := "qiita"
	link := "https://qiita.com/test/test-article"

	article, err := entity.NewArticle(
		"Test Article with All Fields",
		"published",
		entity.WithBody(&bodyMarkdown),
		entity.WithProviderType(&providerType),
		entity.WithLink(&link),
	)
	require.NoError(t, err)
	require.NotNil(t, article)

	// When: Create メソッドを呼び出す
	ctx := context.Background()
	createdArticle, err := repo.Create(ctx, article)

	// Then: エラーがなく、正しく作成されること
	require.NoError(t, err)
	require.NotNil(t, createdArticle)

	// IDが割り当てられていること
	assert.NotZero(t, createdArticle.ID)

	// 全てのフィールドが正しく設定されていること
	assert.Equal(t, "Test Article with All Fields", createdArticle.Title.String())
	assert.Equal(t, "published", string(createdArticle.Status))
	assert.NotNil(t, createdArticle.Body)
	assert.Equal(t, bodyMarkdown, createdArticle.Body.String())
	assert.NotNil(t, createdArticle.ProviderType)
	assert.Equal(t, providerType, createdArticle.ProviderType.String())
	assert.NotNil(t, createdArticle.Link)
	assert.Equal(t, link, createdArticle.Link.String())
}

func TestPostgresArticleRepository_FindByID_Success(t *testing.T) {
	// Given: リポジトリインスタンスと事前に保存された記事
	repo := setupTestRepository(t)
	require.NotNil(t, repo, "Repository should be initialized")

	// 事前に記事を作成
	originalArticle, err := entity.NewArticle(
		"Test Article for FindByID",
		"draft",
	)
	require.NoError(t, err)

	ctx := context.Background()
	createdArticle, err := repo.Create(ctx, originalArticle)
	require.NoError(t, err)
	require.NotNil(t, createdArticle)

	// When: FindByID メソッドを呼び出す
	foundArticle, err := repo.FindByID(ctx, createdArticle.ID)

	// Then: エラーがなく、正しく取得されること
	require.NoError(t, err)
	require.NotNil(t, foundArticle)

	// 取得された記事の内容が正しいこと
	assert.Equal(t, createdArticle.ID, foundArticle.ID)
	assert.Equal(t, createdArticle.Title.String(), foundArticle.Title.String())
	assert.Equal(t, createdArticle.Status, foundArticle.Status)
}

func TestPostgresArticleRepository_FindByID_NotFound(t *testing.T) {
	// Given: リポジトリインスタンスと存在しないID
	repo := setupTestRepository(t)
	require.NotNil(t, repo, "Repository should be initialized")

	// When: 存在しないIDでFindByIDを呼び出す
	ctx := context.Background()
	foundArticle, err := repo.FindByID(ctx, 9999)

	// Then: 記事はnilであること
	assert.NoError(t, err)
	assert.Nil(t, foundArticle)
}

func TestPostgresArticleRepository_Update_Success(t *testing.T) {
	// Given: リポジトリインスタンスと事前に保存された記事
	repo := setupTestRepository(t)
	require.NotNil(t, repo, "Repository should be initialized")

	// 事前に記事を作成
	originalArticle, err := entity.NewArticle(
		"Original Title",
		"draft",
	)
	require.NoError(t, err)

	ctx := context.Background()
	createdArticle, err := repo.Create(ctx, originalArticle)
	require.NoError(t, err)
	require.NotNil(t, createdArticle)

	originalUpdatedAt := createdArticle.UpdatedAt

	// 少し時間を待ってからUpdateを実行（UpdatedAtの変更を確認するため）
	time.Sleep(10 * time.Millisecond)

	// 記事を更新（ここではタイトルを変更）
	updatedTitle := "Updated Title"
	require.NoError(t, createdArticle.Update(
		stringPtr(updatedTitle),
		nil,
		stringPtr("published"),
		nil, nil,
	))

	// When: Update メソッドを呼び出す
	err = repo.Update(ctx, createdArticle)

	// Then: エラーがないこと
	require.NoError(t, err)

	// データベースから再取得して確認
	foundArticle, err := repo.FindByID(ctx, createdArticle.ID)
	require.NoError(t, err)
	require.NotNil(t, foundArticle)

	// 更新された内容が正しいこと
	assert.Equal(t, updatedTitle, foundArticle.Title.String())
	assert.Equal(t, "published", string(foundArticle.Status))

	// UpdatedAtが更新されていること
	assert.True(t, foundArticle.UpdatedAt.After(originalUpdatedAt))
}

func TestPostgresArticleRepository_Update_NotFound(t *testing.T) {
	// Given: リポジトリインスタンスと存在しないIDの記事
	repo := setupTestRepository(t)
	require.NotNil(t, repo, "Repository should be initialized")

	// 存在しないIDの記事
	nonExistentArticle, err := entity.NewArticle(
		"Non-existent Article",
		"draft",
	)
	require.NoError(t, err)
	nonExistentArticle.ID = 9999 // 存在しないID

	// When: Update メソッドを呼び出す
	ctx := context.Background()
	err = repo.Update(ctx, nonExistentArticle)

	// Then: エラーが返されること
	assert.Error(t, err)
}

func TestPostgresArticleRepository_Delete_Success(t *testing.T) {
	// Given: リポジトリインスタンスと事前に保存された記事
	repo := setupTestRepository(t)
	require.NotNil(t, repo, "Repository should be initialized")

	// 事前に記事を作成
	originalArticle, err := entity.NewArticle(
		"Article to Delete",
		"draft",
	)
	require.NoError(t, err)

	ctx := context.Background()
	createdArticle, err := repo.Create(ctx, originalArticle)
	require.NoError(t, err)
	require.NotNil(t, createdArticle)

	// When: Delete メソッドを呼び出す
	err = repo.Delete(ctx, createdArticle.ID)

	// Then: エラーがないこと
	require.NoError(t, err)

	// 削除後にFindByIDで取得できることを確認
	foundArticle, err := repo.FindByID(ctx, createdArticle.ID)
	require.NoError(t, err)
	require.NotNil(t, foundArticle)
	// IDの一致を確認
	assert.Equal(t, foundArticle.ID, createdArticle.ID)
	// DeletedAtが設定されていることを確認
	assert.NotNil(t, foundArticle.DeletedAt)
	// DeletedAtが現在時刻よりも前であることを確認
	assert.True(t, foundArticle.DeletedAt.Before(time.Now()))
}

func TestPostgresArticleRepository_Delete_NotFound(t *testing.T) {
	// Given: リポジトリインスタンスと存在しないID
	repo := setupTestRepository(t)
	require.NotNil(t, repo, "Repository should be initialized")

	// When: 存在しないIDでDeleteを呼び出す
	ctx := context.Background()
	err := repo.Delete(ctx, 9999)

	// Then: エラーが返されること
	assert.Error(t, err)
}

func TestPostgresArticleRepository_FindByCriteria_DefaultConditions(t *testing.T) {
	// Given: リポジトリインスタンスと複数の記事
	repo := setupTestRepository(t)
	require.NotNil(t, repo, "Repository should be initialized")

	ctx := context.Background()

	// 複数の記事を作成
	articles := []*entity.Article{}
	for i := 0; i < 3; i++ {
		article, err := entity.NewArticle(
			"Test Article "+string(rune('A'+i)),
			"draft",
		)
		require.NoError(t, err)

		createdArticle, err := repo.Create(ctx, article)
		require.NoError(t, err)
		articles = append(articles, createdArticle)
	}

	// When: デフォルト条件でFindByCriteriaを呼び出す
	criteria := repository.ArticleQueryCriteria{}
	foundArticles, total, err := repo.FindByCriteria(ctx, criteria)

	// Then: エラーがなく、全ての記事が取得されること
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(foundArticles), 3)
	assert.GreaterOrEqual(t, total, 3)
}

func TestPostgresArticleRepository_FindByCriteria_FilterByStatus(t *testing.T) {
	// Given: リポジトリインスタンスと異なるステータスの記事
	repo := setupTestRepository(t)
	require.NotNil(t, repo, "Repository should be initialized")

	ctx := context.Background()

	// ドラフトの記事を作成
	draftArticle, err := entity.NewArticle(
		"Draft Article",
		"draft",
	)
	require.NoError(t, err)
	_, err = repo.Create(ctx, draftArticle)
	require.NoError(t, err)

	// 公開済みの記事を作成
	publishedArticle, err := entity.NewArticle(
		"Published Article",
		"published",
	)
	require.NoError(t, err)
	_, err = repo.Create(ctx, publishedArticle)
	require.NoError(t, err)

	// When: publishedステータスでフィルタリング
	publishedStatus := "published"
	criteria := repository.ArticleQueryCriteria{
		Status: &publishedStatus,
	}
	foundArticles, total, err := repo.FindByCriteria(ctx, criteria)

	// Then: エラーがなく、公開済みの記事のみが取得されること
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, 1)

	// 全ての取得された記事がpublishedステータスであること
	for _, article := range foundArticles {
		assert.Equal(t, "published", string(article.Status))
	}
}

func TestPostgresArticleRepository_FindByCriteria_Pagination(t *testing.T) {
	// Given: リポジトリインスタンスと多数の記事
	repo := setupTestRepository(t)
	require.NotNil(t, repo, "Repository should be initialized")

	ctx := context.Background()

	// 10個の記事を作成
	for i := 0; i < 10; i++ {
		article, err := entity.NewArticle(
			"Test Article "+string(rune('A'+i)),
			"draft",
		)
		require.NoError(t, err)

		_, err = repo.Create(ctx, article)
		require.NoError(t, err)
	}

	// When: ページネーション条件でFindByCriteriaを呼び出す
	criteria := repository.ArticleQueryCriteria{
		Page:  2,
		Limit: 3,
	}
	foundArticles, total, err := repo.FindByCriteria(ctx, criteria)

	// Then: エラーがなく、指定された件数が取得されること
	require.NoError(t, err)
	assert.LessOrEqual(t, len(foundArticles), 3)
	assert.GreaterOrEqual(t, total, 10)
}

func TestPostgresArticleRepository_FindByCriteria_EmptyResult(t *testing.T) {
	// Given: リポジトリインスタンスと条件に合わない記事のみ
	repo := setupTestRepository(t)
	require.NotNil(t, repo, "Repository should be initialized")

	// When: 存在しない条件でFindByCriteriaを呼び出す
	nonExistentStatus := "invalid_status"
	criteria := repository.ArticleQueryCriteria{
		Status: &nonExistentStatus,
	}
	ctx := context.Background()
	foundArticles, total, err := repo.FindByCriteria(ctx, criteria)

	// Then: エラーがなく、空の結果が返されること
	require.NoError(t, err)
	assert.Empty(t, foundArticles)
	assert.Equal(t, 0, total)
}

func TestPostgresArticleRepository_Update_OptionalFieldsSet(t *testing.T) {
	// Given: リポジトリインスタンスと事前に保存された最小限の記事
	repo := setupTestRepository(t)
	require.NotNil(t, repo, "Repository should be initialized")

	// 最小限の記事を作成
	originalArticle, err := entity.NewArticle(
		"Original Article",
		"draft",
	)
	require.NoError(t, err)

	ctx := context.Background()
	createdArticle, err := repo.Create(ctx, originalArticle)
	require.NoError(t, err)
	require.NotNil(t, createdArticle)

	// オプションフィールドを追加した記事を作成
	bodyMarkdown := "# Updated Content\nThis is updated content."
	providerType := "qiita"
	link := "https://qiita.com/updated/article"

	updatedArticle, err := entity.NewArticle(
		"Updated Article with Options",
		"published",
		entity.WithBody(&bodyMarkdown),
		entity.WithProviderType(&providerType),
		entity.WithLink(&link),
	)
	require.NoError(t, err)
	updatedArticle.ID = createdArticle.ID

	// When: Update メソッドを呼び出す
	err = repo.Update(ctx, updatedArticle)

	// Then: エラーがないこと
	require.NoError(t, err)

	// データベースから再取得して確認
	foundArticle, err := repo.FindByID(ctx, createdArticle.ID)
	require.NoError(t, err)
	require.NotNil(t, foundArticle)

	// オプションフィールドが正しく設定されていること
	assert.NotNil(t, foundArticle.Body)
	assert.Equal(t, bodyMarkdown, foundArticle.Body.String())
	assert.NotNil(t, foundArticle.ProviderType)
	assert.Equal(t, providerType, foundArticle.ProviderType.String())
	assert.NotNil(t, foundArticle.Link)
	assert.Equal(t, link, foundArticle.Link.String())
}

func TestPostgresArticleRepository_Update_ChangeTitle(t *testing.T) {
	// Given: リポジトリインスタンスと事前に保存された最小限の記事
	repo := setupTestRepository(t)
	require.NotNil(t, repo, "Repository should be initialized")

	// 最小限の記事を作成
	originalArticle, err := entity.NewArticle(
		"Original Article",
		"draft",
	)
	require.NoError(t, err)

	ctx := context.Background()
	createdArticle, err := repo.Create(ctx, originalArticle)
	require.NoError(t, err)
	require.NotNil(t, createdArticle)

	// オプションフィールドを追加した記事を作成
	newTitle := "Updated Article Title"

	updatedArticle, err := entity.NewArticle(
		newTitle,
		"draft",
	)
	require.NoError(t, err)
	updatedArticle.ID = createdArticle.ID

	// When: Update メソッドを呼び出す
	err = repo.Update(ctx, updatedArticle)

	// Then: エラーがないこと
	require.NoError(t, err)

	// データベースから再取得して確認
	foundArticle, err := repo.FindByID(ctx, createdArticle.ID)
	require.NoError(t, err)
	require.NotNil(t, foundArticle)

	// タイトルが更新されていること
	assert.Equal(t, newTitle, foundArticle.Title.String())
	// 他のフィールドは変更されていないことを確認
	assert.Equal(t, createdArticle.Body, foundArticle.Body)
	assert.Equal(t, createdArticle.Status, foundArticle.Status)
	assert.Equal(t, createdArticle.ProviderType, foundArticle.ProviderType)
	assert.Equal(t, createdArticle.Link, foundArticle.Link)
	// UpdatedAtが更新されていることを確認
	assert.True(t, foundArticle.UpdatedAt.After(createdArticle.UpdatedAt), "UpdatedAt should be updated")
}

func TestPostgresArticleRepository_Update_OptionalFieldsCleared(t *testing.T) {
	// Given: リポジトリインスタンスと全フィールドが設定された記事
	repo := setupTestRepository(t)
	require.NotNil(t, repo, "Repository should be initialized")

	// 全フィールドが設定された記事を作成
	bodyMarkdown := "# Original Content"
	providerType := "zenn"
	link := "https://zenn.dev/original/article"

	originalArticle, err := entity.NewArticle(
		"Article with All Fields",
		"published",
		entity.WithBody(&bodyMarkdown),
		entity.WithProviderType(&providerType),
		entity.WithLink(&link),
	)
	require.NoError(t, err, "Failed to create original article")

	ctx := context.Background()
	createdArticle, err := repo.Create(ctx, originalArticle)
	require.NoError(t, err, "Failed to create article")
	require.NotNil(t, createdArticle, "Created article should not be nil")

	// オプションフィールドをクリアした記事を作成
	updatedArticle, err := entity.NewArticle(
		"Updated Article Cleared",
		"draft",
		entity.WithBody(stringPtr("")),
		entity.WithProviderType(stringPtr("")),
		entity.WithLink(stringPtr("")),
	)
	require.NoError(t, err)
	updatedArticle.ID = createdArticle.ID

	// When: Update メソッドを呼び出す
	err = repo.Update(ctx, updatedArticle)

	// Then: エラーがないこと
	require.NoError(t, err, "Update should not return an error")

	// データベースから再取得して確認
	foundArticle, err := repo.FindByID(ctx, createdArticle.ID)
	require.NoError(t, err, "Failed to find article by ID", err)
	require.NotNil(t, foundArticle)

	// オプションフィールドがnilになっていること
	// デバッグ
	assert.Nil(t, *foundArticle.Body, "Body should be nil after clearing")
	assert.Nil(t, *foundArticle.ProviderType)
	assert.Nil(t, *foundArticle.Link)
}

func TestPostgresArticleRepository_Delete_LogicalDeletion(t *testing.T) {
	// Given: リポジトリインスタンスと事前に保存された記事
	repo := setupTestRepository(t)
	require.NotNil(t, repo, "Repository should be initialized")

	// 事前に記事を作成
	originalArticle, err := entity.NewArticle(
		"Article to Delete",
		"published",
	)
	require.NoError(t, err)

	ctx := context.Background()
	createdArticle, err := repo.Create(ctx, originalArticle)
	require.NoError(t, err)
	require.NotNil(t, createdArticle)
	require.Nil(t, createdArticle.DeletedAt)

	// When: Delete メソッドを呼び出す
	err = repo.Delete(ctx, createdArticle.ID)

	// Then: エラーがないこと
	require.NoError(t, err)

	// 論理削除後にFindByIDで取得できないこと
	deletedArticle, err := repo.FindByID(ctx, createdArticle.ID)
	require.NoError(t, err)
	assert.NotNil(t, deletedArticle)
	// IDの一致を確認
	assert.Equal(t, deletedArticle.ID, createdArticle.ID)
	// DeletedAtが設定されていることを確認
	assert.NotNil(t, deletedArticle.DeletedAt)
	// DeletedAtが現在時刻よりも前であることを確認
	assert.True(t, deletedArticle.DeletedAt.Before(time.Now()))
}

func TestPostgresArticleRepository_FindByCriteria_FilterByProviderType(t *testing.T) {
	// Given: リポジトリインスタンスと異なるプロバイダタイプの記事
	repo := setupTestRepository(t)
	require.NotNil(t, repo, "Repository should be initialized")

	ctx := context.Background()

	// Qiitaの記事を作成
	qiitaProvider := "qiita"
	qiitaLink := "https://qiita.com/user/article1"
	qiitaArticle, err := entity.NewArticle(
		"Qiita Article",
		"published",
		entity.WithProviderType(&qiitaProvider),
		entity.WithLink(&qiitaLink),
	)
	require.NoError(t, err)
	_, err = repo.Create(ctx, qiitaArticle)
	require.NoError(t, err)

	// Zennの記事を作成
	zennProvider := "zenn"
	zennLink := "https://zenn.dev/user/article1"
	zennArticle, err := entity.NewArticle(
		"Zenn Article",
		"published",
		entity.WithProviderType(&zennProvider),
		entity.WithLink(&zennLink),
	)
	require.NoError(t, err)
	_, err = repo.Create(ctx, zennArticle)
	require.NoError(t, err)

	// When: qiitaプロバイダタイプでフィルタリング
	criteria := repository.ArticleQueryCriteria{
		ProviderType: &qiitaProvider,
	}
	foundArticles, total, err := repo.FindByCriteria(ctx, criteria)

	// Then: エラーがなく、Qiitaの記事のみが取得されること
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, 1)

	// 全ての取得された記事がqiitaプロバイダタイプであること
	for _, article := range foundArticles {
		assert.NotNil(t, article.ProviderType)
		assert.Equal(t, "qiita", article.ProviderType.String())
	}
}

func TestPostgresArticleRepository_FindByCriteria_SortByUpdatedAt(t *testing.T) {
	// Given: リポジトリインスタンスと異なる時刻に作成された記事
	repo := setupTestRepository(t)
	require.NotNil(t, repo, "Repository should be initialized")

	ctx := context.Background()

	// 複数の記事を時間差で作成
	articleIds := []uint64{}
	for i := 0; i < 3; i++ {
		article, err := entity.NewArticle(
			"Test Article "+string(rune('A'+i)),
			"published",
		)
		require.NoError(t, err)

		createdArticle, err := repo.Create(ctx, article)
		require.NoError(t, err)
		articleIds = append(articleIds, createdArticle.ID)

		// 時間差を作るため少し待機
		time.Sleep(10 * time.Millisecond)
	}

	// When: updated_atでDESC(降順)ソート
	sortBy := "updated_at"
	sortOrder := "desc"
	criteria := repository.ArticleQueryCriteria{
		SortBy:    &sortBy,
		SortOrder: &sortOrder,
	}
	foundArticles, _, err := repo.FindByCriteria(ctx, criteria)

	// Then: エラーがなく、降順で記事が取得されること
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(foundArticles), 3)

	// 記事のUpdatedAtが降順になっていることを確認
	for i := 0; i < len(foundArticles)-1; i++ {
		assert.True(
			t, foundArticles[i].UpdatedAt.After(foundArticles[i+1].UpdatedAt),
			"Article %d should be newer than Article %d", i, i+1,
		)
	}
}

func TestPostgresArticleRepository_FindByCriteria_IncludeDeleted(t *testing.T) {
	// Given: リポジトリインスタンスと削除された記事を含む記事
	repo := setupTestRepository(t)
	require.NotNil(t, repo, "Repository should be initialized")

	ctx := context.Background()

	// 通常の記事を作成
	normalArticle, err := entity.NewArticle(
		"Normal Article",
		"published",
	)
	require.NoError(t, err)
	createdNormal, err := repo.Create(ctx, normalArticle)
	require.NoError(t, err)

	// 削除する記事を作成
	deletableArticle, err := entity.NewArticle(
		"Article to Delete",
		"draft",
	)
	require.NoError(t, err)
	createdDeletable, err := repo.Create(ctx, deletableArticle)
	require.NoError(t, err)

	// 記事を削除
	err = repo.Delete(ctx, createdDeletable.ID)
	require.NoError(t, err)

	// When: IncludeDeleted: trueで検索
	criteria := repository.ArticleQueryCriteria{
		IncludeDeleted: true,
	}
	foundArticles, total, err := repo.FindByCriteria(ctx, criteria)

	// Then: エラーがなく、削除された記事も含まれること
	require.NoError(t, err)
	assert.GreaterOrEqual(t, total, 2)

	// 削除された記事と通常の記事の両方が含まれることを確認
	foundNormal := false
	foundDeleted := false
	for _, article := range foundArticles {
		if article.ID == createdNormal.ID {
			foundNormal = true
			assert.Nil(t, article.DeletedAt)
		}
		if article.ID == createdDeletable.ID {
			foundDeleted = true
			assert.NotNil(t, article.DeletedAt)
		}
	}
	assert.True(t, foundNormal, "Normal article should be found")
	assert.True(t, foundDeleted, "Deleted article should be found when IncludeDeleted is true")

	// When: IncludeDeleted: falseで検索（デフォルト）
	criteria = repository.ArticleQueryCriteria{
		IncludeDeleted: false,
	}
	foundArticles, total, err = repo.FindByCriteria(ctx, criteria)

	// Then: 削除された記事は含まれないこと
	require.NoError(t, err)
	foundDeletedInNormalSearch := false
	for _, article := range foundArticles {
		if article.ID == createdDeletable.ID {
			foundDeletedInNormalSearch = true
		}
	}
	assert.False(t, foundDeletedInNormalSearch, "Deleted article should not be found when IncludeDeleted is false")
}

func TestPostgresArticleRepository_FindAll_Success(t *testing.T) {
	// Given: リポジトリインスタンスと複数の記事
	repo := setupTestRepository(t)
	require.NotNil(t, repo, "Repository should be initialized")

	ctx := context.Background()

	// 複数の記事を作成
	expectedCount := 5
	for i := 0; i < expectedCount; i++ {
		article, err := entity.NewArticle(
			"Test Article "+string(rune('A'+i)),
			"published",
		)
		require.NoError(t, err)

		_, err = repo.Create(ctx, article)
		require.NoError(t, err)
	}

	// When: FindAll メソッドを呼び出す
	foundArticles, err := repo.FindAll(ctx)

	// Then: エラーがなく、全ての記事が取得されること
	require.NoError(t, err)
	assert.Equal(t, expectedCount, len(foundArticles))

	// 削除されていない記事のみが含まれていること
	for _, article := range foundArticles {
		assert.Nil(t, article.DeletedAt)
	}
}
