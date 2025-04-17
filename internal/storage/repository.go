package storage

import "github.com/n0tB0b17/isner/internal/models"

type Repository interface {
	AddClient(c *models.Client) error
	GetClient(clientID string) (*models.Client, error)
	DeleteClient(clientID string) error
	UpdateClient(c *models.Client) error
	ListClients() ([]*models.Client, error)

	AddCommentResponse(clientID string, resp *models.CommandResponse) error
	GetCommandHistory(clientID string) ([]*models.CommandResponse, error) // client specific history
	GetCommandResponse(clientID, commandID string) (*models.CommandResponse, error)

	Close() error
}
