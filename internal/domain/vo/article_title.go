package vo

import "fmt"

// ArticleTitle は記事のタイトルを表す値オブジェクトです。
// 不変であり、100文字以内というビジネスルールを持ちます。
type ArticleTitle string

const MaxArticleTitleLength = 100

// NewArticleTitle は新しいArticleTitle値オブジェクトを生成します。
// 生成時にバリデーションを行い、無効な場合はエラーを返します。
func NewArticleTitle(value string) (ArticleTitle, error) {
	if len(value) == 0 {
		return "", fmt.Errorf("article title cannot be empty")
	}
	if len(value) > MaxArticleTitleLength {
		return "", fmt.Errorf("article title exceeds maximum length of %d characters", MaxArticleTitleLength)
	}
	return ArticleTitle(value), nil
}

// String はArticleTitleの文字列表現を返します。
func (t ArticleTitle) String() string {
	return string(t)
}
