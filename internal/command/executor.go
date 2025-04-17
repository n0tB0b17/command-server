package command

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/n0tB0b17/isner/internal/connections"
	"github.com/n0tB0b17/isner/internal/events"
	"github.com/n0tB0b17/isner/internal/models"
)

type CmdExecutor struct {
	conn     *connections.ConnectionManager
	pub      *events.EventPublisher
	registry *CommandRegistry
}

func NewCmdExecutor(
	conn *connections.ConnectionManager,
	pub *events.EventPublisher,
	registry *CommandRegistry,
) *CmdExecutor {
	return &CmdExecutor{
		conn:     conn,
		pub:      pub,
		registry: registry,
	}
}

func (e *CmdExecutor) ExecuteCommand(clientID string, cmdType models.CommandType, args []string) error {
	_, exist := e.registry.Get(cmdType)
	if !exist {
		return ErrCommandNotFound
	}

	conn, exist := e.conn.GetAClient(clientID)
	if !exist {
		return ErrClientNotFound
	}

	cmd := models.NewCommand(cmdType, clientID, "")
	cmd.Content = buildArgs(args)

	x := map[string]interface{}{
		"type":       "command",
		"command_id": cmd.ID,
		"cmd":        string(cmdType),
		"args":       args,
	}

	bytes, err := json.Marshal(x)
	if err != nil {
		return err
	}

	if err := conn.SendMessage(websocket.TextMessage, bytes); err != nil {
		return err
	}

	e.pub.Publish(events.EventCommandIssued, cmd)
	fmt.Printf("[command::ExecuteCommand]::> command send to client: %s, type: %s \n", clientID, cmdType)
	return nil
}

func buildArgs(args []string) string {
	cmd := ""
	for _, arg := range args {
		cmd += arg + " "
	}

	return cmd
}
