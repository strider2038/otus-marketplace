package kafka

import (
	"context"
	"log"
	"time"

	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
)

type Processor interface {
	Process(ctx context.Context, message []byte) error
}

type Consumer struct {
	reader     *kafka.Reader
	processors map[string]Processor
}

func NewConsumer(reader *kafka.Reader, processors map[string]Processor) *Consumer {
	return &Consumer{reader: reader, processors: processors}
}

func (c *Consumer) Run(ctx context.Context) error {
	for {
		message, err := c.reader.FetchMessage(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to fetch message")
		}

		start := time.Now()
		name := getMessageName(message)
		log.Printf(`message "%s" received at %s`, name, start.Format(time.RFC3339))

		processor := c.processors[name]
		if processor == nil {
			log.Printf(`processor not found for message "%s"`, name)
		} else {
			err = processor.Process(ctx, message.Value)
			log.Printf(`message "%s" processed at %s (processing time %s)`, name, time.Now().Format(time.RFC3339), time.Since(start))
			if err != nil {
				log.Printf(`error while processing message "%s": %v`, name, err)
			}
		}

		err = c.reader.CommitMessages(ctx, message)
		if err != nil {
			return errors.Wrap(err, "failed to commit message")
		}

		log.Printf(`message "%s" committed at %s`, name, time.Now().Format(time.RFC3339))
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
