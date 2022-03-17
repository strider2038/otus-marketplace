package messaging

import (
	"context"
	"encoding/json"

	"trading-service/internal/trading"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

type CreateAccrual struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"userId"`
	Amount      float64   `json:"amount"`
	Commission  float64   `json:"commission"`
	Description string    `json:"description"`
}

func (c CreateAccrual) Name() string {
	return "Billing/CreateAccrual"
}

type AccrualApproved struct {
	ID         uuid.UUID `json:"id"`
	Amount     float64   `json:"amount"`
	Commission float64   `json:"commission"`
}

func (p AccrualApproved) Name() string {
	return "Billing/AccrualApproved"
}

type AccrualApprovedProcessor struct {
	dealer *trading.Dealer
}

func NewAccrualApprovedProcessor(dealer *trading.Dealer) *AccrualApprovedProcessor {
	return &AccrualApprovedProcessor{dealer: dealer}
}

func (p *AccrualApprovedProcessor) Process(ctx context.Context, message []byte) error {
	var accrual AccrualApproved
	err := json.Unmarshal(message, &accrual)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal AccrualApproved message")
	}

	err = p.dealer.ApproveAccrual(ctx, &trading.Accrual{
		ID:         accrual.ID,
		Amount:     accrual.Amount,
		Commission: accrual.Commission,
	})
	if err != nil {
		return errors.WithMessage(err, "failed to process AccrualApproved message")
	}

	return nil
}
