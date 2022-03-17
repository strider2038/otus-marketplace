package trading

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

type SellStatus string

const (
	SellPending        SellStatus = "pending"
	SellCanceled       SellStatus = "canceled"
	SellDealPending    SellStatus = "dealPending"
	SellAccrualPending SellStatus = "accrualPending"
	SellApproved       SellStatus = "approved"
)

type SellOrder struct {
	ID         uuid.UUID     `json:"id"`
	UserID     uuid.NullUUID `json:"-"`
	ItemID     uuid.UUID     `json:"itemId,omitempty"`
	AccrualID  uuid.NullUUID `json:"-"`
	DealID     uuid.NullUUID `json:"-"`
	Price      float64       `json:"price"`
	Commission float64       `json:"commission"`
	Status     SellStatus    `json:"status,omitempty"`
	CreatedAt  time.Time     `json:"createdAt,omitempty"`
	UpdatedAt  time.Time     `json:"updatedAt,omitempty"`

	isNew bool
}

func (order SellOrder) IsNew() bool {
	return order.isNew
}

func (order SellOrder) TotalPrice() float64 {
	return order.Price + order.Commission
}

func (order *SellOrder) InitiateDeal(dealID uuid.UUID) error {
	err := order.verifyStatus(SellPending)
	if err != nil {
		return err
	}

	order.DealID = uuid.NullUUID{UUID: dealID, Valid: true}
	order.Status = SellDealPending

	return nil
}

func (order *SellOrder) CancelByUser(userID uuid.UUID) error {
	if order.UserID.Valid && order.UserID.UUID != userID {
		return errors.WithStack(ErrDenied)
	}
	if order.Status != SellPending {
		return errors.WithStack(ErrCannotCancel)
	}

	order.Status = SellCanceled

	return nil
}

func (order *SellOrder) CancelDeal() error {
	err := order.verifyStatus(SellDealPending)
	if err != nil {
		return err
	}

	order.Status = SellPending
	order.DealID = uuid.NullUUID{}

	return nil
}

func (order *SellOrder) InitiateAccrual(accrualID uuid.UUID) error {
	err := order.verifyStatus(SellDealPending)
	if err != nil {
		return err
	}

	order.AccrualID = uuid.NullUUID{UUID: accrualID, Valid: true}
	order.Status = SellAccrualPending

	return nil
}

func (order *SellOrder) Approve() error {
	err := order.verifyStatus(SellAccrualPending)
	if err != nil {
		return err
	}

	order.Status = SellApproved

	return nil
}

func (order *SellOrder) verifyStatus(status SellStatus) error {
	if order.Status != status {
		return newUnexpectedStatusError(string(order.Status), string(status))
	}

	return nil
}

func NewSellOrder(userID uuid.UUID, item *Item, price float64) *SellOrder {
	now := time.Now()

	return &SellOrder{
		isNew:      true,
		ID:         uuid.Must(uuid.NewV4()),
		UserID:     uuid.NullUUID{UUID: userID, Valid: true},
		ItemID:     item.ID,
		Price:      price,
		Commission: item.CalculateCommission(price),
		Status:     SellPending,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

func NewInitialOrder(itemID uuid.UUID, price, commission float64) *SellOrder {
	now := time.Now()

	return &SellOrder{
		ID:         uuid.Must(uuid.NewV4()),
		ItemID:     itemID,
		Price:      price,
		Commission: commission,
		Status:     SellPending,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

type SellOrderRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*SellOrder, error)
	FindByIDForUpdate(ctx context.Context, id uuid.UUID) (*SellOrder, error)
	FindByAccrualForUpdate(ctx context.Context, accrualID uuid.UUID) (*SellOrder, error)
	FindByDealForUpdate(ctx context.Context, dealID uuid.UUID) (*SellOrder, error)
	FindByUser(ctx context.Context, userID uuid.UUID) ([]*SellOrder, error)
	FindForDeal(ctx context.Context, itemID uuid.UUID, maxPrice float64) (*SellOrder, error)
	Save(ctx context.Context, order *SellOrder) error
}
