package models

import "time"

type ConnectionStatus string

const (
	StatusConnected    ConnectionStatus = "connected"
	StatusDisconnected ConnectionStatus = "disconnected"
	StatusActive       ConnectionStatus = "active"
)

type ClientMetadata struct {
	OS        string
	Arch      string
	Hostname  string
	IP        string
	FirstSeen time.Time
	LastSeen  time.Time
}
