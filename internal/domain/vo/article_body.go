package vo

import "fmt"

// ArticleBody は記事の本文を表す値オブジェクトです。
// Markdown形式の文字列を保持し、空ではないという不変条件を持ちます。
type ArticleBody string

// NewArticleBody は新しいArticleBody値オブジェクトを生成します。
// 生成時にバリデーションを行い、無効な場合はエラーを返します。
func NewArticleBody(value string) (*ArticleBody, error) {
	if len(value) == 0 { // 要件定義で「空は許容しない」と判断した場合
		return nil, fmt.Errorf("article body cannot be empty")
	}
	body := ArticleBody(value)
	return &body, nil
}

// String はArticleBodyの文字列表現を返します。
func (b ArticleBody) String() string {
	return string(b)
}
