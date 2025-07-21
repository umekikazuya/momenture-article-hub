package vo

import "fmt"

// Value Object - ProviderType.
// 記事が紐づく外部プロバイダの種類を表す。
// その値によって識別され不変。
type ProviderType string

const (
	ProviderTypeQiita ProviderType = "qiita"
	ProviderTypeZenn  ProviderType = "zenn"
)

// AllProviderTypes は定義されている全てのProviderTypeのリストを返す。
var AllProviderTypes = []ProviderType{
	ProviderTypeQiita,
	ProviderTypeZenn,
}

func NewProviderType(value *string) (*ProviderType, error) {
	if value == nil {
		return nil, nil
	}
	pt := ProviderType(*value)
	if !pt.IsValid() {
		return nil, fmt.Errorf("invalid provider type: %s", *value)
	}
	return &pt, nil
}

// IsValid はプロバイダタイプが有効な値であるかを検証。
func (pt *ProviderType) IsValid() bool {
	if pt == nil {
		return false
	}
	for _, p := range AllProviderTypes {
		if *pt == p {
			return true
		}
	}
	return false
}

// String はProviderTypeの文字列表現を返す。
func (pt *ProviderType) String() string {
	return string(*pt)
}

// このプロバイダタイプが手動投稿を前提とするかどうかを判定(API連携)。
func (pt *ProviderType) IsManual() bool {
	if pt == nil {
		return false
	}
	switch pt {
	default:
		return false
	}
}

// 表示名を返す。
func (pt *ProviderType) DisplayName() string {
	if pt == nil {
		return fmt.Sprintf("Unknown Provider (%s)", pt)
	}
	switch *pt {
	case ProviderTypeQiita:
		return "Qiita"
	case ProviderTypeZenn:
		return "Zenn"
	default:
		return fmt.Sprintf("Unknown Provider (%s)", pt)
	}
}
