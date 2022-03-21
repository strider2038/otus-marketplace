package messaging

import (
	"context"
	"encoding/json"

	"trading-service/internal/trading"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

type CreatePayment struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"userId"`
	Amount      float64   `json:"amount"`
	Commission  float64   `json:"commission"`
	Description string    `json:"description"`
}

func (p CreatePayment) Name() string {
	return "Billing/CreatePayment"
}

type PaymentSucceeded struct {
	ID         uuid.UUID `json:"id"`
	Amount     float64   `json:"amount"`
	Commission float64   `json:"commission"`
}

func (p PaymentSucceeded) Name() string {
	return "Billing/PaymentSucceeded"
}

type PaymentDeclined struct {
	ID         uuid.UUID `json:"id"`
	Amount     float64   `json:"amount"`
	Commission float64   `json:"commission"`
	Reason     string    `json:"reason"`
}

func (p PaymentDeclined) Name() string {
	return "Billing/PaymentDeclined"
}

type PaymentSucceededProcessor struct {
	dealer *trading.Dealer
}

func NewPaymentSucceededProcessor(dealer *trading.Dealer) *PaymentSucceededProcessor {
	return &PaymentSucceededProcessor{dealer: dealer}
}

func (p *PaymentSucceededProcessor) Process(ctx context.Context, message []byte) error {
	var payment PaymentSucceeded
	err := json.Unmarshal(message, &payment)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal PaymentSucceeded message")
	}

	err = p.dealer.ApprovePayment(ctx, &trading.Payment{
		ID:         payment.ID,
		Amount:     payment.Amount,
		Commission: payment.Commission,
	})
	if err != nil {
		return errors.WithMessage(err, "failed to process PaymentSucceeded message")
	}

	return nil
}

type PaymentDeclinedProcessor struct {
	dealer *trading.Dealer
}

func NewPaymentDeclinedProcessor(dealer *trading.Dealer) *PaymentDeclinedProcessor {
	return &PaymentDeclinedProcessor{dealer: dealer}
}

func (p *PaymentDeclinedProcessor) Process(ctx context.Context, message []byte) error {
	var payment PaymentDeclined
	err := json.Unmarshal(message, &payment)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal PaymentDeclined message")
	}

	err = p.dealer.DeclinePayment(ctx, &trading.Payment{ID: payment.ID}, payment.Reason)
	if err != nil {
		return errors.WithMessage(err, "failed to process PaymentDeclined message")
	}

	return nil
}
