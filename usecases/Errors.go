package usecases

import "errors"

var (
	ErrorBadParameters = errors.New("Bad parameters")
	ErrorDatabase      = errors.New("Database error")
)
