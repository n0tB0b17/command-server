package models

import (
	"time"

	"github.com/n0tB0b17/isner/internal/utils"
)

type CommandType string

const (
	CmdEcho   CommandType = "echo"
	CmdSystem CommandType = "system_info"
)

type Command struct {
	ID        string
	ClientID  string
	Type      CommandType
	Content   string
	Timestamp time.Time
}

type CommandResponse struct {
	CommandID string
	ClientID  string
	Output    string
	Success   bool
	Timestamp time.Time
}

func NewCommand(cmdType CommandType, clientID, content string) *Command {
	return &Command{
		ID:        utils.GenerateUUID("cmd"),
		ClientID:  clientID,
		Type:      cmdType,
		Content:   content,
		Timestamp: time.Now(),
	}
}
