package trading

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

type PurchaseStatus string

const (
	PurchasePending          PurchaseStatus = "pending"
	PurchaseCanceled         PurchaseStatus = "canceled"
	PurchasePaymentPending   PurchaseStatus = "paymentPending"
	PurchasePaymentSucceeded PurchaseStatus = "paymentSucceeded"
	PurchasePaymentFailed    PurchaseStatus = "paymentFailed"
	PurchaseApproved         PurchaseStatus = "approved"
)

type PurchaseOrder struct {
	ID         uuid.UUID      `json:"id"`
	UserID     uuid.UUID      `json:"-"`
	ItemID     uuid.UUID      `json:"itemId,omitempty"`
	PaymentID  uuid.NullUUID  `json:"-"`
	DealID     uuid.NullUUID  `json:"-"`
	Price      float64        `json:"price"`
	Commission float64        `json:"commission"`
	Status     PurchaseStatus `json:"status,omitempty"`
	CreatedAt  time.Time      `json:"createdAt,omitempty"`
	UpdatedAt  time.Time      `json:"updatedAt,omitempty"`

	isNew bool
}

func (order PurchaseOrder) IsNew() bool {
	return order.isNew
}

func (order PurchaseOrder) TotalPrice() float64 {
	return order.Price + order.Commission
}

func (order *PurchaseOrder) InitiatePayment(dealID, paymentID uuid.UUID) error {
	err := order.verifyStatus(PurchasePending)
	if err != nil {
		return err
	}

	order.DealID = uuid.NullUUID{UUID: dealID, Valid: true}
	order.PaymentID = uuid.NullUUID{UUID: paymentID, Valid: true}
	order.Status = PurchasePaymentPending

	return nil
}

func (order *PurchaseOrder) ApprovePayment() error {
	err := order.verifyStatus(PurchasePaymentPending)
	if err != nil {
		return err
	}

	order.Status = PurchasePaymentSucceeded

	return nil
}

func (order *PurchaseOrder) DeclinePayment() error {
	err := order.verifyStatus(PurchasePaymentPending)
	if err != nil {
		return err
	}

	order.Status = PurchasePaymentFailed

	return nil
}

func (order *PurchaseOrder) CancelByUser(userID uuid.UUID) error {
	if order.UserID != userID {
		return errors.WithStack(ErrDenied)
	}
	if order.Status != PurchasePending {
		return errors.WithStack(ErrCannotCancel)
	}

	order.Status = PurchaseCanceled

	return nil
}

func (order *PurchaseOrder) Approve() error {
	err := order.verifyStatus(PurchasePaymentSucceeded)
	if err != nil {
		return err
	}

	order.Status = PurchaseApproved

	return nil
}

func (order *PurchaseOrder) verifyStatus(status PurchaseStatus) error {
	if order.Status != status {
		return newUnexpectedStatusError(string(order.Status), string(status))
	}

	return nil
}

func NewPurchaseOrder(userID uuid.UUID, item *Item, price float64) *PurchaseOrder {
	now := time.Now()

	return &PurchaseOrder{
		isNew:      true,
		ID:         uuid.Must(uuid.NewV4()),
		UserID:     userID,
		ItemID:     item.ID,
		Price:      price,
		Commission: item.CalculateCommission(price),
		Status:     PurchasePending,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

type PurchaseOrderRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*PurchaseOrder, error)
	FindByIDForUpdate(ctx context.Context, id uuid.UUID) (*PurchaseOrder, error)
	FindByPaymentForUpdate(ctx context.Context, paymentID uuid.UUID) (*PurchaseOrder, error)
	FindByDealForUpdate(ctx context.Context, dealID uuid.UUID) (*PurchaseOrder, error)
	FindByUser(ctx context.Context, userID uuid.UUID) ([]*PurchaseOrder, error)
	FindForDeal(ctx context.Context, itemID uuid.UUID, minPrice float64) (*PurchaseOrder, error)
	Save(ctx context.Context, order *PurchaseOrder) error
}
