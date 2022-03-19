package kafka

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/segmentio/kafka-go"
)

type Processor interface {
	Process(ctx context.Context, message []byte) error
}

type Consumer struct {
	reader     *kafka.Reader
	logger     zerolog.Logger
	processors map[string]Processor
}

func NewConsumer(reader *kafka.Reader, logger zerolog.Logger, processors map[string]Processor) *Consumer {
	return &Consumer{reader: reader, logger: logger, processors: processors}
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

		processor := c.processors[name]
		if processor == nil {
			c.logger.Error().Str("name", name).Msg("processor not found")
		} else {
			err = processor.Process(ctx, message.Value)
			c.logger.Info().Str("name", name).Dur("processingTime", time.Since(start)).Msg("message processed")
			if err != nil {
				c.logger.Info().Str("name", name).Err(err).Msg("processing failed")
			}
		}

		err = c.reader.CommitMessages(ctx, message)
		if err != nil {
			return errors.Wrap(err, "failed to commit message")
		}

		c.logger.Info().Str("name", name).Dur("processingTime", time.Since(start)).Msg("message committed")
	}
}

func getMessageName(message kafka.Message) string {
	for _, header := range message.Headers {
		if header.Key == "name" {
			return string(header.Value)
		}
	}

	return ""
}
