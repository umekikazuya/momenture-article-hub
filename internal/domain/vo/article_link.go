package vo

import (
	"fmt"
	"net/url"
)

// ValueObject - Link.
// 外部リンクを表す。
// 有効なURL形式(RFC準拠)
type Link string

// NewLink は新しいLink値オブジェクトを生成する。
func NewLink(value *string) (*Link, error) {
	if value == nil {
		return nil, nil
	}
	err := isValid(value)
	if err != nil {
		return nil, fmt.Errorf("invalid link: %w", err)
	}
	link := Link(*value)
	return &link, nil
}

// isValid はリンクの値が有効かどうかを検証する。
func isValid(value *string) error {
	// 空文字列は許容しない。
	if len(*value) == 0 {
		return fmt.Errorf("link cannot be empty")
	}

	// URLのパースと検証。
	_, err := url.ParseRequestURI(*value)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	return nil
}

// String はLinkの文字列表現（URL）を返す。
func (l *Link) String() string {
	return string(*l)
}
