/*
 * Public API
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package api

import (
	"context"
	"net/http"
	"time"

	"trading-service/internal/trading"

	"github.com/bsm/redislock"
	"github.com/gofrs/uuid"
	"github.com/muonsoft/validation"
	"github.com/pkg/errors"
	"github.com/strider2038/pkg/persistence"
)

// TradingApiService is a service that implents the logic for the TradingApiServicer
// This service should implement the business logic for every endpoint for the TradingApi API.
// Include any external packages or services that will be required by this service.
type TradingApiService struct {
	purchaseOrders      trading.PurchaseOrderRepository
	sellOrders          trading.SellOrderRepository
	items               trading.ItemRepository
	userItems           trading.UserItemRepository
	aggregatedUserItems trading.AggregatedUserItemRepository
	transactionManager  persistence.TransactionManager
	dealer              *trading.Dealer
	validator           *validation.Validator
	purchaseState       StateRepository
	sellState           StateRepository
	locker              Locker
	lockTimeout         time.Duration
}

// NewTradingApiService creates a default api service
func NewTradingApiService(
	purchaseOrders trading.PurchaseOrderRepository,
	sellOrders trading.SellOrderRepository,
	items trading.ItemRepository,
	userItems trading.UserItemRepository,
	aggregatedUserItems trading.AggregatedUserItemRepository,
	transactionManager persistence.TransactionManager,
	dealer *trading.Dealer,
	validator *validation.Validator,
	purchaseState StateRepository,
	sellState StateRepository,
	locker Locker,
	lockTimeout time.Duration,
) TradingApiServicer {
	return &TradingApiService{
		purchaseOrders:      purchaseOrders,
		sellOrders:          sellOrders,
		items:               items,
		userItems:           userItems,
		aggregatedUserItems: aggregatedUserItems,
		transactionManager:  transactionManager,
		dealer:              dealer,
		validator:           validator,
		purchaseState:       purchaseState,
		sellState:           sellState,
		locker:              locker,
		lockTimeout:         lockTimeout,
	}
}

// GetPurchaseOrders -
func (s *TradingApiService) GetPurchaseOrders(ctx context.Context, userID uuid.UUID) (ImplResponse, string, error) {
	if userID.IsNil() {
		return newUnauthorizedResponse(), "", nil
	}

	orders, err := s.purchaseOrders.FindByUser(ctx, userID)
	if err != nil {
		return Response(http.StatusInternalServerError, nil), "", err
	}

	state, err := s.purchaseState.Get(ctx, userID)
	if err != nil {
		return Response(http.StatusInternalServerError, nil), "", errors.WithMessagef(
			err,
			`failed to get purchase state of user "%s"`,
			userID,
		)
	}

	return Response(http.StatusOK, orders), state, nil
}

// CreatePurchaseOrder -
func (s *TradingApiService) CreatePurchaseOrder(ctx context.Context, form PurchaseOrder) (ImplResponse, error) {
	if form.UserID.IsNil() {
		return newUnauthorizedResponse(), nil
	}
	if form.IdempotenceKey == "" {
		return newPreconditionRequiredResponse(), nil
	}
	err := s.validator.ValidateValidatable(ctx, form)
	if err != nil {
		return newUnprocessableEntityResponse(err.Error()), nil
	}

	lock, err := s.locker.Obtain(ctx, form.UserID.String()+":create-purchase-order", s.lockTimeout)
	if errors.Is(err, redislock.ErrNotObtained) {
		return newConflictResponse(), nil
	}
	if err != nil {
		return newErrorResponsef(err, "failed to obtain lock")
	}
	defer lock.Release(ctx)

	err = s.purchaseState.Verify(ctx, form.UserID, form.IdempotenceKey)
	if errors.Is(err, trading.ErrOutdated) {
		return newPreconditionFailedResponse(), nil
	}
	if err != nil {
		return newErrorResponsef(err, "failed to verify idempotence key")
	}

	item, err := s.items.FindByID(ctx, form.ItemID)
	if errors.Is(err, trading.ErrItemNotFound) {
		return newUnprocessableEntityResponse("trading item does not exist"), nil
	}
	if err != nil {
		return newErrorResponsef(err, "failed to create purchase order")
	}

	order := trading.NewPurchaseOrder(form.UserID, item, form.Price)
	err = s.dealer.CreatePurchaseOrder(ctx, item, order)
	if err != nil {
		return newErrorResponsef(err, "failed to create purchase order")
	}

	err = s.purchaseState.Refresh(ctx, form.UserID)
	if err != nil {
		return newErrorResponsef(err, `failed to refresh purchase state of user "%s"`, form.UserID)
	}

	return Response(http.StatusAccepted, order), nil
}

// CancelPurchaseOrder -
func (s *TradingApiService) CancelPurchaseOrder(ctx context.Context, userID, orderID uuid.UUID) (ImplResponse, error) {
	if userID.IsNil() {
		return newUnauthorizedResponse(), nil
	}
	var order *trading.PurchaseOrder
	var err error

	err = s.transactionManager.DoTransactionally(ctx, func(ctx context.Context) error {
		order, err = s.purchaseOrders.FindByIDForUpdate(ctx, orderID)
		if err != nil {
			return errors.WithMessage(err, "failed to find order")
		}

		err = order.CancelByUser(userID)
		if err != nil {
			return errors.WithMessage(err, "failed to cancel purchase order")
		}

		err = s.purchaseOrders.Delete(ctx, order.ID)
		if err != nil {
			return errors.WithMessage(err, "failed to delete purchase order")
		}

		return nil
	})
	if errors.Is(err, trading.ErrOrderNotFound) {
		return newNotFoundResponse(), nil
	}
	if errors.Is(err, trading.ErrDenied) {
		return newForbiddenResponse(), nil
	}
	if errors.Is(err, trading.ErrCannotCancel) {
		return newUnprocessableEntityResponse(err.Error()), nil
	}
	if err != nil {
		return Response(http.StatusInternalServerError, nil), err
	}

	return Response(http.StatusOK, order), nil
}

// GetSellOrders -
func (s *TradingApiService) GetSellOrders(ctx context.Context, userID uuid.UUID) (ImplResponse, string, error) {
	if userID.IsNil() {
		return newUnauthorizedResponse(), "", nil
	}

	orders, err := s.sellOrders.FindByUser(ctx, userID)
	if err != nil {
		return Response(http.StatusInternalServerError, nil), "", err
	}

	state, err := s.sellState.Get(ctx, userID)
	if err != nil {
		return Response(http.StatusInternalServerError, nil), "", errors.WithMessagef(
			err,
			`failed to get sell state of user "%s"`,
			userID,
		)
	}

	return Response(http.StatusOK, orders), state, nil
}

// CreateSellOrder -
func (s *TradingApiService) CreateSellOrder(ctx context.Context, form SellOrder) (ImplResponse, error) {
	if form.UserID.IsNil() {
		return newUnauthorizedResponse(), nil
	}
	if form.IdempotenceKey == "" {
		return newPreconditionRequiredResponse(), nil
	}
	err := s.validator.ValidateValidatable(ctx, form)
	if err != nil {
		return newUnprocessableEntityResponse(err.Error()), nil
	}

	lock, err := s.locker.Obtain(ctx, form.UserID.String()+":create-sell-order", s.lockTimeout)
	if errors.Is(err, redislock.ErrNotObtained) {
		return newConflictResponse(), nil
	}
	if err != nil {
		return newErrorResponsef(err, "failed to obtain lock")
	}
	defer lock.Release(ctx)

	err = s.sellState.Verify(ctx, form.UserID, form.IdempotenceKey)
	if errors.Is(err, trading.ErrOutdated) {
		return newPreconditionFailedResponse(), nil
	}
	if err != nil {
		return newErrorResponsef(err, "failed to verify idempotence key")
	}

	item, err := s.items.FindByID(ctx, form.ItemID)
	if errors.Is(err, trading.ErrItemNotFound) {
		return newUnprocessableEntityResponse("trading item does not exist"), nil
	}
	if err != nil {
		return Response(http.StatusInternalServerError, nil), err
	}

	sellOrder, err := s.dealer.CreateSellOrder(ctx, item, form.UserID, form.Price)
	if errors.Is(err, trading.ErrItemNotFound) {
		return newUnprocessableEntityResponse("no trading items for sale"), nil
	}
	if err != nil {
		return Response(http.StatusInternalServerError, nil), err
	}

	err = s.sellState.Refresh(ctx, form.UserID)
	if err != nil {
		return newErrorResponsef(err, `failed to refresh sell state of user "%s"`, form.UserID)
	}

	return Response(http.StatusAccepted, sellOrder), nil
}

// CancelSellOrder -
func (s *TradingApiService) CancelSellOrder(ctx context.Context, userID, orderID uuid.UUID) (ImplResponse, error) {
	if userID.IsNil() {
		return newUnauthorizedResponse(), nil
	}
	var order *trading.SellOrder
	var err error

	err = s.transactionManager.DoTransactionally(ctx, func(ctx context.Context) error {
		order, err = s.sellOrders.FindByIDForUpdate(ctx, orderID)
		if err != nil {
			return errors.WithMessage(err, "failed to find order")
		}

		err = order.CancelByUser(userID)
		if err != nil {
			return errors.WithMessage(err, "failed to cancel sell order")
		}

		err = s.sellOrders.Delete(ctx, order.ID)
		if err != nil {
			return errors.WithMessage(err, "failed to delete sell order")
		}

		return nil
	})
	if errors.Is(err, trading.ErrOrderNotFound) {
		return newNotFoundResponse(), nil
	}
	if errors.Is(err, trading.ErrDenied) {
		return newForbiddenResponse(), nil
	}
	if errors.Is(err, trading.ErrCannotCancel) {
		return newUnprocessableEntityResponse(err.Error()), nil
	}
	if err != nil {
		return Response(http.StatusInternalServerError, nil), err
	}

	return Response(http.StatusOK, order), nil
}

// GetTradingItems -
func (s *TradingApiService) GetTradingItems(ctx context.Context) (ImplResponse, error) {
	items, err := s.items.FindAll(ctx)
	if err != nil {
		return Response(http.StatusInternalServerError, nil), err
	}

	return Response(http.StatusOK, items), nil
}

// CreateTradingItem -
func (s *TradingApiService) CreateTradingItem(ctx context.Context, form TradingItem) (ImplResponse, error) {
	if form.UserID.IsNil() {
		return newUnauthorizedResponse(), nil
	}
	if form.UserRole != "broker" {
		return newForbiddenResponse(), nil
	}

	err := s.validator.ValidateValidatable(ctx, form)
	if err != nil {
		return newUnprocessableEntityResponse(err.Error()), nil
	}

	item := trading.NewItem(form.Name, form.InitialCount, form.InitialPrice, form.CommissionPercent)
	err = s.transactionManager.DoTransactionally(ctx, func(ctx context.Context) error {
		err := s.items.Add(ctx, item)
		if err != nil {
			return errors.WithMessagef(err, `failed to add item "%s"`, item.Name)
		}

		for i := 0; i < int(item.InitialCount); i++ {
			userItem := trading.NewUserItem(item.ID, form.UserID)
			userItem.IsOnSale = true
			err = s.userItems.Save(ctx, userItem)
			if err != nil {
				return errors.WithMessagef(err, `failed to save user item "%s"`, item.Name)
			}

			err = s.sellOrders.Save(ctx, trading.NewInitialOrder(form.UserID, item, userItem))
			if err != nil {
				return errors.WithMessage(err, "failed to save initial order")
			}
		}

		return nil
	})
	if err != nil {
		return Response(http.StatusInternalServerError, nil), err
	}

	return Response(http.StatusCreated, item), nil
}

func (s *TradingApiService) GetUserItems(ctx context.Context, userID uuid.UUID) (ImplResponse, error) {
	if userID.IsNil() {
		return newUnauthorizedResponse(), nil
	}

	items, err := s.aggregatedUserItems.FindByUser(ctx, userID)
	if err != nil {
		return Response(http.StatusInternalServerError, nil), err
	}

	return Response(http.StatusOK, items), nil
}
