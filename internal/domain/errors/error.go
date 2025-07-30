package errors

import (
	"errors"
	"fmt"
)

// ErrNotFound はリソースが見つからないことを示すエラー。
var ErrNotFound = errors.New("not found")

// ErrInvalidInput は入力データが無効であることを示すエラー。
var ErrInvalidInput = errors.New("invalid input")

// ErrConflict はリソースの競合やビジネスルール違反を示すエラー。
var ErrConflict = errors.New("conflict")

// ErrUnauthorized は認証されていないことを示すエラー。
var ErrUnauthorized = errors.New("unauthorized")

// ErrForbidden は権限がないことを示すエラー。
var ErrForbidden = errors.New("forbidden")

// ErrInternal は予期せぬ内部エラーを示すエラー。
var ErrInternal = errors.New("internal error")

// NewNotFound は指定されたメッセージでErrNotFoundをラップ。
func NewNotFound(message string) error {
	return fmt.Errorf("%s: %w", message, ErrNotFound)
}

// NewInvalidInput は指定されたメッセージでErrInvalidInputをラップ。
func NewInvalidInput(message string) error {
	return fmt.Errorf("%s: %w", message, ErrInvalidInput)
}

// NewConflict は指定されたメッセージでErrConflictをラップ。
func NewConflict(message string) error {
	return fmt.Errorf("%s: %w", message, ErrConflict)
}

// NewInternal は指定されたメッセージでErrInternalをラップ。
func NewInternal(message string) error {
	return fmt.Errorf("%s: %w", message, ErrInternal)
}
