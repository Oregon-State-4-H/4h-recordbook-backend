package api

import (
	"4h-recordbook-backend/internal/utils"
	"4h-recordbook-backend/pkg/db"
	"strconv"

	"github.com/beevik/guid"
	"github.com/gin-gonic/gin"
)

type GetFeedsOutput struct {
	Feeds []db.Feed `json:"feeds"`
	Next  string    `json:"next"`
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
// @Param projectID path string true "Project ID"
// @Param page query int false "Page number, default 0"
// @Param per_page query int false "Max number of items to return. Can be [1-200], default 100"
// @Param sort_by_newest query bool false "Sort results by most recently added, default false"
// @Success 200 {object} api.GetFeedsOutput
// @Failure 400
// @Failure 401
// @Router /project/{projectID}/feed [get]
func (e *env) getFeeds(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		return
	}

	projectID := c.Param("projectID")

	var output GetFeedsOutput

	paginationOptions := db.PaginationOptions{
		Page:         c.GetInt(CONTEXT_KEY_PAGE),
		PerPage:      c.GetInt(CONTEXT_KEY_PER_PAGE),
		SortByNewest: c.GetBool(CONTEXT_KEY_SORT_BY_NEWEST),
	}

	output.Feeds, err = e.db.GetFeedsByProject(c.Request.Context(), claims.ID, projectID, paginationOptions)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	if len(output.Feeds) == paginationOptions.PerPage {

		queryParamsMap := make(map[string]string)
		queryParamsMap[CONTEXT_KEY_PAGE] = strconv.Itoa(paginationOptions.Page + 1)
		queryParamsMap[CONTEXT_KEY_PER_PAGE] = strconv.Itoa(paginationOptions.PerPage)
		queryParamsMap[CONTEXT_KEY_SORT_BY_NEWEST] = strconv.FormatBool(paginationOptions.SortByNewest)

		nextUrlInput := utils.NextUrlInput{
			Context:     c,
			QueryParams: queryParamsMap,
		}

		output.Next = utils.BuildNextUrl(nextUrlInput)
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

	output.Feed, err = e.db.GetFeedByID(c.Request.Context(), claims.ID, feedID)
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

	output.Feed, err = e.db.UpsertFeed(c.Request.Context(), feed)
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

	feed, err := e.db.GetFeedByID(c.Request.Context(), claims.ID, feedID)
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

	output.Feed, err = e.db.UpsertFeed(c.Request.Context(), updatedFeed)
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

	response, err := e.db.RemoveFeed(c.Request.Context(), claims.ID, feedID)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(204, response)

}
