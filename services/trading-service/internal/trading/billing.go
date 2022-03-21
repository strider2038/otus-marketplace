package trading

import (
	"context"

	"github.com/gofrs/uuid"
)

type Payment struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Amount      float64
	Commission  float64
	Description string
}

func NewPayment(userID uuid.UUID, amount, commission float64, description string) *Payment {
	return &Payment{
		ID:          uuid.Must(uuid.NewV4()),
		UserID:      userID,
		Amount:      amount,
		Commission:  commission,
		Description: description,
	}
}

type Accrual struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Amount      float64
	Commission  float64
	Description string
}

func NewAccrual(userID uuid.UUID, amount, commission float64, description string) *Accrual {
	return &Accrual{
		ID:          uuid.Must(uuid.NewV4()),
		UserID:      userID,
		Amount:      amount,
		Commission:  commission,
		Description: description,
	}
}

type Billing interface {
	MakePayment(ctx context.Context, payment *Payment) error
	MakeAccrual(ctx context.Context, accrual *Accrual) error
}
