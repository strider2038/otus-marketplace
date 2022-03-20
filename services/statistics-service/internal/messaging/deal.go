package messaging

import (
	"context"
	"encoding/json"
	"time"

	"statistics-service/internal/statistics"

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
	dailyDeals      statistics.DailyDealsRepository
	totalDailyDeals statistics.TotalDailyDealsRepository
	top10Deals      statistics.Top10DealsRepository
}

func NewDealSucceededProcessor(
	dailyDeals statistics.DailyDealsRepository,
	totalDailyDeals statistics.TotalDailyDealsRepository,
	top10Deals statistics.Top10DealsRepository,
) *DealSucceededProcessor {
	return &DealSucceededProcessor{
		dailyDeals:      dailyDeals,
		totalDailyDeals: totalDailyDeals,
		top10Deals:      top10Deals,
	}
}

func (p *DealSucceededProcessor) Process(ctx context.Context, message []byte) error {
	var deal DealSucceeded
	err := json.Unmarshal(message, &deal)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal DealSucceeded event")
	}

	err = p.processDeal(ctx, deal)
	if err != nil {
		return errors.WithMessage(err, "failed to process DealSucceeded event")
	}

	return nil
}

func (p *DealSucceededProcessor) processDeal(ctx context.Context, deal DealSucceeded) error {
	err := p.dailyDeals.Add(ctx, &statistics.DailyDeals{
		Date:     deal.CompletedAt.Format(statistics.DateLayout),
		ItemID:   deal.ItemID,
		ItemName: deal.ItemName,
		Count:    1,
		Amount:   deal.Amount,
	})
	if err != nil {
		return err
	}

	err = p.totalDailyDeals.Add(ctx, &statistics.TotalDailyDeals{
		Date:   deal.CompletedAt.Format(statistics.DateLayout),
		Count:  1,
		Amount: deal.Amount,
	})
	if err != nil {
		return err
	}

	err = p.top10Deals.Add(ctx, &statistics.Top10Deals{
		ItemID:   deal.ItemID,
		ItemName: deal.ItemName,
		Count:    1,
		Amount:   deal.Amount,
	})
	if err != nil {
		return err
	}

	return nil
}
