package redis

import "errors"

var (
	ErrOutOfStock   = errors.New("out of stock")
	ErrInvalidField = errors.New("invalid field")
)
