package connections

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/n0tB0b17/isner/internal/models"
	"github.com/n0tB0b17/isner/internal/utils"
)

type ClientConnection struct {
	ID           string `json:"id"`
	Socket       *websocket.Conn
	Metadata     models.ClientMetadata `json:"meta_data"`
	LastActivity time.Time             `json:"last_activity"`
	mu           sync.Mutex
}

func NewClientConnection(conn *websocket.Conn, meta models.ClientMetadata) *ClientConnection {
	return &ClientConnection{
		ID:           utils.GenerateUUID("client"),
		Socket:       conn,
		Metadata:     meta,
		LastActivity: time.Now(),
	}
}

func (c *ClientConnection) SendMessage(msgType int, docs []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.Socket == nil {
		return ErrConnectionClosed
	}

	c.LastActivity = time.Now()
	return c.Socket.WriteMessage(msgType, docs)
}

func (c *ClientConnection) IsActive() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.Socket != nil
}

func (c *ClientConnection) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.Socket != nil {
		err := c.Socket.Close()
		c.Socket = nil
		return err
	}

	return nil
}
