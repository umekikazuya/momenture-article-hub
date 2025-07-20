package vo

import (
	"fmt"
	"net/url" // URLのパースとバリデーション用
)

// Link は記事の外部URLを表す値オブジェクトです。
// 有効なURL形式であるという不変条件を持ちます。
type Link string

// NewLink は新しいLink値オブジェクトを生成します。
// 生成時にURL形式のバリデーションを行い、無効な場合はエラーを返します。
// RFC準拠。クエリパラメータやハッシュは許容します。
func NewLink(value string) (*Link, error) {
	// 空文字列は許容しない。OptionalなリンクはNewArticle/UpdateFromInputでnilを渡す
	if len(value) == 0 {
		return nil, fmt.Errorf("link cannot be empty")
	}

	// URLのパースと検証
	parsedURL, err := url.ParseRequestURI(value) // RFC 3986 に厳密
	if err != nil {
		return nil, fmt.Errorf("invalid URL format: %w", err)
	}
	// スキーム（httpまたはhttps）が必須
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return nil, fmt.Errorf("URL scheme must be http or https")
	}
	// ホスト名が空でないこと（有効なURLには必要）
	if parsedURL.Host == "" {
		return nil, fmt.Errorf("URL must have a host")
	}

	link := Link(value)
	return &link, nil
}

// String はLinkの文字列表現（URL）を返します。
func (l Link) String() string {
	return string(l)
}
