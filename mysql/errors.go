package mysql

import "errors"

var (
	ErrNilPointer     = errors.New("invalid data or nil pointer")
	ErrMustPtrData    = errors.New("data must be a pointer")
	ErrMissPrimaryKey = errors.New("primary key not found")
)
