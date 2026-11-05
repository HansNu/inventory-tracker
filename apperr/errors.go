// Package apperr holds the domain errors shared by every layer.
//
// It exists so that the repository can translate storage-specific failures
// (pgx sentinels, Postgres error codes) into something the service and handler
// layers can reason about without importing a database driver. Nothing in here
// knows about HTTP or SQL.
package apperr

import "errors"

var (
	ErrNotFound           = errors.New("resource not found")
	ErrDuplicate          = errors.New("resource already exists")
	ErrReferenced         = errors.New("resource is still referenced by other records")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidSortField   = errors.New("invalid sort field")
)
