package ports

import "errors"

var (
	ErrorDuplicateUsername = errors.New("User with this username already exists")
	ErrorDatabase          = errors.New("Database error")
	ErrorNoResult          = errors.New("No results")
	ErrorHashing           = errors.New("Failed to hash password")
	ErrorBadCredentials    = errors.New("Invalid username or password")
)
