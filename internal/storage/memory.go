package storage

import (
	"sync"

	"github.com/n0tB0b17/isner/internal/models"
)

type MemoryStorage struct {
	clients        map[string]*models.Client
	commandHistory map[string][]*models.CommandResponse
	mu             sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		clients:        make(map[string]*models.Client),
		commandHistory: make(map[string][]*models.CommandResponse),
	}
}

func (ms *MemoryStorage) AddClient(c *models.Client) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if _, exist := ms.clients[c.ID]; exist {
		return ErrDuplicateClient
	}

	ms.clients[c.ID] = c
	return nil
}

func (ms *MemoryStorage) GetClient(clientID string) (*models.Client, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	client, exist := ms.clients[clientID]
	if !exist {
		return nil, ErrClientNotFound
	}

	return client, nil
}

func (ms *MemoryStorage) DeleteClient(clientID string) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if _, exist := ms.clients[clientID]; !exist {
		return ErrClientNotFound
	}

	delete(ms.clients, clientID)
	delete(ms.commandHistory, clientID)
	return nil
}

func (ms *MemoryStorage) UpdateClient(c *models.Client) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if _, exist := ms.clients[c.ID]; !exist {
		return ErrClientNotFound
	}

	ms.clients[c.ID] = c
	return nil
}

func (ms *MemoryStorage) ListClients() ([]*models.Client, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	if len(ms.clients) == 0 {
		return nil, ErrClientNotFound
	}

	clients := make([]*models.Client, 0, len(ms.clients))
	for _, client := range ms.clients {
		clients = append(clients, client)
	}

	return clients, nil
}

func (ms *MemoryStorage) AddCommentResponse(clientID string, resp *models.CommandResponse) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if _, exist := ms.clients[clientID]; !exist {
		return ErrClientNotFound
	}

	ms.commandHistory[clientID] = append(ms.commandHistory[clientID], resp)
	return nil
}

func (ms *MemoryStorage) GetCommandResponse(clientID, commandID string) (*models.CommandResponse, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	history, exists := ms.commandHistory[clientID]
	if !exists {
		return nil, ErrClientNotFound
	}

	for _, hist := range history {
		if hist.CommandID == commandID {
			return hist, nil
		}
	}

	return nil, ErrCommandNotFound
}

func (ms *MemoryStorage) GetCommandHistory(clientID string) ([]*models.CommandResponse, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	if _, exist := ms.clients[clientID]; !exist {
		return nil, ErrClientNotFound
	}

	return ms.commandHistory[clientID], nil
}

func (ms *MemoryStorage) Close() error {
	return nil
}
