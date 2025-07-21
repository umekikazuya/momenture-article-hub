package vo_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/umekikazuya/momenture-article-hub/internal/domain/vo"
)

func TestNewArticleBody(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		value     *string
		want      *vo.ArticleBody
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:  "有効な本文で作成成功",
			value: func() *string { s := "This is a valid body."; return &s }(),
			want: func() *vo.ArticleBody {
				b, _ := vo.NewArticleBody(func() *string { s := "This is a valid body."; return &s }())
				return b
			}(),
			assertion: assert.NoError,
		},
		{
			name:      "空文字列の本文で作成成功",
			value:     func() *string { s := ""; return &s }(),
			want:      func() *vo.ArticleBody { b, _ := vo.NewArticleBody(func() *string { s := ""; return &s }()); return b }(),
			assertion: assert.NoError,
		},
		{
			name:      "nilの場合も作成成功 (nilが返る)",
			value:     nil,
			want:      nil,
			assertion: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := vo.NewArticleBody(tt.value)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestArticleBody_String(t *testing.T) {
	t.Parallel()

	t.Run("有効な本文の場合", func(t *testing.T) {
		t.Parallel()
		bodyValue := "My Test Body"
		body, err := vo.NewArticleBody(&bodyValue)
		require.NoError(t, err)
		assert.Equal(t, bodyValue, body.String())
	})

	t.Run("nilの本文の場合", func(t *testing.T) {
		t.Parallel()
		var body *vo.ArticleBody
		assert.Equal(t, "", body.String())
	})
}
