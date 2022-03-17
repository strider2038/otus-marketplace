package mock

import (
	"context"
	"testing"

	"trading-service/internal/messaging"
)

type MessageDispatcher struct {
	messages []messaging.Message
}

func NewMessageDispatcher() *MessageDispatcher {
	return &MessageDispatcher{}
}

func (mock *MessageDispatcher) Dispatch(ctx context.Context, message messaging.Message) error {
	mock.messages = append(mock.messages, message)

	return nil
}

func (mock *MessageDispatcher) AssertMessage(t *testing.T, index int, assert func(t *testing.T, message messaging.Message)) {
	if len(mock.messages) <= index {
		t.Errorf("message at index %d does not exist", index)
		return
	}

	assert(t, mock.messages[index])
}
