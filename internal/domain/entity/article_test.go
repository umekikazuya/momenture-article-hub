package entity_test

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/umekikazuya/momenture-article-hub/internal/domain/entity"
	"github.com/umekikazuya/momenture-article-hub/internal/domain/vo"
)

// --- ヘルパー関数 ---

func ptr[T any](v T) *T {
	return &v
}

// --- テストケース ---

func TestNewArticle(t *testing.T) {
	t.Parallel()

	t.Run("有効な全属性で作成成功", func(t *testing.T) {
		t.Parallel()
		link := "https://example.com"
		body := "This is the body."
		provider := string(vo.ProviderTypeQiita)

		article, err := entity.NewArticle(
			"Valid Title",
			string(vo.ArticleStatusDraft),
			entity.WithLink(&link),
			entity.WithBody(&body),
			entity.WithProviderType(&provider),
		)

		require.NoError(t, err)
		assert.NotNil(t, article)
		assert.Equal(t, vo.ArticleTitle("Valid Title"), article.Title)
		assert.Equal(t, vo.ArticleStatusDraft, article.Status)
		require.NotNil(t, article.Link)
		assert.Equal(t, link, article.Link.String())
		require.NotNil(t, article.Body)
		assert.Equal(t, body, article.Body.String())
		require.NotNil(t, article.ProviderType)
		assert.Equal(t, vo.ProviderTypeQiita, *article.ProviderType)
		assert.WithinDuration(t, time.Now(), article.CreatedAt, time.Second)
		assert.WithinDuration(t, time.Now(), article.UpdatedAt, time.Second)
		assert.Nil(t, article.DeletedAt)
	})

	t.Run("タイトルが100文字を超える場合はエラー", func(t *testing.T) {
		t.Parallel()
		longTitle := strings.Repeat("a", vo.MaxArticleTitleLength+1)
		_, err := entity.NewArticle(longTitle, string(vo.ArticleStatusDraft))
		assert.Error(t, err)
	})

	// NOTE: bodyはポインタ型でnil許容のため、NewArticleではエラーにならない

	t.Run("無効なステータス値の場合はエラー", func(t *testing.T) {
		t.Parallel()
		_, err := entity.NewArticle("Valid Title", "invalid_status")
		assert.Error(t, err)
	})

	t.Run("無効なプロバイダタイプの場合はエラー", func(t *testing.T) {
		t.Parallel()
		provider := "invalid_provider"
		_, err := entity.NewArticle("Valid Title", string(vo.ArticleStatusDraft), entity.WithProviderType(&provider))
		assert.Error(t, err)
	})
}

func TestArticle_Update(t *testing.T) {
	t.Parallel()

	baseArticle, _ := entity.NewArticle("Original Title", string(vo.ArticleStatusDraft))

	t.Run("有効な内容で更新成功", func(t *testing.T) {
		t.Parallel()
		article := *baseArticle // コピーして使う
		originalUpdatedAt := article.UpdatedAt
		time.Sleep(10 * time.Millisecond) // 更新時刻が変わることを確実にする

		newTitle := "Updated Title"
		newBody := "Updated body."
		newStatus := string(vo.ArticleStatusPublished)

		err := article.Update(&newTitle, &newBody, &newStatus, nil, nil)

		require.NoError(t, err)
		assert.Equal(t, vo.ArticleTitle(newTitle), article.Title)
		require.NotNil(t, article.Body)
		assert.Equal(t, newBody, article.Body.String())
		assert.Equal(t, vo.ArticleStatus(newStatus), article.Status)
		assert.True(t, article.UpdatedAt.After(originalUpdatedAt))
	})

	t.Run("無効なタイトルで更新失敗", func(t *testing.T) {
		t.Parallel()
		article := *baseArticle
		invalidTitle := strings.Repeat("b", vo.MaxArticleTitleLength+1)
		err := article.Update(&invalidTitle, nil, nil, nil, nil)
		assert.Error(t, err)
		assert.Equal(t, baseArticle.Title, article.Title) // 変更されていないこと
	})

	t.Run("無効なステータスで更新失敗", func(t *testing.T) {
		t.Parallel()
		article := *baseArticle
		invalidStatus := "invalid"
		err := article.Update(nil, nil, &invalidStatus, nil, nil)
		assert.Error(t, err)
		assert.Equal(t, baseArticle.Status, article.Status)
	})
}

func TestArticle_Publish(t *testing.T) {
	t.Parallel()
	article, _ := entity.NewArticle("T", string(vo.ArticleStatusDraft))

	t.Run("下書きから公開済みに変更成功", func(t *testing.T) {
		t.Parallel()
		art := *article
		originalUpdatedAt := art.UpdatedAt
		time.Sleep(10 * time.Millisecond)

		err := art.Publish()
		require.NoError(t, err)
		assert.True(t, art.Status.IsPublished())
		assert.True(t, art.UpdatedAt.After(originalUpdatedAt))
	})

	t.Run("既に公開済みの場合はエラー", func(t *testing.T) {
		t.Parallel()
		art, _ := entity.NewArticle("T", string(vo.ArticleStatusPublished))
		err := art.Publish()
		assert.Error(t, err)
	})
}

func TestArticle_Draft(t *testing.T) {
	t.Parallel()
	article, _ := entity.NewArticle("T", string(vo.ArticleStatusPublished))

	t.Run("公開済みから下書きに変更成功", func(t *testing.T) {
		t.Parallel()
		art := *article
		originalUpdatedAt := art.UpdatedAt
		time.Sleep(10 * time.Millisecond)

		err := art.Draft()
		require.NoError(t, err)
		assert.True(t, art.Status.IsDraft())
		assert.True(t, art.UpdatedAt.After(originalUpdatedAt))
	})

	t.Run("既に下書きの場合はエラー", func(t *testing.T) {
		t.Parallel()
		art, _ := entity.NewArticle("T", string(vo.ArticleStatusDraft))
		err := art.Draft()
		assert.Error(t, err)
	})
}

func TestArticle_SoftDelete_And_Restore(t *testing.T) {
	t.Parallel()
	baseArticle, _ := entity.NewArticle("T", string(vo.ArticleStatusDraft))

	t.Run("SoftDelete成功", func(t *testing.T) {
		t.Parallel()
		article := *baseArticle
		originalUpdatedAt := article.UpdatedAt
		time.Sleep(10 * time.Millisecond)

		err := article.SoftDelete()
		require.NoError(t, err)
		assert.NotNil(t, article.DeletedAt)
		assert.True(t, article.UpdatedAt.After(originalUpdatedAt))

		t.Run("既に削除済みの場合はエラー", func(t *testing.T) {
			err := article.SoftDelete()
			assert.Error(t, err)
		})

		t.Run("Restore成功", func(t *testing.T) {
			deletedAtTime := *article.DeletedAt
			originalUpdatedAt := article.UpdatedAt
			time.Sleep(10 * time.Millisecond)

			err := article.Restore()
			require.NoError(t, err)
			assert.Nil(t, article.DeletedAt)
			assert.True(t, article.UpdatedAt.After(originalUpdatedAt))
			assert.True(t, article.UpdatedAt.After(deletedAtTime))

			t.Run("まだ削除されていない場合はエラー", func(t *testing.T) {
				err := article.Restore()
				assert.Error(t, err)
			})
		})
	})
}

func TestArticle_ChangeProvider(t *testing.T) {
	t.Parallel()

	t.Run("下書き記事のプロバイダ変更成功", func(t *testing.T) {
		t.Parallel()
		provider := string(vo.ProviderTypeZenn)
		article, _ := entity.NewArticle("T", string(vo.ArticleStatusDraft), entity.WithProviderType(&provider))
		originalUpdatedAt := article.UpdatedAt
		time.Sleep(10 * time.Millisecond)

		err := article.ChangeProvider(vo.ProviderTypeQiita)

		require.NoError(t, err)
		require.NotNil(t, article.ProviderType)
		assert.Equal(t, vo.ProviderTypeQiita, *article.ProviderType)
		assert.True(t, article.UpdatedAt.After(originalUpdatedAt))
	})

	t.Run("公開済み記事のプロバイダ変更はエラー", func(t *testing.T) {
		t.Parallel()
		article, _ := entity.NewArticle("T", string(vo.ArticleStatusPublished))
		err := article.ChangeProvider(vo.ProviderTypeQiita)
		assert.Error(t, err)
	})
}
