/*
 * Billing service
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

// BillingApiRouter defines the required methods for binding the api requests to a responses for the BillingApi
// The BillingApiRouter implementation should parse necessary information from the http request,
// pass the data to a BillingApiServicer to perform the required actions, then write the service results to the http response.
type BillingApiRouter interface {
	DepositMoney(http.ResponseWriter, *http.Request)
	GetBillingAccount(http.ResponseWriter, *http.Request)
	GetBillingOperations(http.ResponseWriter, *http.Request)
	WithdrawMoney(http.ResponseWriter, *http.Request)
}

// BillingApiServicer defines the api actions for the BillingApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type BillingApiServicer interface {
	DepositMoney(context.Context, BillingOperation) (ImplResponse, error)
	GetBillingAccount(context.Context, uuid.UUID) (ImplResponse, string, error)
	GetBillingOperations(context.Context, uuid.UUID) (ImplResponse, error)
	WithdrawMoney(context.Context, BillingOperation) (ImplResponse, error)
}
