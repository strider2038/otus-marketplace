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

	"github.com/gofrs/uuid"
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
)

type SellOrder struct {
	ItemID         uuid.UUID `json:"itemId,omitempty"`
	UserID         uuid.UUID `json:"-"`
	Price          float64   `json:"price"`
	IdempotenceKey string    `json:"-"`
}

func (order SellOrder) Validate(ctx context.Context, validator *validation.Validator) error {
	return validator.Validate(ctx,
		validation.NumberProperty("price", order.Price, it.IsGreaterThanFloat(1.0)),
	)
}
