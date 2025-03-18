package api

import (
	"4h-recordbook-backend/pkg/upc"

	"github.com/gin-gonic/gin"
)

type GetUpcProductOutput struct {
	Product upc.Product `json:"product"`
}

// GetUpcProduct godoc
// @Summary Get a UPC product
// @Description
// @Tags UPC
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param code path string true "UPC code"
// @Success 200 {object} upc.Product
// @Failure 400
// @Failure 401
// @Router /upc/{code} [get]
func (e *env) getUpcProduct(c *gin.Context) {

	_, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	code := c.Param("code")

	var output GetUpcProductOutput

	output.Product, err = e.upc.GetProductByCode(code)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

	c.JSON(200, output)

}
