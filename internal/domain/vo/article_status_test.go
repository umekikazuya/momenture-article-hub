package vo_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/umekikazuya/momenture-article-hub/internal/domain/vo"
)

func TestArticleStatus_IsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		as   vo.ArticleStatus
		want bool
	}{
		{name: "Draftは有効", as: vo.ArticleStatusDraft, want: true},
		{name: "Publishedは有効", as: vo.ArticleStatusPublished, want: true},
		{name: "無効な値はfalse", as: vo.ArticleStatus("invalid_status"), want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, tt.as.IsValid())
		})
	}
}

func TestArticleStatus_IsDraft(t *testing.T) {
	t.Parallel()

	assert.True(t, vo.ArticleStatusDraft.IsDraft())
	assert.False(t, vo.ArticleStatusPublished.IsDraft())
}

func TestArticleStatus_IsPublished(t *testing.T) {
	t.Parallel()

	assert.True(t, vo.ArticleStatusPublished.IsPublished())
	assert.False(t, vo.ArticleStatusDraft.IsPublished())
}

func TestArticleStatus_String(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "draft", vo.ArticleStatusDraft.String())
	assert.Equal(t, "published", vo.ArticleStatusPublished.String())
}
