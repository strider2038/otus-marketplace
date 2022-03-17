package billing

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
)

type OperationType string

const (
	DepositOperation    OperationType = "deposit"
	WithdrawOperation   OperationType = "withdraw"
	PaymentOperation    OperationType = "payment"
	AccrualOperation    OperationType = "accrual"
	CommissionOperation OperationType = "commission"
)

type Operation struct {
	ID          uuid.UUID     `json:"id"`
	AccountID   uuid.UUID     `json:"-"`
	Type        OperationType `json:"type"`
	Amount      float64       `json:"amount"`
	Description string        `json:"description"`
	CreatedAt   time.Time     `json:"createdAt"`
}

func NewDeposit(accountID uuid.UUID, amount float64, description string) *Operation {
	return newOperation(uuid.Must(uuid.NewV4()), accountID, DepositOperation, amount, description)
}

func NewWithdrawal(accountID uuid.UUID, amount float64, description string) *Operation {
	return newOperation(uuid.Must(uuid.NewV4()), accountID, WithdrawOperation, amount, description)
}

func NewPayment(id, accountID uuid.UUID, amount float64, description string) *Operation {
	return newOperation(id, accountID, PaymentOperation, amount, description)
}

func NewAccrual(id, accountID uuid.UUID, amount float64, description string) *Operation {
	return newOperation(id, accountID, AccrualOperation, amount, description)
}

func NewCommission(accountID uuid.UUID, amount float64, description string) *Operation {
	return newOperation(uuid.Must(uuid.NewV4()), accountID, CommissionOperation, amount, description)
}

type OperationRepository interface {
	FindByAccount(ctx context.Context, accountID uuid.UUID) ([]*Operation, error)
	Add(ctx context.Context, operation *Operation) error
}

func newOperation(id, accountID uuid.UUID, opType OperationType, amount float64, description string) *Operation {
	return &Operation{
		ID:          id,
		AccountID:   accountID,
		Type:        opType,
		Amount:      amount,
		Description: description,
		CreatedAt:   time.Now(),
	}
}
