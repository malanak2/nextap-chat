package handlers

import "errors"

var (
	ErrorTooLong           = errors.New("Parameter is too long")
	ErrorTooShort          = errors.New("Parameter is too short")
	ErrorInvalidId         = errors.New("Invalid id format")
	ErrorInvalidLimit      = errors.New("Invalid limit format")
	ErrorInvalidPage       = errors.New("Invalid page format")
	ErrorDatabase          = errors.New("Internal database error")
	ErrorMarshal           = errors.New("Failed to marshal value")
	ErrorInvalidParameters = errors.New("Invalid parameters")
	ErrorInvalidBody       = errors.New("Invalid body content")
	ErrorBadJwt            = errors.New("Invalid jwt")
	ErrorNoExist           = errors.New("This id does not exist")
)
