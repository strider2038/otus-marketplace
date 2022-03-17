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

	"trading-service/internal/messaging"
	"trading-service/internal/trading"

	"github.com/gofrs/uuid"
	"github.com/muonsoft/validation"
	"github.com/pkg/errors"
	"github.com/strider2038/pkg/persistence"
)

// TradingApiService is a service that implents the logic for the TradingApiServicer
// This service should implement the business logic for every endpoint for the TradingApi API.
// Include any external packages or services that will be required by this service.
type TradingApiService struct {
	purchaseOrders     trading.PurchaseOrderRepository
	sellOrders         trading.SellOrderRepository
	items              trading.ItemRepository
	transactionManager persistence.TransactionManager
	dealer             *trading.Dealer
	dispatcher         messaging.Dispatcher
	validator          *validation.Validator
}

// NewTradingApiService creates a default api service
func NewTradingApiService(
	purchaseOrders trading.PurchaseOrderRepository,
	sellOrders trading.SellOrderRepository,
	items trading.ItemRepository,
	transactionManager persistence.TransactionManager,
	dealer *trading.Dealer,
	validator *validation.Validator,
) TradingApiServicer {
	return &TradingApiService{
		purchaseOrders:     purchaseOrders,
		sellOrders:         sellOrders,
		items:              items,
		transactionManager: transactionManager,
		dealer:             dealer,
		validator:          validator,
	}
}

// CreatePurchaseOrder -
func (s *TradingApiService) CreatePurchaseOrder(ctx context.Context, form PurchaseOrder) (ImplResponse, error) {
	if form.UserID == uuid.Nil {
		return newUnauthorizedResponse(), nil
	}
	err := s.validator.ValidateValidatable(ctx, form)
	if err != nil {
		return newUnprocessableEntityResponse(err.Error()), nil
	}

	item, err := s.items.FindByID(ctx, form.ItemID)
	if errors.Is(err, trading.ErrItemNotFound) {
		return newUnprocessableEntityResponse("trading item does not exist"), nil
	}
	if err != nil {
		return Response(http.StatusInternalServerError, nil), err
	}

	order := trading.NewPurchaseOrder(form.UserID, item, form.Price)
	err = s.transactionManager.DoTransactionally(ctx, func(ctx context.Context) error {
		return s.dealer.CreatePurchaseOrder(ctx, item, order)
	})
	if err != nil {
		return Response(http.StatusInternalServerError, nil), err
	}

	return Response(http.StatusAccepted, order), nil
}

// CancelPurchaseOrder -
func (s *TradingApiService) CancelPurchaseOrder(ctx context.Context, userID, orderID uuid.UUID) (ImplResponse, error) {
	if userID == uuid.Nil {
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

		err = s.purchaseOrders.Save(ctx, order)
		if err != nil {
			return errors.WithMessage(err, "failed to save purchase order")
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

// CreateSellOrder -
func (s *TradingApiService) CreateSellOrder(ctx context.Context, form SellOrder) (ImplResponse, error) {
	if form.UserID == uuid.Nil {
		return newUnauthorizedResponse(), nil
	}
	err := s.validator.ValidateValidatable(ctx, form)
	if err != nil {
		return newUnprocessableEntityResponse(err.Error()), nil
	}

	item, err := s.items.FindByID(ctx, form.ItemID)
	if errors.Is(err, trading.ErrItemNotFound) {
		return newUnprocessableEntityResponse("trading item does not exist"), nil
	}
	if err != nil {
		return Response(http.StatusInternalServerError, nil), err
	}

	order := trading.NewSellOrder(form.UserID, item, form.Price)
	err = s.transactionManager.DoTransactionally(ctx, func(ctx context.Context) error {
		return s.dealer.CreateSellOrder(ctx, item, order)
	})
	if err != nil {
		return Response(http.StatusInternalServerError, nil), err
	}

	return Response(http.StatusAccepted, order), nil
}

// CancelSellOrder -
func (s *TradingApiService) CancelSellOrder(ctx context.Context, userID, orderID uuid.UUID) (ImplResponse, error) {
	if userID == uuid.Nil {
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
			return errors.WithMessage(err, "failed to cancel purchase order")
		}

		err = s.sellOrders.Save(ctx, order)
		if err != nil {
			return errors.WithMessage(err, "failed to save purchase order")
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

// GetPurchaseOrders -
func (s *TradingApiService) GetPurchaseOrders(ctx context.Context, userID uuid.UUID) (ImplResponse, error) {
	if userID == uuid.Nil {
		return newUnauthorizedResponse(), nil
	}

	orders, err := s.purchaseOrders.FindByUser(ctx, userID)
	if err != nil {
		return Response(http.StatusInternalServerError, nil), err
	}

	return Response(http.StatusOK, orders), nil
}

// GetSellOrders -
func (s *TradingApiService) GetSellOrders(ctx context.Context, userID uuid.UUID) (ImplResponse, error) {
	if userID == uuid.Nil {
		return newUnauthorizedResponse(), nil
	}

	orders, err := s.sellOrders.FindByUser(ctx, userID)
	if err != nil {
		return Response(http.StatusInternalServerError, nil), err
	}

	return Response(http.StatusOK, orders), nil
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
	if form.UserRole != "broker" {
		return newForbiddenResponse(), nil
	}

	err := s.transactionManager.DoTransactionally(ctx, func(ctx context.Context) error {
		item := trading.NewItem(form.Name, form.InitialCount, form.InitialPrice, form.Commission)
		err := s.items.Add(ctx, item)
		if err != nil {
			return errors.WithMessagef(err, `failed to add item "%s"`, item.Name)
		}

		for i := 0; i < int(item.InitialCount); i++ {
			order := trading.NewInitialOrder(item.ID, item.InitialPrice, item.CalculateCommission(item.InitialPrice))
			err = s.sellOrders.Save(ctx, order)
			if err != nil {
				return errors.WithMessage(err, "failed to save initial order")
			}
		}

		return nil
	})
	if err != nil {
		return Response(http.StatusInternalServerError, nil), err
	}

	return Response(http.StatusNoContent, nil), nil
}
