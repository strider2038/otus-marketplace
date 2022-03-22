package api

import (
	"net/http"

	"github.com/pkg/errors"
)

func newUnauthorizedResponse() ImplResponse {
	return Response(
		http.StatusUnauthorized,
		Error{Code: http.StatusUnauthorized, Message: "Authorization is required"},
	)
}

func newForbiddenResponse() ImplResponse {
	return Response(
		http.StatusForbidden,
		Error{Code: http.StatusForbidden, Message: "Operation denied"},
	)
}

func newNotFoundResponse() ImplResponse {
	return Response(
		http.StatusNotFound,
		Error{Code: http.StatusNotFound, Message: "Not found"},
	)
}

func newConflictResponse() ImplResponse {
	return Response(http.StatusConflict, Error{
		Code:    http.StatusConflict,
		Message: "Resource is locked",
	})
}

func newPreconditionFailedResponse() ImplResponse {
	return Response(http.StatusPreconditionFailed, Error{
		Code:    http.StatusPreconditionFailed,
		Message: "State of resource has changed",
	})
}

func newUnprocessableEntityResponse(message string) ImplResponse {
	return Response(
		http.StatusUnprocessableEntity,
		Error{Code: http.StatusUnprocessableEntity, Message: message},
	)
}

func newPreconditionRequiredResponse() ImplResponse {
	return Response(http.StatusPreconditionRequired, Error{
		Code:    http.StatusPreconditionRequired,
		Message: "Missing If-Match header",
	})
}

func newErrorResponsef(err error, format string, args ...interface{}) (ImplResponse, error) {
	return Response(http.StatusInternalServerError, nil), errors.WithMessagef(err, format, args...)
}
