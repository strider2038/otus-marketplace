package messaging

import (
	"context"

	"trading-service/internal/trading"

	"github.com/pkg/errors"
)

type Message interface {
	Name() string
}

type Dispatcher interface {
	Dispatch(ctx context.Context, message Message) error
}

type BillingAdapter struct {
	dispatcher Dispatcher
}

func NewBillingAdapter(dispatcher Dispatcher) *BillingAdapter {
	return &BillingAdapter{dispatcher: dispatcher}
}

func (billing *BillingAdapter) MakePayment(ctx context.Context, payment *trading.Payment) error {
	err := billing.dispatcher.Dispatch(ctx, CreatePayment{
		ID:          payment.ID,
		UserID:      payment.UserID,
		Amount:      payment.Amount,
		Commission:  payment.Commission,
		Description: payment.Description,
	})
	if err != nil {
		return errors.WithMessage(err, "failed to dispatch CreatePayment message")
	}

	return nil
}

func (billing *BillingAdapter) MakeAccrual(ctx context.Context, accrual *trading.Accrual) error {
	err := billing.dispatcher.Dispatch(ctx, CreateAccrual{
		ID:          accrual.ID,
		UserID:      accrual.UserID,
		Amount:      accrual.Amount,
		Commission:  accrual.Commission,
		Description: accrual.Description,
	})
	if err != nil {
		return errors.WithMessage(err, "failed to dispatch CreateAccrual message")
	}

	return nil
}

type TradingAdapter struct {
	dispatcher Dispatcher
}

func NewTradingAdapter(dispatcher Dispatcher) *TradingAdapter {
	return &TradingAdapter{dispatcher: dispatcher}
}

func (trading *TradingAdapter) Publish(ctx context.Context, event trading.Event) error {
	err := trading.dispatcher.Dispatch(ctx, event)
	if err != nil {
		return errors.WithMessagef(err, "failed to dispatch %s message", event.Name())
	}

	return nil
}
