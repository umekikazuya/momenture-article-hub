package entity

import (
	"fmt"
	"time"

	"github.com/umekikazuya/momenture-article-hub/internal/domain/vo"
)

// Article is the main entity representing an article in the system.
type Article struct {
	ID           uint64           // 記事のユニークID
	Title        vo.ArticleTitle  // 記事タイトル (値オブジェクト)
	BodyMarkdown vo.ArticleBody   // 記事本文 (値オブジェクト)
	Status       vo.ArticleStatus // 記事ステータス (値オブジェクト)
	ProviderType vo.ProviderType  // 投稿先プロバイダ (値オブジェクト)

	// オプショナルな属性はポインタ型で表現し、値が存在しない場合はnil
	Link *vo.Link // 外部URL (値オブジェクト, ポインタでNULL許容)

	CreatedAt time.Time  // 作成日時
	UpdatedAt time.Time  // 更新日時
	DeletedAt *time.Time // 論理削除日時 (ポインタでNULL許容。削除されていない場合はnil)
}

// ArticleOption はArticleを生成・設定するための関数型オプションです。
type ArticleOption func(*Article) error

// WithLink は記事のリンクを設定するArticleOptionです。
func WithLink(link string) ArticleOption {
	return func(a *Article) error {
		l, err := vo.NewLink(link)
		if err != nil {
			return fmt.Errorf("invalid link for article option: %w", err)
		}
		a.Link = l
		return nil
	}
}

// NewArticle は新しいArticleエンティティを生成するファクトリメソッドです。
// IDとタイムスタンプはここで初期化されます。
// オプショナルな設定はFunctional Optionで渡されます。
func NewArticle(
	title string,
	bodyMarkdown string,
	status string, // vo.ArticleStatusのstring表現
	providerType string, // vo.ProviderTypeのstring表現
	opts ...ArticleOption,
) (*Article, error) {
	// 値オブジェクトの生成とバリデーション (必須フィールド)
	artTitle, err := vo.NewArticleTitle(title)
	if err != nil {
		return nil, fmt.Errorf("failed to create article title: %w", err)
	}
	artBody, err := vo.NewArticleBody(bodyMarkdown)
	if err != nil {
		return nil, fmt.Errorf("failed to create article body: %w", err)
	}
	artStatus := vo.ArticleStatus(status)
	if !artStatus.IsValid() {
		return nil, fmt.Errorf("invalid article status: %s", status)
	}
	provType := vo.ProviderType(providerType)
	if !provType.IsValid() {
		return nil, fmt.Errorf("invalid provider type: %s", providerType)
	}

	now := time.Now()
	article := &Article{
		// IDは永続化時にDBから払い出されるため、ここでは0または初期値
		Title:        artTitle,
		BodyMarkdown: artBody,
		Status:       artStatus,
		ProviderType: provType,
		CreatedAt:    now,
		UpdatedAt:    now,
		DeletedAt:    nil, // 初期状態では削除されていない
	}

	// Functional Optionの適用
	for _, opt := range opts {
		if err := opt(article); err != nil {
			return nil, fmt.Errorf("failed to apply article option: %w", err)
		}
	}

	return article, nil
}

// Publish は記事のステータスを公開済みに変更します。
func (a *Article) Publish() error {
	if a.Status == vo.ArticleStatusPublished {
		return fmt.Errorf("article is already published")
	}
	a.Status = vo.ArticleStatusPublished
	a.UpdatedAt = time.Now() // 更新日時を更新
	return nil
}

// Draft は記事のステータスを下書きに変更します。
func (a *Article) Draft() error {
	if a.Status == vo.ArticleStatusDraft {
		return fmt.Errorf("article is already in draft status")
	}
	a.Status = vo.ArticleStatusDraft
	a.UpdatedAt = time.Now() // 更新日時を更新
	return nil
}

// SoftDelete は記事を論理削除します。
func (a *Article) SoftDelete() error {
	if a.DeletedAt != nil { // DeletedAtがnilでない場合は既に削除済み
		return fmt.Errorf("article is already soft deleted")
	}
	now := time.Now()
	a.DeletedAt = &now // ポインタに現在日時を設定
	a.UpdatedAt = now  // 更新日時も更新
	return nil
}

// Restore は論理削除された記事を復元します。
func (a *Article) Restore() error {
	if a.DeletedAt == nil { // DeletedAtがnilの場合は削除されていない
		return fmt.Errorf("article is not soft deleted")
	}
	a.DeletedAt = nil        // nilに戻す
	a.UpdatedAt = time.Now() // 更新日時も更新
	return nil
}

// ChangeProvider は記事のプロバイダタイプを変更します。
// ビジネスルール: 例として、既に公開済みの記事のプロバイダ変更は許可しない
func (a *Article) ChangeProvider(newProviderType vo.ProviderType) error {
	if !newProviderType.IsValid() {
		return fmt.Errorf("invalid new provider type: %s", newProviderType)
	}
	if a.Status == vo.ArticleStatusPublished {
		return fmt.Errorf("cannot change provider for a published article")
	}
	a.ProviderType = newProviderType
	a.UpdatedAt = time.Now() // 更新日時を更新
	return nil
}

// UpdateFromInput はユースケース層からの入力に基づいてArticleの属性を更新します。
// 各属性の更新ロジックやバリデーションをカプセル化します。
// このメソッドは、更新可能な属性のみを受け取るように設計します。
// 引数はポインタ型でnilの場合、そのフィールドは更新されないことを意味します。
func (a *Article) UpdateFromInput(
	title *string,
	bodyMarkdown *string,
	status *string,
	providerType *string,
	link *string,
	externalPlatformID *string,
	externalMetadata map[string]interface{},
) error {
	if title != nil {
		newTitle, err := vo.NewArticleTitle(*title)
		if err != nil {
			return fmt.Errorf("failed to update title: %w", err)
		}
		a.Title = newTitle
	}
	if bodyMarkdown != nil {
		newBody, err := vo.NewArticleBody(*bodyMarkdown)
		if err != nil {
			return fmt.Errorf("failed to update body: %w", err)
		}
		a.BodyMarkdown = newBody
	}
	if status != nil {
		newStatus := vo.ArticleStatus(*status)
		if !newStatus.IsValid() {
			return fmt.Errorf("invalid status provided for update: %s", *status)
		}
		// @todo ステータス変更のビジネスルールはPublish/Draftメソッド経由で適用するのが理想的だが、
		// ここでは簡略化のため直接代入。厳密にはa.Publish()/a.Draft()を呼ぶべきか検討
		a.Status = newStatus
	}
	if providerType != nil {
		newProvider := vo.ProviderType(*providerType)
		// ChangeProviderメソッドを使ってビジネスルールを適用
		if err := a.ChangeProvider(newProvider); err != nil {
			return fmt.Errorf("failed to change provider: %w", err)
		}
	}
	if link != nil {
		// linkが空文字列の場合にnilを設定するビジネスルールはここで考慮
		if *link == "" {
			a.Link = nil
		} else {
			newLink, err := vo.NewLink(*link)
			if err != nil {
				return fmt.Errorf("failed to update link: %w", err)
			}
			a.Link = newLink
		}
	}
	a.UpdatedAt = time.Now() // 最終更新日時を更新
	return nil
}
