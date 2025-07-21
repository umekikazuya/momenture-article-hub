package vo

import "fmt"

// ProviderType は記事投稿先プロバイダの種類を表すValue Object
type ProviderType string

const (
	ProviderTypeQiita ProviderType = "qiita"
	ProviderTypeZenn  ProviderType = "zenn"
)

var AllProviderTypes = []ProviderType{
	ProviderTypeQiita,
	ProviderTypeZenn,
}

func NewProviderType(value *string) (*ProviderType, error) {
	if value == nil {
		return nil, nil
	}
	if *value == "" {
		return nil, nil
	}
	pt := ProviderType(*value)
	if !pt.IsValid() {
		return nil, fmt.Errorf("invalid provider type: %s", *value)
	}
	return &pt, nil
}

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

func (pt *ProviderType) String() string {
	if pt == nil {
		return ""
	}
	return string(*pt)
}

// API連携の可否を判定
func (pt *ProviderType) IsManual() bool {
	if pt == nil {
		return false
	}
	switch pt {
	default:
		return false
	}
}

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
