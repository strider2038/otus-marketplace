package postgres

import (
	"context"

	notificationspkg "notification-service/internal/notifications"
	"notification-service/internal/postgres/database"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

type NotificationRepository struct {
	db *database.Queries
}

func NewNotificationRepository(db *database.Queries) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (repository *NotificationRepository) FindByUser(ctx context.Context, userID uuid.UUID) ([]*notificationspkg.Notification, error) {
	dbNotifications, err := repository.db.FindNotificationsByUser(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to find notifications")
	}

	notifications := make([]*notificationspkg.Notification, len(dbNotifications))
	for i := range dbNotifications {
		notifications[i] = &notificationspkg.Notification{
			ID:        dbNotifications[i].ID,
			UserID:    dbNotifications[i].UserID,
			Message:   dbNotifications[i].Message,
			CreatedAt: dbNotifications[i].CreatedAt,
		}
	}

	return notifications, nil
}

func (repository *NotificationRepository) Add(ctx context.Context, notification *notificationspkg.Notification) error {
	_, err := repository.db.CreateNotification(ctx, database.CreateNotificationParams{
		ID:      notification.ID,
		UserID:  notification.UserID,
		Message: notification.Message,
	})
	if err != nil {
		return errors.Wrap(err, "failed to create notification")
	}

	return nil
}
