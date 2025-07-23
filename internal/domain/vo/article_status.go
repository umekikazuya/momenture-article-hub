package vo

// ArticleStatus は記事の公開ステータスを表すValue Object
type ArticleStatus string

const (
	ArticleStatusDraft     ArticleStatus = "draft"
	ArticleStatusPublished ArticleStatus = "published"
)

var AllArticleStatuses = []ArticleStatus{
	ArticleStatusDraft,
	ArticleStatusPublished,
}

func (as ArticleStatus) IsValid() bool {
	for _, s := range AllArticleStatuses {
		if as == s {
			return true
		}
	}
	return false
}

func (as ArticleStatus) IsDraft() bool {
	return as == ArticleStatusDraft
}

func (as ArticleStatus) IsPublished() bool {
	return as == ArticleStatusPublished
}

func (as ArticleStatus) String() string {
	return string(as)
}
