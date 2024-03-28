package mysql

import "errors"

var (
	ErrMustPtrData    = errors.New("data must be a pointer")
	ErrMissPrimaryKey = errors.New("primary key not found")
)
