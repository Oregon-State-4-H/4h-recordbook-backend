package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"4h-recordbook-backend/internal/utils"
	"4h-recordbook-backend/pkg/db"
	"github.com/beevik/guid"
)

type GetSuppliesOutput struct {
	Supplies []db.Supply `json:"supplies"`
}

type GetSupplyOutput struct {
	Supply db.Supply `json:"supply"`
}

type UpsertSupplyInput struct {
	Description string `json:"description" validate:"required"`
	StartValue *float64 `json:"start_value" validate:"required"`
	EndValue *float64 `json:"end_value" validate:"required"`
	ProjectID string `json:"project_id" validate:"required"`
}

type UpsertSupplyOutput GetSupplyOutput

// GetSupplies godoc
// @Summary Get supplies by project
// @Description Gets all of a user's supplies given a project ID
// @Tags Supply
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param projectID query string true "Project ID"
// @Success 200 {object} api.GetSuppliesOutput
// @Failure 400
// @Failure 401
// @Router /supply [get]
func (e *env) getSupplies(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	projectID := c.DefaultQuery("projectID", "")
	if projectID == "" {
		c.JSON(400, gin.H{
			"message": ErrNoQuery,
		})
		return
	}

	var output GetSuppliesOutput

	output.Supplies, err = e.db.GetSuppliesByProject(context.TODO(), claims.ID, projectID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

// GetSupply godoc
// @Summary Get a supply
// @Description Get a user's supply by ID
// @Tags Supply
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param supplyID path string true "Supply ID"
// @Success 200 {object} api.GetSupplyOutput
// @Failure 401
// @Failure 404
// @Router /supply/{supplyID} [get]
func (e *env) getSupply(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	supplyID := c.Param("supplyID")

	var output GetSupplyOutput

	output.Supply, err = e.db.GetSupplyByID(context.TODO(), claims.ID, supplyID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

// AddSupply godoc
// @Summary Add a supply
// @Description Adds a supply to a user's personal records
// @Tags Supply
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param UpsertSupplyInput body api.UpsertSupplyInput true "Supply information"
// @Success 201 {object} api.UpsertSupplyOutput
// @Failure 400
// @Failure 401
// @Router /supply [post]
func (e *env) addSupply(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	var input UpsertSupplyInput
	err = c.BindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{
			"message": ErrBadRequest,
		})
		return
	}

	err = e.validator.Struct(input)
	if err != nil {
		c.JSON(400, gin.H{
			"message": ErrMissingFields,
		})
		return
	}

	g := guid.New()
	timestamp := utils.TimeNow()

	supply := db.Supply{
		ID: g.String(),
		Description: input.Description,
		StartValue: *input.StartValue,
		EndValue: *input.EndValue,
		ProjectID: input.ProjectID,
		UserID: claims.ID,
		GenericDatabaseInfo: db.GenericDatabaseInfo {
			Created: timestamp.String(),
			Updated: timestamp.String(),
		},
	}

	var output UpsertSupplyOutput

	output.Supply, err = e.db.UpsertSupply(context.TODO(), supply)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(201, output)

}

// UpdateSupply godoc
// @Summary Update a supply
// @Description Updates a user's supply information
// @Tags Supply
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param supplyID path string true "Supply ID"
// @Param UpsertSupplyInput body api.UpsertSupplyInput true "Supply information"
// @Success 200 {object} api.UpsertSupplyOutput
// @Failure 400
// @Failure 401
// @Failure 404
// @Router /supply/{supplyID} [put]
func (e *env) updateSupply(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	var input UpsertSupplyInput
	err = c.BindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{
			"message": ErrBadRequest,
		})
		return
	}
	
	err = e.validator.Struct(input)
	if err != nil {
		c.JSON(400, gin.H{
			"message": ErrMissingFields,
		})
		return
	}

	supplyID := c.Param("supplyID")

	supply, err := e.db.GetSupplyByID(context.TODO(), claims.ID, supplyID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	timestamp := utils.TimeNow()

	updatedSupply := db.Supply{
		ID: supply.ID,
		Description: input.Description,
		StartValue: *input.StartValue,
		EndValue: *input.EndValue,
		ProjectID: supply.ProjectID,
		UserID: claims.ID,
		GenericDatabaseInfo: db.GenericDatabaseInfo {
			Created: supply.Created,
			Updated: timestamp.String(),
		},
	}

	var output UpsertSupplyOutput

	output.Supply, err = e.db.UpsertSupply(context.TODO(), updatedSupply)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

// DeleteSupply godoc
// @Summary Removes a supply
// @Description Deletes a user's supply given the supply ID
// @Tags Supply
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param supplyID path string true "Supply ID"
// @Success 204
// @Failure 401
// @Failure 404 
// @Router /supply/{supplyID} [delete]
func (e *env) deleteSupply(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	supplyID := c.Param("supplyID")

	response, err := e.db.RemoveSupply(context.TODO(), claims.ID, supplyID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(204, response)

}