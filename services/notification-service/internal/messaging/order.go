package messaging

import (
	"context"
	"encoding/json"
	"fmt"

	"notification-service/internal/notifications"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

type OrderSucceeded struct {
	UserID uuid.UUID `json:"userId"`
	Price  float64   `json:"price"`
}

type OrderFailed struct {
	UserID uuid.UUID `json:"userId"`
	Price  float64   `json:"price"`
	Reason string    `json:"reason"`
}

type OrderSucceededProcessor struct {
	users         notifications.UserRepository
	notifications notifications.Repository
}

func NewOrderSucceededProcessor(
	users notifications.UserRepository,
	notifications notifications.Repository,
) *OrderSucceededProcessor {
	return &OrderSucceededProcessor{users: users, notifications: notifications}
}

func (processor *OrderSucceededProcessor) Process(ctx context.Context, message []byte) error {
	var order OrderSucceeded
	err := json.Unmarshal(message, &order)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal OrderSucceeded event")
	}

	user, err := processor.users.FindByID(ctx, order.UserID)
	if err != nil {
		return errors.WithMessagef(err, "failed to find user %s", order.UserID)
	}

	err = processor.notifications.Add(ctx, notifications.NewNotification(
		order.UserID,
		fmt.Sprintf(
			"Dear, %s %s! Your order for %.2f$ has been succesfully processed.",
			user.FirstName,
			user.LastName,
			order.Price,
		),
	))
	if err != nil {
		return errors.WithMessagef(err, "failed to add OrderSucceeded notification for user %s", order.UserID)
	}

	return nil
}

type OrderFailedProcessor struct {
	users         notifications.UserRepository
	notifications notifications.Repository
}

func NewOrderFailedProcessor(
	users notifications.UserRepository,
	notifications notifications.Repository,
) *OrderFailedProcessor {
	return &OrderFailedProcessor{users: users, notifications: notifications}
}

func (processor *OrderFailedProcessor) Process(ctx context.Context, message []byte) error {
	var order OrderFailed
	err := json.Unmarshal(message, &order)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal OrderFailed event")
	}

	user, err := processor.users.FindByID(ctx, order.UserID)
	if err != nil {
		return errors.WithMessagef(err, "failed to find user %s", order.UserID)
	}

	err = processor.notifications.Add(ctx, notifications.NewNotification(
		order.UserID,
		fmt.Sprintf(
			"Dear, %s %s! Your order for %.2f$ was not fulfilled due to: %s.",
			user.FirstName,
			user.LastName,
			order.Price,
			order.Reason,
		),
	))
	if err != nil {
		return errors.WithMessagef(err, "failed to add OrderFailed notification for user %s", order.UserID)
	}

	return nil
}
