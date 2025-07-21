package vo_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/umekikazuya/momenture-article-hub/internal/domain/vo"
)

func TestNewArticleTitle(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		value     string
		want      vo.ArticleTitle
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "有効なタイトルで作成成功",
			value:     "Valid Title",
			want:      vo.ArticleTitle("Valid Title"),
			assertion: assert.NoError,
		},
		{
			name:      "ちょうど100文字のタイトルで作成成功",
			value:     strings.Repeat("a", vo.MaxArticleTitleLength),
			want:      vo.ArticleTitle(strings.Repeat("a", vo.MaxArticleTitleLength)),
			assertion: assert.NoError,
		},
		{
			name:      "空文字列の場合はエラー",
			value:     "",
			want:      "",
			assertion: assert.Error,
		},
		{
			name:      "100文字を超える場合はエラー",
			value:     strings.Repeat("a", vo.MaxArticleTitleLength+1),
			want:      "",
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := vo.NewArticleTitle(tt.value)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestArticleTitle_String(t *testing.T) {
	t.Parallel()

	titleValue := "My Test Title"
	title, err := vo.NewArticleTitle(titleValue)
	require.NoError(t, err)

	assert.Equal(t, titleValue, title.String())
}
