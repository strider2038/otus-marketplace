package history

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
)

type DealType string

const (
	PurchaseDeal DealType = "purchase"
	SaleDeal     DealType = "sale"
)

type Deal struct {
	ID          uuid.UUID `json:"id,omitempty"`
	UserID      uuid.UUID `json:"-"`
	ItemID      uuid.UUID `json:"itemId,omitempty"`
	ItemName    string    `json:"itemName,omitempty"`
	Type        DealType  `json:"type,omitempty"`
	Amount      float64   `json:"amount,omitempty"`
	Commission  float64   `json:"commission,omitempty"`
	CompletedAt time.Time `json:"completedAt,omitempty"`
}

func NewPurchase(
	userID, itemID uuid.UUID,
	itemName string,
	amount, commission float64,
	completedAt time.Time,
) *Deal {
	return newDeal(userID, itemID, itemName, PurchaseDeal, amount, commission, completedAt)
}

func NewSale(
	userID, itemID uuid.UUID,
	itemName string,
	amount, commission float64,
	completedAt time.Time,
) *Deal {
	return newDeal(userID, itemID, itemName, SaleDeal, amount, commission, completedAt)
}

func newDeal(
	userID, itemID uuid.UUID,
	itemName string,
	dealType DealType,
	amount, commission float64,
	completedAt time.Time,
) *Deal {
	return &Deal{
		ID:          uuid.Must(uuid.NewV4()),
		UserID:      userID,
		ItemID:      itemID,
		ItemName:    itemName,
		Type:        dealType,
		Amount:      amount,
		Commission:  commission,
		CompletedAt: completedAt,
	}
}

type DealRepository interface {
	FindByUser(ctx context.Context, userID uuid.UUID) ([]*Deal, error)
	Add(ctx context.Context, deal *Deal) error
}
