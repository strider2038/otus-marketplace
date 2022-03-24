package kafka

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/segmentio/kafka-go"
)

var errProcessorNotFound = errors.New("processor not found")

type Processor interface {
	Process(ctx context.Context, name string, message []byte) error
}

type Consumer struct {
	reader    *kafka.Reader
	logger    zerolog.Logger
	processor Processor
}

func NewConsumer(reader *kafka.Reader, logger zerolog.Logger, processor Processor) *Consumer {
	return &Consumer{reader: reader, logger: logger, processor: processor}
}

func (c *Consumer) Run(ctx context.Context) error {
	for {
		message, err := c.reader.FetchMessage(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to fetch message")
		}

		start := time.Now()
		name := getMessageName(message)
		c.logger.Info().Str("name", name).Msg("message received")

		err = c.processor.Process(ctx, name, message.Value)
		if err == nil {
			c.logger.Info().Str("name", name).Dur("processingTime", time.Since(start)).Msg("message processed")
		} else {
			c.logger.Info().Str("name", name).Err(err).Msg("processing failed")
		}

		err = c.reader.CommitMessages(ctx, message)
		if err != nil {
			return errors.Wrap(err, "failed to commit message")
		}

		c.logger.Info().Str("name", name).Dur("processingTime", time.Since(start)).Msg("message committed")
	}
}

type Mux struct {
	processors map[string]Processor
	logger     zerolog.Logger
}

func NewMux(logger zerolog.Logger, processors map[string]Processor) *Mux {
	return &Mux{processors: processors, logger: logger}
}

func (m *Mux) Process(ctx context.Context, name string, message []byte) error {
	processor := m.processors[name]
	if processor == nil {
		return errors.WithStack(errProcessorNotFound)
	}

	return processor.Process(ctx, name, message)
}

func getMessageName(message kafka.Message) string {
	for _, header := range message.Headers {
		if header.Key == "name" {
			return string(header.Value)
		}
	}

	return ""
}
