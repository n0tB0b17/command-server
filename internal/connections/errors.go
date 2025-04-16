package connections

import "errors"

var (
	ErrInvalidMessageFormat = errors.New("invalid message format")
	ErrClientNotRegistered  = errors.New("client not registered")
	ErrConnectionClosed     = errors.New("connection closed")
	ErrClientValidation     = errors.New("client validation failed")
)
