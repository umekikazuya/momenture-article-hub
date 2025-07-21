package vo_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/umekikazuya/momenture-article-hub/internal/domain/vo"
)

func TestNewProviderType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		value     *string
		want      *vo.ProviderType
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:  "有効なプロバイダタイプ (qiita) で作成成功",
			value: func() *string { s := string(vo.ProviderTypeQiita); return &s }(),
			want: func() *vo.ProviderType {
				pt, _ := vo.NewProviderType(func() *string { s := string(vo.ProviderTypeQiita); return &s }())
				return pt
			}(),
			assertion: assert.NoError,
		},
		{
			name:  "有効なプロバイダタイプ (zenn) で作成成功",
			value: func() *string { s := string(vo.ProviderTypeZenn); return &s }(),
			want: func() *vo.ProviderType {
				pt, _ := vo.NewProviderType(func() *string { s := string(vo.ProviderTypeZenn); return &s }())
				return pt
			}(),
			assertion: assert.NoError,
		},
		{
			name:      "無効なプロバイダタイプの場合はエラー",
			value:     func() *string { s := "invalid_provider"; return &s }(),
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
			got, err := vo.NewProviderType(tt.value)
			tt.assertion(t, err)
			// ポインタの比較ではなく値の比較を行う
			if tt.want != nil && got != nil {
				assert.Equal(t, *tt.want, *got)
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestProviderType_IsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		pt   vo.ProviderType
		want bool
	}{
		{name: "Qiitaは有効", pt: vo.ProviderTypeQiita, want: true},
		{name: "Zennは有効", pt: vo.ProviderTypeZenn, want: true},
		{name: "無効な値はfalse", pt: vo.ProviderType("invalid"), want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, tt.pt.IsValid())
		})
	}
}

func TestProviderType_IsManual(t *testing.T) {
	t.Parallel()
	// 現在の実装では常にfalseが返ることをテスト
	qiita := vo.ProviderTypeQiita
	zenn := vo.ProviderTypeZenn
	assert.False(t, qiita.IsManual())
	assert.False(t, zenn.IsManual())
}

func TestProviderType_DisplayName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		pt   vo.ProviderType
		want string
	}{
		{name: "Qiitaの表示名", pt: vo.ProviderTypeQiita, want: "Qiita"},
		{name: "Zennの表示名", pt: vo.ProviderTypeZenn, want: "Zenn"},
		{name: "不明なプロバイダの表示名", pt: "unknown", want: fmt.Sprintf("Unknown Provider (%s)", "unknown")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, tt.pt.DisplayName())
		})
	}
}
