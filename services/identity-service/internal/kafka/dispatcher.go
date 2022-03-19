package kafka

import (
	"context"
	"encoding/json"

	"identity-service/internal/events"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/segmentio/kafka-go"
)

type Dispatcher struct {
	writer *kafka.Writer
	logger zerolog.Logger
}

func NewDispatcher(writer *kafka.Writer, logger zerolog.Logger) *Dispatcher {
	return &Dispatcher{writer: writer, logger: logger}
}

func (d *Dispatcher) Dispatch(ctx context.Context, event events.Event) error {
	value, err := json.Marshal(event)
	if err != nil {
		return errors.Wrapf(err, `failed to marshal event "%s"`, event.Name())
	}

	err = d.writer.WriteMessages(ctx, kafka.Message{
		Value: value,
		Headers: []kafka.Header{
			{Key: "name", Value: []byte(event.Name())},
		},
	})
	if err != nil {
		return errors.Wrapf(err, `failed to dispatch event "%s"`, event.Name())
	}

	d.logger.Info().
		Str("name", event.Name()).
		RawJSON("body", value).
		Msg("Message dispatched")

	return nil
}
