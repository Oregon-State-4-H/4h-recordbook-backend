package api

import (
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

var HTTPResponseCodeMap = map[int]string {
	401: "Unauthorized",
	403: "Forbidden",
	404: "Item not found",
	409: "Conflict",
	500: "Internal server error",
}

type HTTPResponseCode struct {
	Code 	int
	Message string
}

func InterpretCosmosError(err error) HTTPResponseCode {
	
	var responseError *azcore.ResponseError
	errors.As(err, &responseError)

	response := HTTPResponseCode{
		Code: responseError.StatusCode,
		Message: HTTPResponseCodeMap[responseError.StatusCode],
	}

	return response

}