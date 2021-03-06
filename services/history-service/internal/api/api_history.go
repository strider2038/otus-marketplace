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
	"net/http"
	"strings"

	"github.com/gofrs/uuid"
)

// A HistoryApiController binds http requests to an api service and writes the service results to the http response
type HistoryApiController struct {
	service HistoryApiServicer
}

// NewHistoryApiController creates a default api controller
func NewHistoryApiController(s HistoryApiServicer) Router {
	return &HistoryApiController{service: s}
}

// Routes returns all of the api route for the HistoryApiController
func (c *HistoryApiController) Routes() Routes {
	return Routes{
		{
			"GetDealsHistory",
			strings.ToUpper("Get"),
			"/api/v1/deals",
			c.GetDealsHistory,
		},
	}
}

// GetDealsHistory -
func (c *HistoryApiController) GetDealsHistory(w http.ResponseWriter, r *http.Request) {
	userID := uuid.FromStringOrNil(r.Header.Get("X-User-Id"))

	result, err := c.service.GetDealsHistory(r.Context(), userID)
	// If an error occured, encode the error with the status code
	if err != nil {
		EncodeJSONResponse(err.Error(), &result.Code, w)
		return
	}
	// If no error, encode the body and the result code
	EncodeJSONResponse(result.Body, &result.Code, w)

}
