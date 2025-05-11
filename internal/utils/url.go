package utils

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func BaseUrl(c *gin.Context) string {

	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}

	return fmt.Sprintf("%s://%s", scheme, c.Request.Host)

}

func NextUrl(c *gin.Context, category string, page int, perPage int, sortByNewest bool) string {
	return fmt.Sprintf("%s/%s?page=%d&per_page=%d&sort_by_newest=%t", BaseUrl(c), category, page+1, perPage, sortByNewest)
}
