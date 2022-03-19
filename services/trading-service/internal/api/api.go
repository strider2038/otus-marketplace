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

	"github.com/gofrs/uuid"
)

// TradingApiRouter defines the required methods for binding the api requests to a responses for the TradingApi
// The TradingApiRouter implementation should parse necessary information from the http request,
// pass the data to a TradingApiServicer to perform the required actions, then write the service results to the http response.
type TradingApiRouter interface {
	CancelPurchaseOrder(http.ResponseWriter, *http.Request)
	CancelSellOrder(http.ResponseWriter, *http.Request)
	CreatePurchaseOrder(http.ResponseWriter, *http.Request)
	CreateSellOrder(http.ResponseWriter, *http.Request)
	CreateTradingItem(http.ResponseWriter, *http.Request)
	GetPurchaseOrders(http.ResponseWriter, *http.Request)
	GetSellOrders(http.ResponseWriter, *http.Request)
	GetTradingItems(http.ResponseWriter, *http.Request)
	GetUserItems(http.ResponseWriter, *http.Request)
}

// TradingApiServicer defines the api actions for the TradingApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type TradingApiServicer interface {
	CancelPurchaseOrder(context.Context, uuid.UUID, uuid.UUID) (ImplResponse, error)
	CancelSellOrder(context.Context, uuid.UUID, uuid.UUID) (ImplResponse, error)
	CreatePurchaseOrder(context.Context, PurchaseOrder) (ImplResponse, error)
	CreateSellOrder(context.Context, SellOrder) (ImplResponse, error)
	CreateTradingItem(context.Context, TradingItem) (ImplResponse, error)
	GetPurchaseOrders(context.Context, uuid.UUID) (ImplResponse, error)
	GetSellOrders(context.Context, uuid.UUID) (ImplResponse, error)
	GetTradingItems(context.Context) (ImplResponse, error)
	GetUserItems(context.Context, uuid.UUID) (ImplResponse, error)
}
