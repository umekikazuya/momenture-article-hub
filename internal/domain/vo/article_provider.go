package vo

import "fmt"

// ProviderType は記事が紐づく外部プロバイダの種類を表す値オブジェクトです。
// その値によって識別され、不変です。
type ProviderType string

const (
	ProviderTypeQiita               ProviderType = "qiita"
	ProviderTypeZenn                ProviderType = "zenn"
	ProviderTypeNote                ProviderType = "note"
	ProviderTypePersonalBlog        ProviderType = "personal_blog"
	ProviderTypeOtherManualPlatform ProviderType = "other_manual_platform"
)

// AllProviderTypes は定義されている全てのProviderTypeのリストを返します。
var AllProviderTypes = []ProviderType{
	ProviderTypeQiita,
	ProviderTypeZenn,
	ProviderTypeNote,
	ProviderTypePersonalBlog,
	ProviderTypeOtherManualPlatform,
}

// IsValid はプロバイダタイプが有効な値であるかを検証します。
func (pt ProviderType) IsValid() bool {
	for _, p := range AllProviderTypes {
		if pt == p {
			return true
		}
	}
	return false
}

// String はProviderTypeの文字列表現を返します。
func (pt ProviderType) String() string {
	return string(pt)
}

// HasAPI このプロバイダタイプがAPI連携可能かどうかを判定します。
func (pt ProviderType) HasAPI() bool {
	switch pt {
	case ProviderTypeQiita, ProviderTypeZenn:
		return true
	default:
		return false
	}
}

// IsManual このプロバイダタイプが手動投稿を前提とするかどうかを判定します。
func (pt ProviderType) IsManual() bool {
	switch pt {
	case ProviderTypePersonalBlog, ProviderTypeOtherManualPlatform:
		return true
	default:
		return false
	}
}

// DisplayName は人間が読める表示名を返します。
func (pt ProviderType) DisplayName() string {
	switch pt {
	case ProviderTypeQiita:
		return "Qiita"
	case ProviderTypeZenn:
		return "Zenn"
	case ProviderTypeNote:
		return "Note"
	case ProviderTypePersonalBlog:
		return "個人ブログ"
	case ProviderTypeOtherManualPlatform:
		return "その他手動投稿"
	default:
		return fmt.Sprintf("Unknown Provider (%s)", pt) // fmtパッケージをインポート
	}
}
