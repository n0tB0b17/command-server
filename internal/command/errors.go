package command

import "errors"

var (
	ErrCommandNotFound = errors.New("command not found in registry")
	ErrInvalidCommand  = errors.New("invalid command format")
	ErrClientNotFound  = errors.New("client not found")
	ErrCommandNotSent  = errors.New("failed to send command to client")
	ErrExecutionFailed = errors.New("command execution failed")
)
