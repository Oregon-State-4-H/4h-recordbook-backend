package api

import (
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// for automatically mapping azcosmos response code to message
var HTTPResponseCodeMap = map[int]string {
	400: "Bad request",
	401: "Unauthorized",
	403: "Forbidden",
	404: "Item not found",
	409: "Conflict",
}
 
const (
	//400
	ErrBadRequest = "Bad request"
	ErrMissingFields = "One or more required fields is missing"
	ErrNoQuery = "No query provided"
	ErrBadDate = "The date(s) provided do not conform to the RFC3339 format."

	//401
	ErrNoToken = "No authentication token provided"
	ErrBadToken = "Bad token"
	
	//404
	ErrNotFound = "Item not found"

	//409
	ErrBookmarkConflict = "Bookmark with that link already exists"
)

type HTTPResponseCode struct {
	Code 	int
	Message string
}

func InterpretCosmosError(err error) HTTPResponseCode {
	
	var responseError *azcore.ResponseError
	errors.As(err, &responseError)

	code := responseError.StatusCode
	message, ok := HTTPResponseCodeMap[code]
	if !ok {
		message = "Unexpected error"
	}

	response := HTTPResponseCode{
		Code: code,
		Message: message,
	}

	return response

}