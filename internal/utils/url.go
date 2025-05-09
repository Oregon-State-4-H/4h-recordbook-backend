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

func NextUrl(c *gin.Context, page int, perPage int) string {
	return fmt.Sprintf("%s/bookmarks?page=%d&per_page=%d", BaseUrl(c), page+1, perPage)
}
