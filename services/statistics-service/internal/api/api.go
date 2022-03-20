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
)

// StatisticsApiRouter defines the required methods for binding the api requests to a responses for the StatisticsApi
// The StatisticsApiRouter implementation should parse necessary information from the http request,
// pass the data to a StatisticsApiServicer to perform the required actions, then write the service results to the http response.
type StatisticsApiRouter interface {
	GetDailyStats(http.ResponseWriter, *http.Request)
	GetTop10(http.ResponseWriter, *http.Request)
	GetTotalDailyStats(http.ResponseWriter, *http.Request)
}

// StatisticsApiServicer defines the api actions for the StatisticsApi service
// This interface intended to stay up to date with the openapi yaml used to generate it,
// while the service implementation can ignored with the .openapi-generator-ignore file
// and updated with the logic required for the API.
type StatisticsApiServicer interface {
	GetDailyStats(context.Context) (ImplResponse, error)
	GetTop10(context.Context) (ImplResponse, error)
	GetTotalDailyStats(context.Context) (ImplResponse, error)
}