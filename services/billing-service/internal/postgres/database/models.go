// Code generated by sqlc. DO NOT EDIT.

package database

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"
)

type OperationType string

const (
	OperationTypeDeposit  OperationType = "deposit"
	OperationTypeWithdraw OperationType = "withdraw"
	OperationTypePayment  OperationType = "payment"
	OperationTypeAccrual  OperationType = "accrual"
)

func (e *OperationType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = OperationType(s)
	case string:
		*e = OperationType(s)
	default:
		return fmt.Errorf("unsupported scan type for OperationType: %T", src)
	}
	return nil
}

type Account struct {
	ID        uuid.UUID
	Amount    float64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Operation struct {
	ID          uuid.UUID
	AccountID   uuid.UUID
	Type        OperationType
	Amount      float64
	Description string
	CreatedAt   time.Time
}
