package api

import "net/http"

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

func newUnprocessableEntityResponse(message string) ImplResponse {
	return Response(
		http.StatusUnprocessableEntity,
		Error{Code: http.StatusUnprocessableEntity, Message: message},
	)
}
