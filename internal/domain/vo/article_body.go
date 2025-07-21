package vo

// Value Object - ArticleBody.
// 記事の本文を表す。
type ArticleBody string

// NewArticleBody は新しいArticleBody値オブジェクトを生成する。
func NewArticleBody(value *string) (*ArticleBody, error) {
	if value == nil {
		return nil, nil
	}
	body := ArticleBody(*value)
	return &body, nil
}

// String returns the string representation of ArticleBody.
func (b *ArticleBody) String() string {
	if b == nil {
		return ""
	}
	return string(*b)
}
