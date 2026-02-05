package domain

import "errors"

var (
	ErrUserNotFound = errors.Join(errors.New("product not found"))
)
