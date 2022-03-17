/*
 * Public API
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package api

type TradingItem struct {
	UserRole     string  `json:"-"`
	Name         string  `json:"name,omitempty"`
	InitialCount int32   `json:"initialCount,omitempty"`
	InitialPrice float64 `json:"initialPrice,omitempty"`
	Commission   float64 `json:"commission,omitempty"`
}
