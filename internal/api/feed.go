package api

import (
	"4h-recordbook-backend/internal/utils"
	"4h-recordbook-backend/pkg/db"
	"context"

	"github.com/beevik/guid"
	"github.com/gin-gonic/gin"
)

type GetFeedsOutput struct {
	Feeds []db.Feed `json:"feeds"`
}

type GetFeedOutput struct {
	Feed db.Feed `json:"feed"`
}

type UpsertFeedInput struct {
	Name      string `json:"name" validate:"required"`
	ProjectID string `json:"project_id" validate:"required"`
}

type UpsertFeedOutput GetFeedOutput

// GetFeeds godoc
// @Summary Get feeds by project
// @Description Gets all of a user's feeds given a project ID
// @Tags Feed
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param projectID query string true "Project ID"
// @Success 200 {object} api.GetFeedsOutput
// @Failure 400
// @Failure 401
// @Router /feed [get]
func (e *env) getFeeds(c *gin.Context) {

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

	var output GetFeedsOutput

	output.Feeds, err = e.db.GetFeedsByProject(context.TODO(), claims.ID, projectID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

// GetFeed godoc
// @Summary Get a feed
// @Description Get a user's feed by ID
// @Tags Feed
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param feedID path string true "Feed ID"
// @Success 200 {object} api.GetFeedOutput
// @Failure 401
// @Failure 404
// @Router /feed/{feedID} [get]
func (e *env) getFeed(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	feedID := c.Param("feedID")

	var output GetFeedOutput

	output.Feed, err = e.db.GetFeedByID(context.TODO(), claims.ID, feedID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

// AddFeed godoc
// @Summary Add a feed
// @Description Adds a feed to a user's personal records
// @Tags Feed
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param UpsertFeedInput body api.UpsertFeedInput true "Feed information"
// @Success 201 {object} api.UpsertFeedOutput
// @Failure 400
// @Failure 401
// @Router /feed [post]
func (e *env) addFeed(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	var input UpsertFeedInput
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

	feed := db.Feed{
		ID:        g.String(),
		Name:      input.Name,
		ProjectID: input.ProjectID,
		UserID:    claims.ID,
		GenericDatabaseInfo: db.GenericDatabaseInfo{
			Created: timestamp.String(),
			Updated: timestamp.String(),
		},
	}

	var output UpsertFeedOutput

	output.Feed, err = e.db.UpsertFeed(context.TODO(), feed)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(201, output)

}

// UpdateFeed godoc
// @Summary Update a feed
// @Description Updates a user's feed information
// @Tags Feed
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param feedID path string true "Feed ID"
// @Param UpsertFeedInput body api.UpsertFeedInput true "Feed information"
// @Success 200 {object} api.UpsertFeedOutput
// @Failure 400
// @Failure 401
// @Failure 404
// @Router /feed/{feedID} [put]
func (e *env) updateFeed(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	var input UpsertFeedInput
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

	feedID := c.Param("feedID")

	feed, err := e.db.GetFeedByID(context.TODO(), claims.ID, feedID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	timestamp := utils.TimeNow()

	updatedFeed := db.Feed{
		ID:        feed.ID,
		Name:      input.Name,
		ProjectID: feed.ProjectID,
		UserID:    claims.ID,
		GenericDatabaseInfo: db.GenericDatabaseInfo{
			Created: feed.Created,
			Updated: timestamp.String(),
		},
	}

	var output UpsertFeedOutput

	output.Feed, err = e.db.UpsertFeed(context.TODO(), updatedFeed)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(200, output)

}

// DeleteFeed godoc
// @Summary Removes a feed
// @Description Deletes a user's feed given the feed ID
// @Tags Feed
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param feedID path string true "Feed ID"
// @Success 204
// @Failure 401
// @Failure 404
// @Router /feed/{feedID} [delete]
func (e *env) deleteFeed(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	feedID := c.Param("feedID")

	response, err := e.db.RemoveFeed(context.TODO(), claims.ID, feedID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(204, response)

}
