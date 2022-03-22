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

	"github.com/gofrs/uuid"
	"github.com/muonsoft/validation"
	"github.com/muonsoft/validation/it"
)

type BillingOperation struct {
	AccountID      uuid.UUID `json:"-"`
	IdempotenceKey string    `json:"-"`
	Amount         float64   `json:"amount"`
}

func (operation BillingOperation) Validate(ctx context.Context, validator *validation.Validator) error {
	return validator.Validate(ctx,
		validation.NumberProperty("amount", operation.Amount, it.IsBetweenFloats(1.0, 10000.0)),
	)
}
