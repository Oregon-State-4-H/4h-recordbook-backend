package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"4h-recordbook-backend/internal/utils"
	"4h-recordbook-backend/pkg/db"
	"github.com/beevik/guid"
)

/*******************************
* ANIMALS
********************************/

func getAnimals(c *gin.Context) {
	
}

func getAnimal(c *gin.Context) {
	
}

func addAnimal(c *gin.Context) {
	
}

func updateAnimal(c *gin.Context) {
	
}

func updateRateOfGain(c *gin.Context) {
	
}

/*******************************
* FEED
********************************/

type UpsertFeedInput struct {
	Name string `json:"name"`
	ProjectID string `json:"projectid`
}

type GetFeedsOutput struct {
	Feeds []db.Feed `json:"feeds"`
}

type GetFeedOutput struct {
	Feed db.Feed `json:"feed"`
}

// GetFeeds godoc
// @Summary Get feeds by project
// @Description Gets all of a user's feeds given a project ID
// @Tags Feed
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param projectId query string true "Project ID"
// @Success 200 {object} api.GetFeedsOutput
// @Failure 400
// @Failure 401
// @Router /feed [get]
func (e *env) getFeeds(c *gin.Context) {
	
	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	projectID := c.DefaultQuery("projectId", "")
	if projectID == "" {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
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
// @Param feedId path string true "Feed ID"
// @Success 200 {object} api.GetFeedOutput
// @Failure 401
// @Failure 404
// @Router /feed/{feedId} [get]
func (e *env) getFeed(c *gin.Context) {
	
	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	id := c.Param("feedId")

	var output GetFeedOutput

	output.Feed, err = e.db.GetFeedByID(context.TODO(), claims.ID, id)
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
// @Success 204
// @Failure 400
// @Failure 401
// @Router /feed [post]
func (e *env) addFeed(c *gin.Context) {
	
	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	var input UpsertFeedInput
	err = c.BindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

	err = e.validator.Struct(input)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

	g := guid.New()
	timestamp := utils.TimeNow()

	feed := db.Feed{
		ID: g.String(),
		Name: input.Name,
		ProjectID: input.ProjectID,
		UserID: claims.ID,
		GenericDatabaseInfo: db.GenericDatabaseInfo {
			Created: timestamp.ToString(),
			Updated: timestamp.ToString(),
		},
	}

	response, err := e.db.UpsertFeed(context.TODO(), feed)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(204, response)

}

// UpdateFeed godoc
// @Summary Update a feed
// @Description Updates a user's feed information
// @Tags Feed
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param feedId path string true "Feed ID"
// @Param UpsertFeedInput body api.UpsertFeedInput true "Feed information"
// @Success 204 
// @Failure 400
// @Failure 401
// @Failure 404
// @Router /feed/{feedId} [put]
func (e *env) updateFeed(c *gin.Context) {
	
	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	var input UpsertFeedInput
	err = c.BindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}
	
	err = e.validator.Struct(input)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

	id := c.Param("feedId")

	feed, err := e.db.GetFeedByID(context.TODO(), claims.ID, id)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	timestamp := utils.TimeNow()

	updatedFeed := db.Feed{
		ID: feed.ID,
		Name: input.Name,
		ProjectID: feed.ProjectID,
		UserID: claims.ID,
		GenericDatabaseInfo: db.GenericDatabaseInfo {
			Created: feed.Created,
			Updated: timestamp.ToString(),
		},
	}

	response, err := e.db.UpsertFeed(context.TODO(), updatedFeed)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(204, response)

}

/*******************************
* FEED PURCHASE
********************************/

type UpsertFeedPurchaseInput struct {
	DatePurchased string `json:"date_purchased"`
	AmountPurchased float64 `json:"amount_purchased"`
	TotalCost float64 `json:"total_cost"`
	FeedID string `json:"feedid"`
	ProjectID string `json:"projectid"`
}

type GetFeedPurchasesOutput struct {
	FeedPurchases []db.FeedPurchase `json:"feed_purchases"`
}

type GetFeedPurchaseOutput struct {
	FeedPurchase db.FeedPurchase `json:"feed_purchase"`
}

// GetFeedPurchases godoc
// @Summary Get feed purchases by project
// @Description Gets all of a user's feed purchass given a project ID
// @Tags Feed Purchase
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param projectId query string true "Project ID"
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

	projectID := c.DefaultQuery("projectId", "")
	if projectID == "" {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

	e.logger.Info("1")

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
// @Param feedPurchaseId path string true "Feed Purchase ID"
// @Success 200 {object} api.GetFeedPurchaseOutput
// @Failure 401
// @Failure 404
// @Router /feed-purchase/{feedPurchaseId} [get]
func (e *env) getFeedPurchase(c *gin.Context) {

	claims, err := decodeJWT(c)
	if err != nil {
		c.JSON(401, gin.H{
			"message": HTTPResponseCodeMap[401],
		})
		return
	}

	id := c.Param("feedPurchaseId")

	var output GetFeedPurchaseOutput

	output.FeedPurchase, err = e.db.GetFeedPurchaseByID(context.TODO(), claims.ID, id)
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
// @Success 204
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
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

	err = e.validator.Struct(input)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

	datePurchased, err := utils.StringToTimestamp(input.DatePurchased)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

	g := guid.New()
	timestamp := utils.TimeNow()

	feedPurchase := db.FeedPurchase{
		ID: g.String(),
		DatePurchased: datePurchased.ToString(),
		AmountPurchased: input.AmountPurchased,
		TotalCost: input.TotalCost,
		FeedID: input.FeedID,
		ProjectID: input.ProjectID,
		UserID: claims.ID,
		GenericDatabaseInfo: db.GenericDatabaseInfo {
			Created: timestamp.ToString(),
			Updated: timestamp.ToString(),
		},
	}

	response, err := e.db.UpsertFeedPurchase(context.TODO(), feedPurchase)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(204, response)

}

// UpdateFeedPurchase godoc
// @Summary Update a feed purchase
// @Description Updates a user's feed purchase information
// @Tags Feed Purchase
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param feedPurchaseId path string true "Feed Purchase ID"
// @Param UpsertFeedPurchaseInput body api.UpsertFeedPurchaseInput true "Feed purchase information"
// @Success 204 
// @Failure 400
// @Failure 401
// @Failure 404
// @Router /feed-purchase/{feedPurchaseId} [put]
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
			"message": HTTPResponseCodeMap[400],
		})
		return
	}
	
	err = e.validator.Struct(input)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

	datePurchased, err := utils.StringToTimestamp(input.DatePurchased)
	if err != nil {
		c.JSON(400, gin.H{
			"message": HTTPResponseCodeMap[400],
		})
		return
	}

	id := c.Param("feedPurchaseId")

	feedPurchase, err := e.db.GetFeedPurchaseByID(context.TODO(), claims.ID, id)
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
		DatePurchased: datePurchased.ToString(),
		AmountPurchased: input.AmountPurchased,
		TotalCost: input.TotalCost,
		FeedID: feedPurchase.FeedID,
		ProjectID: feedPurchase.ProjectID,
		UserID: claims.ID,
		GenericDatabaseInfo: db.GenericDatabaseInfo {
			Created: feedPurchase.Created,
			Updated: timestamp.ToString(),
		},
	}

	response, err := e.db.UpsertFeedPurchase(context.TODO(), updatedFeedPurchase)
	if err != nil {
		response := InterpretCosmosError(err)
		c.JSON(response.Code, gin.H{
			"message": response.Message,
		})
		return
	}

	c.JSON(204, response)

}

/*******************************
* DAILY FEED
********************************/

func getDailyFeeds(c *gin.Context) {

} 

func getDailyFeed(c *gin.Context) {

}

func addDailyFeed(c *gin.Context) {

}

func updateDailyFeed(c *gin.Context) {

}

/*******************************
* EXPENSES
********************************/

func getExpenses(c *gin.Context) {

}

func getExpense(c *gin.Context) {

}

func addExpense(c *gin.Context) {

}

/*******************************
* SUPPLY INVENTORY
********************************/
func getSupplies(c *gin.Context) {

}

func getSupply(c *gin.Context) {

}

func addSupply(c *gin.Context) {

}

func updateSupply(c *gin.Context) {

}

func deleteSupply(c *gin.Context) {

}