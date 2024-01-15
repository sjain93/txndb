package common

import "errors"

var (
	ErrNoDatastore       = errors.New("no datastore provided")
	ErrRecordNotFound    = errors.New("record not found")
	ErrUniqueKeyViolated = errors.New("duplicated key not allowed")
)

const (
	UniqueViolationErr = "23505"
)
