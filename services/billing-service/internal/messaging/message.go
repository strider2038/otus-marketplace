package messaging

import "context"

type Message interface {
	Name() string
}

type Dispatcher interface {
	Dispatch(ctx context.Context, message Message) error
}
