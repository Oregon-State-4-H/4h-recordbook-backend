package api

import (
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// for automatically mapping azcosmos response code to message
var HTTPResponseCodeMap = map[int]string{
	400: "bad request",
	401: "unauthorized",
	403: "forbidden",
	404: "item not found",
	409: "conflict",
}

const (
	//400
	ErrBadRequest           = "bad request"
	ErrMissingFields        = "one or more required fields is missing"
	ErrBadDate              = "the date(s) provided do not conform to the RFC3339 format."
	ErrInvalidSectionNumber = "section number must be in the range [1-14] inclusive"
	ErrQueryMustBeInt       = "query param must be an integer value"

	//401
	ErrNoToken  = "no authentication token provided"
	ErrBadToken = "bad token"

	//404
	ErrNotFound = "item not found"

	//409
	ErrBookmarkConflict     = "bookmark with that link already exists"
	ErrEventSectionConflict = "event already has this section"
)

type HTTPResponseCode struct {
	Code    int
	Message string
}

func InterpretCosmosError(err error) HTTPResponseCode {

	var responseError *azcore.ResponseError
	errors.As(err, &responseError)

	//catch errors from db package that aren't cosmos errors
	if responseError == nil {
		return HTTPResponseCode{
			Code:    500,
			Message: err.Error(),
		}
	}

	code := responseError.StatusCode
	message, ok := HTTPResponseCodeMap[code]
	if !ok {
		message = "unexpected error"
	}

	return HTTPResponseCode{
		Code:    code,
		Message: message,
	}

}
