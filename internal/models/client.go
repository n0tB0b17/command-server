package models

import (
	"time"

	"github.com/n0tB0b17/isner/internal/utils"
)

type Client struct {
	ID            string
	ConnectionID  string
	Status        ConnectionStatus
	MetaData      ClientMetadata
	CommandOutput map[string]string // key: commandID, value: output
}

func NewClient(connID string) *Client {
	return &Client{
		ID:           utils.GenerateUUID("client"),
		ConnectionID: connID,
		Status:       StatusConnected,
		MetaData: ClientMetadata{
			FirstSeen: time.Now(),
			LastSeen:  time.Now(),
		},
		CommandOutput: make(map[string]string),
	}
}

func (c *Client) UpdateLastSeen() {
	c.MetaData.LastSeen = time.Now()
}

func (c *Client) AddCommand(id, output string) {
	c.CommandOutput[id] = output
	c.UpdateLastSeen()
}
