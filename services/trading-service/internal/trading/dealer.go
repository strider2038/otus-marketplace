package trading

import (
	"context"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/strider2038/pkg/persistence"
)

type Dealer struct {
	items              ItemRepository
	userItems          UserItemRepository
	purchaseOrders     PurchaseOrderRepository
	sellOrders         SellOrderRepository
	transactionManager persistence.TransactionManager
	billing            Billing
	events             EventPublisher
}

func NewDealer(
	items ItemRepository,
	userItems UserItemRepository,
	purchaseOrders PurchaseOrderRepository,
	sellOrders SellOrderRepository,
	transactionManager persistence.TransactionManager,
	billing Billing,
	events EventPublisher,
) *Dealer {
	return &Dealer{
		items:              items,
		userItems:          userItems,
		purchaseOrders:     purchaseOrders,
		sellOrders:         sellOrders,
		transactionManager: transactionManager,
		billing:            billing,
		events:             events,
	}
}

func (dealer *Dealer) CreatePurchaseOrder(ctx context.Context, item *Item, purchaseOrder *PurchaseOrder) error {
	return dealer.transactionManager.DoTransactionally(ctx, func(ctx context.Context) error {
		return dealer.createPurchaseOrder(ctx, item, purchaseOrder)
	})
}

func (dealer *Dealer) CreateSellOrder(ctx context.Context, item *Item, userID uuid.UUID, price float64) (*SellOrder, error) {
	var sellOrder *SellOrder
	var err error

	err = dealer.transactionManager.DoTransactionally(ctx, func(ctx context.Context) error {
		sellOrder, err = dealer.createSellOrder(ctx, item, userID, price)

		return err
	})
	if err != nil {
		return nil, err
	}

	return sellOrder, nil
}

func (dealer *Dealer) ApprovePayment(ctx context.Context, payment *Payment) error {
	return dealer.transactionManager.DoTransactionally(ctx, func(ctx context.Context) error {
		return dealer.approvePayment(ctx, payment)
	})
}

func (dealer *Dealer) DeclinePayment(ctx context.Context, payment *Payment, reason string) error {
	return dealer.transactionManager.DoTransactionally(ctx, func(ctx context.Context) error {
		return dealer.declinePayment(ctx, payment, reason)
	})
}

func (dealer *Dealer) ApproveAccrual(ctx context.Context, accrual *Accrual) error {
	return dealer.transactionManager.DoTransactionally(ctx, func(ctx context.Context) error {
		return dealer.approveAccrual(ctx, accrual)
	})
}

func (dealer *Dealer) createPurchaseOrder(ctx context.Context, item *Item, purchaseOrder *PurchaseOrder) error {
	sellOrder, err := dealer.sellOrders.FindForDeal(ctx, purchaseOrder.UserID, purchaseOrder.ItemID, purchaseOrder.Price)
	if errors.Is(err, ErrOrderNotFound) {
		err = dealer.purchaseOrders.Save(ctx, purchaseOrder)
		if err != nil {
			return errors.WithMessagef(err, "failed to save purchase order of user %s", purchaseOrder.UserID)
		}

		return nil
	}
	if err != nil {
		return errors.WithMessage(err, "failed to find sell order")
	}

	err = dealer.startDeal(ctx, sellOrder.Price, sellOrder.Commission, purchaseOrder, sellOrder, item)
	if err != nil {
		return errors.WithMessage(err, "failed to start deal")
	}

	return nil
}

func (dealer *Dealer) createSellOrder(ctx context.Context, item *Item, userID uuid.UUID, price float64) (*SellOrder, error) {
	userItem, err := dealer.userItems.FindForSale(ctx, userID)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to find user item for sale")
	}

	sellOrder, err := NewSellOrder(userID, item, userItem, price)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create sell order")
	}

	purchaseOrder, err := dealer.purchaseOrders.FindForDeal(ctx, sellOrder.UserID, sellOrder.ItemID, sellOrder.Price)
	if errors.Is(err, ErrOrderNotFound) {
		err = dealer.sellOrders.Save(ctx, sellOrder)
		if err != nil {
			return nil, errors.WithMessagef(err, "failed to save sell order of user %s", sellOrder.UserID)
		}

		userItem.IsOnSale = true
		err = dealer.userItems.Save(ctx, userItem)
		if err != nil {
			return nil, errors.WithMessage(err, "failed to save user item")
		}

		return sellOrder, nil
	}
	if err != nil {
		return nil, errors.WithMessage(err, "failed to find purchase order")
	}

	err = dealer.startDeal(ctx, purchaseOrder.Price, purchaseOrder.Commission, purchaseOrder, sellOrder, item)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to start deal")
	}

	userItem.IsOnSale = true
	err = dealer.userItems.Save(ctx, userItem)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to save user item")
	}

	return sellOrder, nil
}

func (dealer *Dealer) startDeal(
	ctx context.Context,
	amount float64,
	commission float64,
	purchaseOrder *PurchaseOrder,
	sellOrder *SellOrder,
	item *Item,
) error {
	payment := NewPayment(
		purchaseOrder.UserID,
		amount,
		commission,
		fmt.Sprintf(`buying the item "%s" on the marketplace (with commission %2.f)`, item.Name, commission),
	)

	dealID := uuid.Must(uuid.NewV4())
	err := purchaseOrder.InitiatePayment(dealID, payment.ID)
	if err != nil {
		return errors.WithMessagef(err, "failed to initiate payment for purchase order %s", purchaseOrder.ID)
	}
	err = sellOrder.InitiateDeal(dealID)
	if err != nil {
		return errors.WithMessagef(err, "failed to initiate deal for sell order %s", sellOrder.ID)
	}

	err = dealer.saveOrders(ctx, sellOrder, purchaseOrder)
	if err != nil {
		return err
	}

	err = dealer.billing.MakePayment(ctx, payment)
	if err != nil {
		return errors.WithMessage(err, "failed to make payment")
	}

	return nil
}

func (dealer *Dealer) approvePayment(ctx context.Context, payment *Payment) error {
	purchaseOrder, err := dealer.purchaseOrders.FindByPaymentForUpdate(ctx, payment.ID)
	if err != nil {
		return errors.WithMessagef(err, "failed to find purchase order by payment %s", payment.ID)
	}

	sellOrder, err := dealer.sellOrders.FindByDealForUpdate(ctx, purchaseOrder.DealID.UUID)
	if err != nil {
		return errors.WithMessagef(err, "failed to find purchase order by deal %s", purchaseOrder.DealID.UUID)
	}

	item, err := dealer.items.FindByID(ctx, purchaseOrder.ItemID)
	if err != nil {
		return errors.WithMessagef(err, "failed to find item %s", item.ID)
	}

	err = purchaseOrder.ApprovePayment()
	if err != nil {
		return errors.WithMessagef(err, "failed to approve payment for purchase order %s", purchaseOrder.ID)
	}

	accrual := NewAccrual(
		sellOrder.UserID,
		payment.Amount,
		payment.Commission,
		fmt.Sprintf(`selling the item "%s" on the marketplace (with commission %.2f)`, item.Name, payment.Commission),
	)
	err = sellOrder.InitiateAccrual(accrual.ID)
	if err != nil {
		return errors.WithMessagef(err, "failed to initiate accrual for sell order %s", sellOrder.ID)
	}

	err = dealer.saveOrders(ctx, sellOrder, purchaseOrder)
	if err != nil {
		return err
	}

	err = dealer.billing.MakeAccrual(ctx, accrual)
	if err != nil {
		return errors.WithMessage(err, "failed to make accrual")
	}

	return nil
}

func (dealer *Dealer) declinePayment(ctx context.Context, payment *Payment, reason string) error {
	purchaseOrder, err := dealer.purchaseOrders.FindByPaymentForUpdate(ctx, payment.ID)
	if err != nil {
		return errors.WithMessagef(err, "failed to find purchase order by payment %s", payment.ID)
	}

	sellOrder, err := dealer.sellOrders.FindByDealForUpdate(ctx, purchaseOrder.DealID.UUID)
	if err != nil {
		return errors.WithMessagef(err, "failed to find purchase order by deal %s", purchaseOrder.DealID.UUID)
	}

	err = purchaseOrder.DeclinePayment()
	if err != nil {
		return errors.WithMessagef(err, "failed to decline payment for purchase order %s", purchaseOrder.ID)
	}
	err = sellOrder.CancelDeal()
	if err != nil {
		return errors.WithMessagef(err, "failed to cancel deal for sell order %s", sellOrder.ID)
	}

	err = dealer.sellOrders.Save(ctx, sellOrder)
	if err != nil {
		return errors.WithMessagef(err, "failed to save sell order %s", sellOrder.ID)
	}

	err = dealer.purchaseOrders.Delete(ctx, purchaseOrder.ID)
	if err != nil {
		return errors.WithMessagef(err, "failed to delete purchase order %s", purchaseOrder.ID)
	}

	item, err := dealer.items.FindByID(ctx, purchaseOrder.ItemID)
	if err != nil {
		return errors.WithMessagef(err, "failed to find item for purchase order %s", purchaseOrder.ID)
	}

	err = dealer.events.Publish(ctx, PurchaseFailedEvent{
		PurchaserID: purchaseOrder.UserID,
		ItemID:      item.ID,
		ItemName:    item.Name,
		Amount:      payment.Amount,
		Reason:      reason,
	})

	return nil
}

func (dealer *Dealer) approveAccrual(ctx context.Context, accrual *Accrual) error {
	sellOrder, err := dealer.sellOrders.FindByAccrualForUpdate(ctx, accrual.ID)
	if err != nil {
		return errors.WithMessagef(err, "failed to find sell order by accrual %s", accrual.ID)
	}

	soldItem, err := dealer.userItems.FindByIDForUpdate(ctx, sellOrder.UserItemID)
	if err != nil {
		return errors.WithMessagef(err, "failed to find sold item for sell order %s", sellOrder.ID)
	}

	purchaseOrder, err := dealer.purchaseOrders.FindByDealForUpdate(ctx, sellOrder.DealID.UUID)
	if err != nil {
		return errors.WithMessagef(err, "failed to find purchase order by deal %s", sellOrder.DealID.UUID)
	}

	err = purchaseOrder.Approve()
	if err != nil {
		return errors.WithMessagef(err, "failed to approve purchase order %s", purchaseOrder.ID)
	}
	err = sellOrder.Approve()
	if err != nil {
		return errors.WithMessagef(err, "failed to approve sell order %s", sellOrder.ID)
	}

	err = dealer.sellOrders.Delete(ctx, sellOrder.ID)
	if err != nil {
		return errors.WithMessagef(err, "failed to delete sell order %s", sellOrder.ID)
	}

	err = dealer.purchaseOrders.Delete(ctx, purchaseOrder.ID)
	if err != nil {
		return errors.WithMessagef(err, "failed to delete purchase order %s", purchaseOrder.ID)
	}

	err = dealer.userItems.Delete(ctx, soldItem.ID)
	if err != nil {
		return errors.WithMessagef(err, "failed to delete sold item %s", soldItem.ID)
	}

	purchasedItem := NewUserItem(purchaseOrder.ItemID, purchaseOrder.UserID)
	err = dealer.userItems.Save(ctx, purchasedItem)
	if err != nil {
		return errors.WithMessagef(err, "failed to save purchased item %s for user %s", purchaseOrder.ItemID, purchaseOrder.UserID)
	}

	item, err := dealer.items.FindByID(ctx, purchaseOrder.ItemID)
	if err != nil {
		return errors.WithMessagef(err, "failed to find item for purchase order %s", purchaseOrder.ID)
	}

	err = dealer.events.Publish(ctx, DealSucceededEvent{
		SellerID:            sellOrder.UserID,
		PurchaserID:         purchaseOrder.UserID,
		ItemID:              item.ID,
		ItemName:            item.Name,
		Amount:              accrual.Amount,
		SellerCommission:    sellOrder.Commission,
		PurchaserCommission: purchaseOrder.Commission,
		CompletedAt:         time.Now(),
	})

	return nil
}

func (dealer *Dealer) saveOrders(ctx context.Context, sellOrder *SellOrder, purchaseOrder *PurchaseOrder) error {
	err := dealer.sellOrders.Save(ctx, sellOrder)
	if err != nil {
		return errors.WithMessagef(err, "failed to save sell order %s", sellOrder.ID)
	}

	err = dealer.purchaseOrders.Save(ctx, purchaseOrder)
	if err != nil {
		return errors.WithMessagef(err, "failed to save purchase order of user %s", purchaseOrder.UserID)
	}

	return nil
}
