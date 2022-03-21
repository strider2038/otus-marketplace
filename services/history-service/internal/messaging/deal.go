package messaging

import (
	"context"
	"encoding/json"
	"time"

	"history-service/internal/history"

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
	CompletedAt         time.Time `json:"completedAt"`
}

func (m DealSucceeded) Name() string {
	return "Trading/DealSucceeded"
}

type DealSucceededProcessor struct {
	deals history.DealRepository
}

func NewDealSucceededProcessor(deals history.DealRepository) *DealSucceededProcessor {
	return &DealSucceededProcessor{deals: deals}
}

func (p *DealSucceededProcessor) Process(ctx context.Context, message []byte) error {
	var deal DealSucceeded
	err := json.Unmarshal(message, &deal)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal DealSucceeded event")
	}

	err = p.deals.Add(ctx, history.NewPurchase(
		deal.PurchaserID,
		deal.ItemID,
		deal.ItemName,
		deal.Amount,
		deal.PurchaserCommission,
		deal.CompletedAt,
	))
	if err != nil {
		return errors.WithMessagef(err, "failed to add purchase for user %s", deal.PurchaserID)
	}

	err = p.deals.Add(ctx, history.NewSale(
		deal.SellerID,
		deal.ItemID,
		deal.ItemName,
		deal.Amount,
		deal.SellerCommission,
		deal.CompletedAt,
	))
	if err != nil {
		return errors.WithMessagef(err, "failed to add sale for user %s", deal.SellerID)
	}

	return nil
}
