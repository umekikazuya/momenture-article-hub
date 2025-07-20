package vo

// ArticleStatus は記事の公開ステータスを表す値オブジェクトです。
type ArticleStatus string

const (
	ArticleStatusDraft     ArticleStatus = "draft"
	ArticleStatusPublished ArticleStatus = "published"
)

// AllArticleStatuses は定義されている全てのArticleStatusのリストを返します。
var AllArticleStatuses = []ArticleStatus{
	ArticleStatusDraft,
	ArticleStatusPublished,
}

// IsValid はArticleStatusが有効な値であるかを検証します。
func (as ArticleStatus) IsValid() bool {
	for _, s := range AllArticleStatuses {
		if as == s {
			return true
		}
	}
	return false
}

// IsDraft はステータスが下書きであるかを判定します。
func (as ArticleStatus) IsDraft() bool {
	return as == ArticleStatusDraft
}

// IsPublished はステータスが公開済みであるかを判定します。
func (as ArticleStatus) IsPublished() bool {
	return as == ArticleStatusPublished
}

// String はArticleStatusの文字列表現を返します。
func (as ArticleStatus) String() string {
	return string(as)
}
