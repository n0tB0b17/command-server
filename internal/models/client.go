package models

import (
	"time"

	"github.com/n0tB0b17/isner/internal/utils"
)

type Client struct {
	ID            string            `json:"client_id"`
	ConnectionID  string            `json:"connection_id"`
	Status        ConnectionStatus  `json:"status"`
	MetaData      ClientMetadata    `json:"meta_data"`
	CommandOutput map[string]string `json:"command_output"` // key: commandID, value: output
}

func NewClient(connID string, meta ClientMetadata) *Client {
	return &Client{
		ID:            utils.GenerateUUID("client"),
		ConnectionID:  connID,
		Status:        StatusConnected,
		MetaData:      meta,
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
