package vo

import "fmt"

// ArticleTitle は記事タイトルを表すValue Object
type ArticleTitle string

// タイトルの最大文字数制限
const MaxArticleTitleLength = 100

func NewArticleTitle(value string) (ArticleTitle, error) {
	if len(value) == 0 {
		return "", fmt.Errorf("article title cannot be empty")
	}
	if len(value) > MaxArticleTitleLength {
		return "", fmt.Errorf("article title exceeds maximum length of %d characters", MaxArticleTitleLength)
	}
	return ArticleTitle(value), nil
}

func (t ArticleTitle) String() string {
	return string(t)
}
