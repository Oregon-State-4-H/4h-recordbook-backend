package api

import (
	"4h-recordbook-backend/internal/utils"
	"4h-recordbook-backend/pkg/db"
	"context"

	"github.com/beevik/guid"
	"github.com/gin-gonic/gin"
)

type GetDailyFeedsOutput struct {
	DailyFeeds []db.DailyFeed `json:"daily_feeds"`
}

type GetDailyFeedOutput struct {
	DailyFeed db.DailyFeed `json:"daily_feed"`
}

type UpsertDailyFeedInput struct {
	FeedDate       string   `json:"feed_date" validate:"required"`
	FeedAmount     *float64 `json:"feed_amount" validate:"required"`
	AnimalID       string   `json:"animal_id" validate:"required"`
	FeedID         string   `json:"feed_id" validate:"required"`
	FeedPurchaseID string   `json:"feed_purchase_id" validate:"required"`
	ProjectID      string   `json:"project_id" validate:"required"`
}

type UpsertDailyFeedOutput GetDailyFeedOutput

// GetDailyFeeds godoc
// @Summary Get daily feeds by project and animal
// @Description Gets all of a user's daily feeds for a given project and animal
// @Tags Daily Feed
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param projectID query string true "Project ID"
// @Param animalID query string true "Animal ID"
// @Success 200 {object} api.GetDailyFeedsOutput
// @Failure 400
// @Failure 401
// @Router /daily-feed [get]
func (e *env) getDailyFeeds(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
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

	animalID := c.DefaultQuery("animalID", "")
	if animalID == "" {
		c.JSON(400, gin.H{
			"message": ErrNoQuery,
		})
		return
	}

	var output GetDailyFeedsOutput

	output.DailyFeeds, err = e.db.GetDailyFeedsByProjectAndAnimal(context.TODO(), claims.ID, projectID, animalID)
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

// GetDailyFeed godoc
// @Summary Get a daily feed
// @Description Get a user's daily feed by ID
// @Tags Daily Feed
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param dailyFeedID path string true "Daily Feed ID"
// @Success 200 {object} api.GetDailyFeedOutput
// @Failure 401
// @Failure 404
// @Router /daily-feed/{dailyFeedID} [get]
func (e *env) getDailyFeed(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	dailyFeedID := c.Param("dailyFeedID")

	var output GetDailyFeedOutput

	output.DailyFeed, err = e.db.GetDailyFeedByID(context.TODO(), claims.ID, dailyFeedID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

// AddDailyFeed godoc
// @Summary Add a daily feed
// @Description Adds a daily feed to a user's personal records
// @Tags Daily Feed
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param UpsertDailyFeedInput body api.UpsertDailyFeedInput true "Daily Feed information"
// @Success 201 {object} api.UpsertDailyFeedOutput
// @Failure 400
// @Failure 401
// @Router /daily-feed [post]
func (e *env) addDailyFeed(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	var input UpsertDailyFeedInput
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

	feedDate, err := utils.StringToTimestamp(input.FeedDate)
	if err != nil {
		c.JSON(400, gin.H{
			"message": ErrBadDate,
		})
		return
	}

	g := guid.New()
	timestamp := utils.TimeNow()

	dailyFeed := db.DailyFeed{
		ID:             g.String(),
		FeedDate:       feedDate.String(),
		FeedAmount:     *input.FeedAmount,
		AnimalID:       input.AnimalID,
		FeedID:         input.FeedID,
		FeedPurchaseID: input.FeedPurchaseID,
		ProjectID:      input.ProjectID,
		UserID:         claims.ID,
		GenericDatabaseInfo: db.GenericDatabaseInfo{
			Created: timestamp.String(),
			Updated: timestamp.String(),
		},
	}

	var output UpsertDailyFeedOutput

	output.DailyFeed, err = e.db.UpsertDailyFeed(context.TODO(), dailyFeed)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(201, output)

}

// UpdateDailyFeed godoc
// @Summary Update a daily feed
// @Description Updates a user's daily feed information
// @Tags Daily Feed
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param dailyFeedID path string true "Daily Feed ID"
// @Param UpsertDailyFeedInput body api.UpsertDailyFeedInput true "DailyFeed information"
// @Success 200 {object} api.UpsertDailyFeedOutput
// @Failure 400
// @Failure 401
// @Failure 404
// @Router /daily-feed/{dailyFeedID} [put]
func (e *env) updateDailyFeed(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	var input UpsertDailyFeedInput
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

	feedDate, err := utils.StringToTimestamp(input.FeedDate)
	if err != nil {
		c.JSON(400, gin.H{
			"message": ErrBadDate,
		})
		return
	}

	dailyFeedID := c.Param("dailyFeedID")

	dailyFeed, err := e.db.GetDailyFeedByID(context.TODO(), claims.ID, dailyFeedID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	timestamp := utils.TimeNow()

	updatedDailyFeed := db.DailyFeed{
		ID:             dailyFeed.ID,
		FeedDate:       feedDate.String(),
		FeedAmount:     *input.FeedAmount,
		AnimalID:       dailyFeed.AnimalID,
		FeedID:         dailyFeed.FeedID,
		FeedPurchaseID: dailyFeed.FeedPurchaseID,
		ProjectID:      dailyFeed.ProjectID,
		UserID:         claims.ID,
		GenericDatabaseInfo: db.GenericDatabaseInfo{
			Created: dailyFeed.Created,
			Updated: timestamp.String(),
		},
	}

	var output UpsertDailyFeedOutput

	output.DailyFeed, err = e.db.UpsertDailyFeed(context.TODO(), updatedDailyFeed)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

// DeleteDailyFeed godoc
// @Summary Removes a daily feed
// @Description Deletes a user's daily feed given the daily feed ID
// @Tags Daily Feed
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param dailyFeedID path string true "Daily Feed ID"
// @Success 204
// @Failure 401
// @Failure 404
// @Router /daily-feed/{dailyFeedID} [delete]
func (e *env) deleteDailyFeed(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	dailyFeedID := c.Param("dailyFeedID")

	response, err := e.db.RemoveDailyFeed(context.TODO(), claims.ID, dailyFeedID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(204, response)

}
