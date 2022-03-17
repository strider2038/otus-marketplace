package trading

import (
	"context"

	"github.com/gofrs/uuid"
)

type PurchaseFailedEvent struct {
	PurchaserID uuid.UUID `json:"purchaserId"`
	ItemID      uuid.UUID `json:"itemId"`
	ItemName    string    `json:"itemName"`
	Amount      float64   `json:"amount"`
	Reason      string    `json:"reason"`
}

func (event PurchaseFailedEvent) Name() string {
	return "Trading/PurchaseFailed"
}

type DealSucceededEvent struct {
	SellerID            uuid.UUID `json:"sellerId"`
	PurchaserID         uuid.UUID `json:"purchaserId"`
	ItemID              uuid.UUID `json:"itemId"`
	ItemName            string    `json:"itemName"`
	Amount              float64   `json:"amount"`
	SellerCommission    float64   `json:"sellerCommission"`
	PurchaserCommission float64   `json:"purchaserCommission"`
}

func (event DealSucceededEvent) Name() string {
	return "Trading/DealSucceeded"
}

type Event interface {
	Name() string
}

type EventPublisher interface {
	Publish(ctx context.Context, event Event) error
}
