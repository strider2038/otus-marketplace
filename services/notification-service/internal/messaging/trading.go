package messaging

import (
	"context"
	"encoding/json"
	"fmt"

	"notification-service/internal/notifications"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

type DealSucceeded struct {
	SellerID            uuid.UUID `json:"sellerId"`
	PurchaserID         uuid.UUID `json:"purchaserId"`
	ItemID              uuid.UUID `json:"itemId"`
	ItemName            string    `json:"itemName"`
	Amount              float64   `json:"amount"`
	SellerCommission    float64   `json:"sellerCommission"`
	PurchaserCommission float64   `json:"purchaserCommission"`
}

func (m DealSucceeded) Name() string {
	return "Trading/DealSucceeded"
}

type DealSucceededProcessor struct {
	users         notifications.UserRepository
	notifications notifications.Repository
}

func NewDealSucceededProcessor(
	users notifications.UserRepository,
	notifications notifications.Repository,
) *DealSucceededProcessor {
	return &DealSucceededProcessor{users: users, notifications: notifications}
}

func (p *DealSucceededProcessor) Process(ctx context.Context, message []byte) error {
	var order DealSucceeded
	err := json.Unmarshal(message, &order)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal DealSucceeded event")
	}

	seller, err := p.users.FindByID(ctx, order.SellerID)
	if err != nil {
		return errors.WithMessagef(err, "failed to find seller %s", order.SellerID)
	}

	err = p.notifications.Add(ctx, notifications.NewNotification(
		order.SellerID,
		fmt.Sprintf(
			`Dear, %s %s! You have successfully sold item "%s" for %.2f$ (commission %.2f$) on the marketplace.`,
			seller.FirstName,
			seller.LastName,
			order.ItemName,
			order.Amount,
			order.SellerCommission,
		),
	))
	if err != nil {
		return errors.WithMessagef(err, "failed to add DealSucceeded notification for seller %s", order.SellerID)
	}

	purchaser, err := p.users.FindByID(ctx, order.PurchaserID)
	if err != nil {
		return errors.WithMessagef(err, "failed to find purchaser %s", order.PurchaserID)
	}

	err = p.notifications.Add(ctx, notifications.NewNotification(
		order.PurchaserID,
		fmt.Sprintf(
			`Dear, %s %s! You have successfully bought item "%s" for %.2f$ (commission %.2f$) on the marketplace.`,
			purchaser.FirstName,
			purchaser.LastName,
			order.ItemName,
			order.Amount,
			order.PurchaserCommission,
		),
	))
	if err != nil {
		return errors.WithMessagef(err, "failed to add DealSucceeded notification for seller %s", order.PurchaserID)
	}

	return nil
}

type PurchaseFailed struct {
	PurchaserID uuid.UUID `json:"purchaserId"`
	ItemID      uuid.UUID `json:"itemId"`
	ItemName    string    `json:"itemName"`
	Amount      float64   `json:"amount"`
	Reason      string    `json:"reason"`
}

func (m PurchaseFailed) Name() string {
	return "Trading/PurchaseFailed"
}

type PurchaseFailedProcessor struct {
	users         notifications.UserRepository
	notifications notifications.Repository
}

func NewPurchaseFailedProcessor(
	users notifications.UserRepository,
	notifications notifications.Repository,
) *PurchaseFailedProcessor {
	return &PurchaseFailedProcessor{users: users, notifications: notifications}
}

func (p *PurchaseFailedProcessor) Process(ctx context.Context, message []byte) error {
	var purchase PurchaseFailed
	err := json.Unmarshal(message, &purchase)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal PurchaseFailed event")
	}

	user, err := p.users.FindByID(ctx, purchase.PurchaserID)
	if err != nil {
		return errors.WithMessagef(err, "failed to find user %s", purchase.PurchaserID)
	}

	err = p.notifications.Add(ctx, notifications.NewNotification(
		purchase.PurchaserID,
		fmt.Sprintf(
			`Dear, %s %s! Your purchase of item "%s" for %.2f$ was not completed due to: %s.`,
			user.FirstName,
			user.LastName,
			purchase.ItemName,
			purchase.Amount,
			purchase.Reason,
		),
	))
	if err != nil {
		return errors.WithMessagef(err, "failed to add PurchaseFailed notification for user %s", purchase.PurchaserID)
	}

	return nil
}
