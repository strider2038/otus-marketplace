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
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gofrs/uuid"
)

// A BillingApiController binds http requests to an api service and writes the service results to the http response
type BillingApiController struct {
	service BillingApiServicer
}

// NewBillingApiController creates a default api controller
func NewBillingApiController(s BillingApiServicer) Router {
	return &BillingApiController{service: s}
}

// Routes returns all of the api route for the BillingApiController
func (c *BillingApiController) Routes() Routes {
	return Routes{
		{
			"DepositMoney",
			strings.ToUpper("Post"),
			"/api/v1/account/deposit",
			c.DepositMoney,
		},
		{
			"GetBillingAccount",
			strings.ToUpper("Get"),
			"/api/v1/account",
			c.GetBillingAccount,
		},
		{
			"GetBillingOperations",
			strings.ToUpper("Get"),
			"/api/v1/operations",
			c.GetBillingOperations,
		},
		{
			"WithdrawMoney",
			strings.ToUpper("Post"),
			"/api/v1/account/withdraw",
			c.WithdrawMoney,
		},
	}
}

// GetBillingAccount -
func (c *BillingApiController) GetBillingAccount(w http.ResponseWriter, r *http.Request) {
	id := uuid.FromStringOrNil(r.Header.Get("X-User-Id"))

	result, err := c.service.GetBillingAccount(r.Context(), id)
	// If an error occured, encode the error with the status code
	if err != nil {
		EncodeJSONResponse(err.Error(), &result.Code, w)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// DepositMoney -
func (c *BillingApiController) DepositMoney(w http.ResponseWriter, r *http.Request) {
	form := &BillingOperation{}
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	form.AccountID = uuid.FromStringOrNil(r.Header.Get("X-User-Id"))

	result, err := c.service.DepositMoney(r.Context(), *form)
	// If an error occured, encode the error with the status code
	if err != nil {
		EncodeJSONResponse(err.Error(), &result.Code, w)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// WithdrawMoney -
func (c *BillingApiController) WithdrawMoney(w http.ResponseWriter, r *http.Request) {
	form := &BillingOperation{}
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	form.AccountID = uuid.FromStringOrNil(r.Header.Get("X-User-Id"))

	result, err := c.service.WithdrawMoney(r.Context(), *form)
	// If an error occured, encode the error with the status code
	if err != nil {
		EncodeJSONResponse(err.Error(), &result.Code, w)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}

// GetBillingOperations -
func (c *BillingApiController) GetBillingOperations(w http.ResponseWriter, r *http.Request) {
	id := uuid.FromStringOrNil(r.Header.Get("X-User-Id"))

	result, err := c.service.GetBillingOperations(r.Context(), id)
	// If an error occured, encode the error with the status code
	if err != nil {
		EncodeJSONResponse(err.Error(), &result.Code, w)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)
}
