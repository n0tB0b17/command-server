package command

import (
	"encoding/json"

	"github.com/n0tB0b17/isner/internal/events"
	"github.com/n0tB0b17/isner/internal/models"
	"github.com/n0tB0b17/isner/internal/storage"
)

type CommandProcessor struct {
	event *events.EventPublisher
	store storage.Repository
}

func NewCommandProcessor(e *events.EventPublisher, repo storage.Repository) *CommandProcessor {
	return &CommandProcessor{
		event: e,
		store: repo,
	}
}

func (cp *CommandProcessor) ProcessCommandResponse(resp models.CommandResponse) error {
	client, err := cp.store.GetClient(resp.ClientID)
	if err != nil {
		return err
	}

	client.AddCommand(resp.CommandID, resp.Output)
	if err := cp.store.UpdateClient(client); err != nil {
		return err
	}

	return nil
}

// command related consumer
func (cp *CommandProcessor) RegisterEventHandler(consumer *events.EventConsumer) {
	consumer.RegisterEvent(events.EventCommandCompleted, func(event events.RabbitMQEvent) error {
		var resp models.CommandResponse
		if err := json.Unmarshal(event.Event.Payload.([]byte), &resp); err != nil {
			return err
		}

		return cp.ProcessCommandResponse(resp)
	})
}
