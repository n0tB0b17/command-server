package connections

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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

	wsh.eventPub.Publish(events.EventClientConnected, map[string]interface{}{
		"client_id": client.ID,
		"meta_data": meta,
	})

	go wsh.handleClientMessage(client)
}

func (wsh *WebSocketHandler) handleClientMessage(client *ClientConnection) {
	defer func() {
		wsh.connManager.Unregister(client.ID) // removes from manager's map
		client.Close()
	}()

	for {
		if !client.IsActive() {
			break
		}

		msgType, msg, err := client.Socket.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("[connections::handleClientMessage]::> Unexpected closing error: %v \n", err)
			}

			break
		}

		if msgType == websocket.TextMessage {
			wsh.handleTextMessage(client, msg)
		}
	}

	wsh.eventPub.Publish(events.EventClientDisconnected, map[string]interface{}{
		"client_id": client.ID,
	})
}

func (wsh *WebSocketHandler) handleTextMessage(client *ClientConnection, msg []byte) {
	var payload struct {
		Type      string      `json:"type"`
		Data      interface{} `json:"data"`
		CommandID string      `json:"command_id,omitempty"`
	}

	if err := json.Unmarshal(msg, &payload); err != nil {
		fmt.Printf("[connections::handleTextMessage]::>  Unable to decode message: %v \n", err)
		return
	}

	switch payload.Type {
	case "ready":
		fmt.Printf("[connections::handleTextMessage]::> Client: %s is ready for command \n", client.ID)
	case "response":
		wsh.handleCommandResponse(client.ID, payload.CommandID, payload.Data)
	default:
		fmt.Printf("[connections::handleTextMessage]::> Running default case for handleTextMessage \n")
	}
}

func (wsh *WebSocketHandler) handleCommandResponse(clientID, commandID string, resp interface{}) {
	respByte, err := json.Marshal(resp)
	if err != nil {
		fmt.Printf("[connection::handleCommandResponse]::> Unable to marshal interface to bytes: %v \n", err)
		return
	}

	wsh.eventPub.Publish(events.EventCommandCompleted, models.CommandResponse{
		CommandID: commandID,
		ClientID:  clientID,
		Output:    string(respByte),
		Success:   true,
		Timestamp: time.Now(),
	})
}
