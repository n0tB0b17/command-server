package command

import (
	"sync"

	"github.com/n0tB0b17/isner/internal/models"
)

type CommandHandler func(clientID string, args []string) (string, error)

type CommandRegistry struct {
	commands map[models.CommandType]CommandHandler
	mu       sync.Mutex
}

func NewCommandRegistry() *CommandRegistry {
	return &CommandRegistry{
		commands: make(map[models.CommandType]CommandHandler),
	}
}

func (r *CommandRegistry) Register(t models.CommandType, h CommandHandler) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.commands[t] = h
}

func (r *CommandRegistry) Get(t models.CommandType) (CommandHandler, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	handler, exist := r.commands[t]
	return handler, exist
}

func (r *CommandRegistry) InitializeDefaultCmd() {
	r.Register(models.CmdEcho, func(clientID string, args []string) (string, error) {
		if len(args) == 0 {
			return "", ErrInvalidCommand
		}

		return args[0], nil
	})

	r.Register(models.CmdSystem, func(clientID string, args []string) (string, error) {
		return "system_info", nil
	})
}
