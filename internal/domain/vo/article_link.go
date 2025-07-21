package vo

import (
	"fmt"
	"net/url"
)

// Link は記事の外部リンクを表すValue Object
type Link string

func NewLink(value *string) (*Link, error) {
	if value == nil {
		return nil, nil
	}
	err := isValid(value)
	if err != nil {
		return nil, fmt.Errorf("invalid link: %w", err)
	}
	link := Link(*value)
	return &link, nil
}

func isValid(value *string) error {
	if len(*value) == 0 {
		return fmt.Errorf("link cannot be empty")
	}

	_, err := url.ParseRequestURI(*value)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	return nil
}

func (l *Link) String() string {
	return string(*l)
}
