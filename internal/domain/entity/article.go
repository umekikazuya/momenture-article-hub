package entity

import (
	"fmt"
	"time"

	"github.com/umekikazuya/momenture-article-hub/internal/domain/vo"
)

// 記事コンテンツのドメインエンティティ。
type Article struct {
	ID           uint64           // 記事のユニークID (永続化時に払い出される)
	Title        vo.ArticleTitle  // 記事タイトル
	Body         *vo.ArticleBody  // 記事本文
	Status       vo.ArticleStatus // 記事ステータス
	ProviderType *vo.ProviderType // 投稿先プロバイダ
	Link         *vo.Link         // 外部URL
	CreatedAt    time.Time        // 作成日時
	UpdatedAt    time.Time        // 更新日時
	DeletedAt    *time.Time       // 論理削除日時
}

// 記事の属性を更新するためFunctional Optionパターンを使用。
type ArticleOption func(*Article) error

// オプション関数。
func WithBody(body *string) ArticleOption {
	return func(a *Article) error {
		b, err := vo.NewArticleBody(body)
		if err != nil {
			return fmt.Errorf("invalid body for article option: %w", err)
		}
		a.Body = b
		return nil
	}
}

// オプション関数。
func WithLink(link *string) ArticleOption {
	return func(a *Article) error {
		l, err := vo.NewLink(link)
		if err != nil {
			return fmt.Errorf("invalid link for article option: %w", err)
		}
		a.Link = l
		return nil
	}
}

// オプション関数。
func WithProviderType(providerType *string) ArticleOption {
	return func(a *Article) error {
		pt, err := vo.NewProviderType(providerType)
		if err != nil {
			return fmt.Errorf("invalid provider type for article option: %w", err)
		}
		a.ProviderType = pt
		return nil
	}
}

// Factory Method.
// 新しいArticleエンティティを生成する。
func NewArticle(
	title string,
	status string,
	opts ...ArticleOption,
) (*Article, error) {
	// 必須フィールドの検証と組み立て。
	artTitle, err := vo.NewArticleTitle(title)
	if err != nil {
		return nil, fmt.Errorf("failed to create article title: %w", err)
	}
	artStatus := vo.ArticleStatus(status)
	if !artStatus.IsValid() {
		return nil, fmt.Errorf("invalid article status: %s", status)
	}
	now := time.Now()

	article := &Article{
		Title:     artTitle,
		Status:    artStatus,
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil,
	}

	// Functional Optionの適用
	for _, opt := range opts {
		if err := opt(article); err != nil {
			return nil, fmt.Errorf("failed to apply article option: %w", err)
		}
	}

	return article, nil
}

// Factory Method.
// データベースなどから読み込んだ既存のArticleエンティティを再構築。
// - IDとタイムスタンプは既に存在するものとして受け取る。
// - 永続化層（リポジトリ実装）からのみ呼び出されることを想定。
func ReconstituteArticle(
	id uint64,
	title string,
	status string,
	body *string,
	providerType *string,
	link *string,
	createdAt time.Time,
	updatedAt time.Time,
	deletedAt *time.Time,
) (*Article, error) {
	artTitle, err := vo.NewArticleTitle(title)
	if err != nil {
		return nil, fmt.Errorf("failed to reconstitute article title: %w", err)
	}

	var artBody *vo.ArticleBody
	if body != nil {
		ab, err := vo.NewArticleBody(body)
		if err != nil {
			return nil, fmt.Errorf("failed to reconstitute article body: %w", err)
		}
		artBody = ab
	}

	artStatus := vo.ArticleStatus(status)
	if !artStatus.IsValid() {
		return nil, fmt.Errorf("invalid article status for reconstitution: %s", status)
	}

	var provType *vo.ProviderType
	if providerType != nil {
		pt, err := vo.NewProviderType(providerType)
		if err != nil {
			return nil, fmt.Errorf("failed to reconstitute article provider type: %w", err)
		}
		provType = pt
	}

	var artLink *vo.Link
	if link != nil {
		al, err := vo.NewLink(link)
		if err != nil {
			return nil, fmt.Errorf("failed to reconstitute article link: %w", err)
		}
		artLink = al
	}

	article := &Article{
		ID:           id,
		Title:        artTitle,
		Body:         artBody,
		Status:       artStatus,
		ProviderType: provType,
		Link:         artLink,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
		DeletedAt:    deletedAt,
	}

	// 全てのフィールドが引数で提供される前提のためFunctional Optionは使わない。
	return article, nil
}

// 記事のステータスを公開済みに変更。
func (a *Article) Publish() error {
	if a.Status.IsPublished() {
		return fmt.Errorf("article is already published")
	}
	a.Status = vo.ArticleStatusPublished
	a.UpdatedAt = time.Now()
	return nil
}

// 記事のステータスを下書きに変更。
func (a *Article) Draft() error {
	if a.Status.IsDraft() {
		return fmt.Errorf("article is already in draft status")
	}
	a.Status = vo.ArticleStatusDraft
	a.UpdatedAt = time.Now()
	return nil
}

// 記事を論理削除。
func (a *Article) SoftDelete() error {
	if a.DeletedAt != nil {
		return fmt.Errorf("article is already soft deleted")
	}
	now := time.Now()
	a.DeletedAt = &now
	a.UpdatedAt = now
	return nil
}

// 論理削除された記事を復元。
func (a *Article) Restore() error {
	if a.DeletedAt == nil {
		return fmt.Errorf("article is not soft deleted")
	}
	a.DeletedAt = nil
	a.UpdatedAt = time.Now()
	return nil
}

// ChangeProvider は記事のプロバイダタイプを変更。
// 既に公開済みの記事のプロバイダ変更は許可しない。
func (a *Article) ChangeProvider(newProviderType *vo.ProviderType) error {
	if a.Status.IsPublished() {
		return fmt.Errorf("cannot change provider for a published article")
	}
	a.ProviderType = newProviderType
	a.UpdatedAt = time.Now()
	return nil
}

// Articleの属性を更新する。
func (a *Article) Update(
	title *string,
	body *string,
	status *string,
	providerType *string,
	link *string,
) error {
	if title != nil {
		newTitle, err := vo.NewArticleTitle(*title)
		if err != nil {
			return fmt.Errorf("failed to update title: %w", err)
		}
		a.Title = newTitle
	}
	if body != nil {
		newBody, err := vo.NewArticleBody(body)
		if err != nil {
			return fmt.Errorf("failed to update body: %w", err)
		}
		a.Body = newBody
	} else {
		a.Body = nil
	}
	if status != nil {
		newStatus := vo.ArticleStatus(*status)
		if !newStatus.IsValid() {
			return fmt.Errorf("invalid status provided for update: %s", *status)
		}
		a.Status = newStatus
	}
	if providerType != nil {
		newProvider, err := vo.NewProviderType(providerType)
		if err != nil {
			return fmt.Errorf("failed to change provider: %w", err)
		}
		if err := a.ChangeProvider(newProvider); err != nil {
			return fmt.Errorf("failed to change provider: %w", err)
		}
		a.ProviderType = newProvider
	} else {
		a.ProviderType = nil
	}
	if link != nil {
		newLink, err := vo.NewLink(link)
		if err != nil {
			return fmt.Errorf("failed to update link: %w", err)
		}
		a.Link = newLink
	} else {
		a.Link = nil
	}

	a.UpdatedAt = time.Now()
	return nil
}
