package vo

// Value Object - ArticleStatus.
// 記事の公開ステータスを表す。
type ArticleStatus string

const (
	ArticleStatusDraft     ArticleStatus = "draft"
	ArticleStatusPublished ArticleStatus = "published"
)

// AllArticleStatuses は定義されている全てのArticleStatusのリストを返す。
var AllArticleStatuses = []ArticleStatus{
	ArticleStatusDraft,
	ArticleStatusPublished,
}

// IsValid はArticleStatusが有効な値であるかを検証。
func (as ArticleStatus) IsValid() bool {
	for _, s := range AllArticleStatuses {
		if as == s {
			return true
		}
	}
	return false
}

// IsDraft はステータスが下書きであるかを判定。
func (as ArticleStatus) IsDraft() bool {
	return as == ArticleStatusDraft
}

// IsPublished はステータスが公開済みであるかを判定。
func (as ArticleStatus) IsPublished() bool {
	return as == ArticleStatusPublished
}

// String はArticleStatusの文字列表現を返す。
func (as ArticleStatus) String() string {
	return string(as)
}
