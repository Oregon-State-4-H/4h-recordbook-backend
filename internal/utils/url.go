package utils

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

type NextUrlInput struct {
	Context     *gin.Context
	QueryParams map[string]string
}

func BaseUrl(c *gin.Context) string {

	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}

	return fmt.Sprintf("%s://%s", scheme, c.Request.Host)

}

// this function is responsible for building the url given the context + query params. this function is NOT responsible for determining the correct values
// e.g. the calling function is responsible for incrementing the page number
func BuildNextUrl(input NextUrlInput) string {

	var queryStringBuilder strings.Builder
	queryStringBuilder.WriteString("?")

	for key, value := range input.QueryParams {
		param := fmt.Sprintf("%s=%s&", key, value)
		queryStringBuilder.WriteString(param)
	}

	queryString := queryStringBuilder.String()
	queryString = strings.TrimRight(queryString, "&")

	return fmt.Sprintf("%s%s%s", BaseUrl(input.Context), input.Context.Request.URL.Path, queryString)

}
