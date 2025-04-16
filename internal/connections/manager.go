package connections

import (
	"fmt"
	"sync"
	"time"
)

type ConnectionManager struct {
	connections map[string]*ClientConnection
	mu          sync.Mutex
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[string]*ClientConnection),
	}
}

func (cm *ConnectionManager) Register(c *ClientConnection) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	cm.connections[c.ID] = c
	fmt.Printf("[connections::ConnectionManager::Register]::> Client has been registered \n")
}

func (cm *ConnectionManager) Unregister(clientID string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if client, exist := cm.connections[clientID]; exist {
		client.Close()
		delete(cm.connections, clientID)
		fmt.Printf("[connections::ConnectionManager::Unregister]::> Client has been unregistered \n")
	}
}

func (cm *ConnectionManager) GetAClient(clientID string) (*ClientConnection, bool) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	client, exist := cm.connections[clientID]
	return client, exist
}

func (cm *ConnectionManager) GetAllClient() []*ClientConnection {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	clients := make([]*ClientConnection, 0, len(cm.connections))
	for _, client := range cm.connections {
		clients = append(clients, client)
	}

	return clients
}

func (cm *ConnectionManager) BroadCast(msgType int, docs []byte) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	for _, client := range cm.connections {
		if err := client.SendMessage(msgType, docs); err != nil {
			fmt.Printf("[connections::ConnectionManager::BroadCast]::> Broadcast failed: %v \n", err)
		}
	}
}

func (cm *ConnectionManager) ClearInActive(timeout time.Duration) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	now := time.Now()
	for id, client := range cm.connections {
		if now.Sub(client.LastActivity) > timeout {
			client.Close()
			delete(cm.connections, id)
			fmt.Printf("[connections::ConnectionManager::ClearInActive]::> Cleared InActive connection \n")
		}
	}
}
