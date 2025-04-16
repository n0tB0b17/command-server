package connections

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/n0tB0b17/isner/internal/events"
	"github.com/n0tB0b17/isner/internal/models"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocketHandler struct {
	connManager *ConnectionManager
	eventPub    *events.EventPublisher
}

func NewWebSocketHandler(conn *ConnectionManager, pub *events.EventPublisher) *WebSocketHandler {
	return &WebSocketHandler{
		connManager: conn,
		eventPub:    pub,
	}
}

func (wsh *WebSocketHandler) ServeHttp(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("[wsh::ServeHttp]::> error while upgrading websocket connection: %v \n", err)
		return
	}

	_, msg, err := ws.ReadMessage()
	if err != nil {
		fmt.Printf("[wsh::ServeHttp]::> error while reading message: %v \n", err)
		ws.Close()
		return
	}

	var meta models.ClientMetadata
	if err := json.Unmarshal(msg, &meta); err != nil {
		fmt.Printf("[wsh::ServeHttp]::> error while parsing message: %v \n", err)
		ws.Close()
		return
	}

	client := NewClientConnection(ws, meta)
	wsh.connManager.Register(client)
	// publish to rabbitMQ
	go wsh.handleClientMessage(client)
}

func (wsh *WebSocketHandler) handleClientMessage(client *ClientConnection)           {}
func (wsh *WebSocketHandler) handleTextMessage(client *ClientConnection, msg []byte) {}
func (wsh *WebSocketHandler) handleCommandResponse(clientID, commandID string, resp interface{})
