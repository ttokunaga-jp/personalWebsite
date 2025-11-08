package repository

import "errors"

var (
	// ErrNotFound indicates the requested entity does not exist.
	ErrNotFound = errors.New("repository: not found")
	// ErrInvalidInput indicates the supplied data violates repository constraints.
	ErrInvalidInput = errors.New("repository: invalid input")
	// ErrDuplicate indicates the entity violates unique constraints.
	ErrDuplicate = errors.New("repository: duplicate")
	// ErrConflict indicates a concurrent modification or stale data was supplied.
	ErrConflict = errors.New("repository: conflict")
	// ErrNotImplemented indicates the repository does not support the requested operation.
	ErrNotImplemented = errors.New("repository: not implemented")
)
