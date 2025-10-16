package repository

import "errors"

var (
	// ErrNotFound indicates the requested entity does not exist.
	ErrNotFound = errors.New("repository: not found")
	// ErrInvalidInput indicates the supplied data violates repository constraints.
	ErrInvalidInput = errors.New("repository: invalid input")
	// ErrDuplicate indicates the entity violates unique constraints.
	ErrDuplicate = errors.New("repository: duplicate")
)
