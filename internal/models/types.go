package models

import "time"

type ConnectionStatus string

const (
	StatusConnected    ConnectionStatus = "connected"
	StatusDisconnected ConnectionStatus = "disconnected"
	StatusActive       ConnectionStatus = "active"
)

type ClientMetadata struct {
	OS        string    `json:"os"`
	Arch      string    `json:"arch"`
	Hostname  string    `json:"host_name"`
	IP        string    `json:"ip"`
	FirstSeen time.Time `json:"first_seen"`
	LastSeen  time.Time `json:"last_seen"`
}
