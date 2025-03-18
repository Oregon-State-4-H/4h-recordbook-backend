package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"4h-recordbook-backend/internal/utils"
	"4h-recordbook-backend/pkg/db"
	"github.com/beevik/guid"
)

type UpsertFeedPurchaseInput struct {
	DatePurchased string `json:"date_purchased" validate:"required"`
	AmountPurchased *float64 `json:"amount_purchased" validate:"required"`
	TotalCost *float64 `json:"total_cost" validate:"required"`
	FeedID string `json:"feed_id" validate:"required"`
	ProjectID string `json:"project_id" validate:"required"`
}

type GetFeedPurchasesOutput struct {
	FeedPurchases []db.FeedPurchase `json:"feed_purchases"`
}

type GetFeedPurchaseOutput struct {
	FeedPurchase db.FeedPurchase `json:"feed_purchase"`
}

type UpsertFeedPurchaseOutput GetFeedPurchaseOutput

// GetFeedPurchases godoc
// @Summary Get feed purchases by project
// @Description Gets all of a user's feed purchases given a project ID
// @Tags Feed Purchase
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param projectID query string true "Project ID"
// @Success 200 {object} api.GetFeedPurchasesOutput
// @Failure 400
// @Failure 401
// @Router /feed-purchase [get]
func (e *env) getFeedPurchases(c *gin.Context) {

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

	var output GetFeedPurchasesOutput

	output.FeedPurchases, err = e.db.GetFeedPurchasesByProject(context.TODO(), claims.ID, projectID)
	if err != nil {
		e.logger.Info(err)
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

// GetFeedPurchase godoc
// @Summary Get a feed purchase
// @Description Get a user's feed purchase by ID
// @Tags Feed Purchase
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param feedPurchaseID path string true "Feed Purchase ID"
// @Success 200 {object} api.GetFeedPurchaseOutput
// @Failure 401
// @Failure 404
// @Router /feed-purchase/{feedPurchaseID} [get]
func (e *env) getFeedPurchase(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	feedPurchaseID := c.Param("feedPurchaseID")

	var output GetFeedPurchaseOutput

	output.FeedPurchase, err = e.db.GetFeedPurchaseByID(context.TODO(), claims.ID, feedPurchaseID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

// AddFeedPurchase godoc
// @Summary Add a feed purchase
// @Description Adds a feed purchase to a user's personal records
// @Tags Feed Purchase
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param UpsertFeedPurchaseInput body api.UpsertFeedPurchaseInput true "Feed Purchase information"
// @Success 201 {object} api.UpsertFeedPurchaseOutput
// @Failure 400
// @Failure 401
// @Router /feed-purchase [post]
func (e *env) addFeedPurchase(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	var input UpsertFeedPurchaseInput
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

	datePurchased, err := utils.StringToTimestamp(input.DatePurchased)
	if err != nil {
		c.JSON(400, gin.H{
			"message": ErrBadDate,
		})
		return
	}

	g := guid.New()
	timestamp := utils.TimeNow()

	feedPurchase := db.FeedPurchase{
		ID: g.String(),
		DatePurchased: datePurchased.String(),
		AmountPurchased: *input.AmountPurchased,
		TotalCost: *input.TotalCost,
		FeedID: input.FeedID,
		ProjectID: input.ProjectID,
		UserID: claims.ID,
		GenericDatabaseInfo: db.GenericDatabaseInfo {
			Created: timestamp.String(),
			Updated: timestamp.String(),
		},
	}

	var output UpsertFeedPurchaseOutput

	output.FeedPurchase, err = e.db.UpsertFeedPurchase(context.TODO(), feedPurchase)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(201, output)

}

// UpdateFeedPurchase godoc
// @Summary Update a feed purchase
// @Description Updates a user's feed purchase information
// @Tags Feed Purchase
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param feedPurchaseID path string true "Feed Purchase ID"
// @Param UpsertFeedPurchaseInput body api.UpsertFeedPurchaseInput true "Feed purchase information"
// @Success 200 {object} api.UpsertFeedPurchaseOutput
// @Failure 400
// @Failure 401
// @Failure 404
// @Router /feed-purchase/{feedPurchaseID} [put]
func (e *env) updateFeedPurchase(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	var input UpsertFeedPurchaseInput
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

	datePurchased, err := utils.StringToTimestamp(input.DatePurchased)
	if err != nil {
		c.JSON(400, gin.H{
			"message": ErrBadDate,
		})
		return
	}

	feedPurchaseID := c.Param("feedPurchaseID")

	feedPurchase, err := e.db.GetFeedPurchaseByID(context.TODO(), claims.ID, feedPurchaseID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	timestamp := utils.TimeNow()

	updatedFeedPurchase := db.FeedPurchase{
		ID: feedPurchase.ID,
		DatePurchased: datePurchased.String(),
		AmountPurchased: *input.AmountPurchased,
		TotalCost: *input.TotalCost,
		FeedID: feedPurchase.FeedID,
		ProjectID: feedPurchase.ProjectID,
		UserID: claims.ID,
		GenericDatabaseInfo: db.GenericDatabaseInfo {
			Created: feedPurchase.Created,
			Updated: timestamp.String(),
		},
	}

	var output UpsertFeedPurchaseOutput

	output.FeedPurchase, err = e.db.UpsertFeedPurchase(context.TODO(), updatedFeedPurchase)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

// DeleteFeedPurchase godoc
// @Summary Removes a feed purchase
// @Description Deletes a user's feed purchase given the feed purchase ID
// @Tags Feed Purchase
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param feedPurchaseID path string true "Feed Purchase ID"
// @Success 204
// @Failure 401
// @Failure 404 
// @Router /feed-purchase/{feedPurchaseID} [delete]
func (e *env) deleteFeedPurchase(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	feedPurchaseID := c.Param("feedPurchaseID")

	response, err := e.db.RemoveFeedPurchase(context.TODO(), claims.ID, feedPurchaseID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(204, response)

}