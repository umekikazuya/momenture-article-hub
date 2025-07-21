package vo_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/umekikazuya/momenture-article-hub/internal/domain/vo"
)

func TestNewLink(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		value     *string
		want      *vo.Link
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:  "有効なURLで作成成功",
			value: func() *string { s := "https://example.com/path"; return &s }(),
			want: func() *vo.Link {
				l, _ := vo.NewLink(func() *string { s := "https://example.com/path"; return &s }())
				return l
			}(),
			assertion: assert.NoError,
		},
		{
			name:      "無効なURL形式の場合はエラー",
			value:     func() *string { s := "://invalid"; return &s }(),
			want:      nil,
			assertion: assert.Error,
		},
		{
			name:      "空文字列の場合はエラー",
			value:     func() *string { s := ""; return &s }(),
			want:      nil,
			assertion: assert.Error,
		},
		{
			name:      "nilの場合は作成成功 (nilが返る)",
			value:     nil,
			want:      nil,
			assertion: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := vo.NewLink(tt.value)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestLink_String(t *testing.T) {
	t.Parallel()

	linkValue := "https://example.com/my-link"
	link, err := vo.NewLink(&linkValue)
	require.NoError(t, err)

	assert.Equal(t, linkValue, link.String())
}
