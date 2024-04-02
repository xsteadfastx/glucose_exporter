package config

import "errors"

var (
	ErrTooManyPasswords = errors.New("too many passwords")
	ErrMissingPassword  = errors.New("missing password")
)
