package vo

// ArticleBody は記事本文を表すValue Object
type ArticleBody string

// NewArticleBody は記事本文を作成する
// 空文字列は値なしとして扱う(必須ではない)
func NewArticleBody(value *string) (*ArticleBody, error) {
	if value == nil {
		return nil, nil
	}
	if *value == "" {
		return nil, nil
	}
	body := ArticleBody(*value)
	return &body, nil
}

func (b *ArticleBody) String() string {
	if b == nil {
		return ""
	}
	return string(*b)
}
