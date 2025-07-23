package entity

import (
	"fmt"
	"time"

	"github.com/umekikazuya/momenture-article-hub/internal/domain/vo"
)

// Article は記事のドメインエンティティ
type Article struct {
	ID           uint64
	Title        vo.ArticleTitle
	Body         *vo.ArticleBody
	Status       vo.ArticleStatus
	ProviderType *vo.ProviderType
	Link         *vo.Link
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

// ArticleOption は記事作成時のオプション設定用
type ArticleOption func(*Article) error

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

// NewArticle は新しい記事を作成する
func NewArticle(
	title string,
	status string,
	opts ...ArticleOption,
) (*Article, error) {
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

	for _, opt := range opts {
		if err := opt(article); err != nil {
			return nil, fmt.Errorf("failed to apply article option: %w", err)
		}
	}

	return article, nil
}

// ReconstituteArticle は永続化層から読み込んだデータから記事を再構築する
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

	// 全てのフィールドが引数で提供される前提
	return article, nil
}

// Publish は記事を公開状態に変更する
func (a *Article) Publish() error {
	if a.Status.IsPublished() {
		return fmt.Errorf("article is already published")
	}
	a.Status = vo.ArticleStatusPublished
	a.UpdatedAt = time.Now()
	return nil
}

// Draft は記事を下書き状態に変更する
func (a *Article) Draft() error {
	if a.Status.IsDraft() {
		return fmt.Errorf("article is already in draft status")
	}
	a.Status = vo.ArticleStatusDraft
	a.UpdatedAt = time.Now()
	return nil
}

// SoftDelete は記事を論理削除する
func (a *Article) SoftDelete() error {
	if a.DeletedAt != nil {
		return fmt.Errorf("article is already soft deleted")
	}
	now := time.Now()
	a.DeletedAt = &now
	a.UpdatedAt = now
	return nil
}

// Restore は論理削除された記事を復元する
func (a *Article) Restore() error {
	if a.DeletedAt == nil {
		return fmt.Errorf("article is not soft deleted")
	}
	a.DeletedAt = nil
	a.UpdatedAt = time.Now()
	return nil
}

// ChangeProvider は記事のプロバイダを変更する
// 公開済みの記事は変更不可
func (a *Article) ChangeProvider(newProviderType *vo.ProviderType) error {
	if a.Status.IsPublished() {
		return fmt.Errorf("cannot change provider for a published article")
	}
	a.ProviderType = newProviderType
	a.UpdatedAt = time.Now()
	return nil
}

// Update は記事の属性を更新する
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
