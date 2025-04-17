package storage

import "errors"

var (
	ErrClientNotFound    = errors.New("client not found")
	ErrDuplicateClient   = errors.New("client already exists")
	ErrInvalidClientData = errors.New("invalid client data")
	ErrCommandNotFound   = errors.New("command not found")
)
