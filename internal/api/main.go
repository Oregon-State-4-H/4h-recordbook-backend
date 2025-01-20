package api

import (
	"net/http"
	"4h-recordbook-backend/internal/config"
	"4h-recordbook-backend/pkg/db"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

const API_VERSION = "1.0"

type Api interface {
	RunLocal() error
}

type env struct {
	logger	  *zap.SugaredLogger  `validate:"required"`
	validator *validator.Validate `validate:"required"`
	config    *config.Config	  `validate:"required"`
	db 		  db.Db				  `validate:"required"`
	api		  *gin.Engine		  `validate:"required"`
}

func (e *env) RunLocal() error {
	return http.ListenAndServe("localhost:8080", e.api)
}

func New(logger *zap.SugaredLogger, cfg *config.Config, dbInstance db.Db) (Api, error) {

	logger.Info("Setting up API")

	validator := validator.New()

	e := &env {
		validator: validator,
		logger:    logger,
		config:    cfg,
		db:        dbInstance,
	}

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowHeaders: 	 []string{"Authorization", "Content-Type"},
		AllowAllOrigins: true,
		AllowMethods:	 []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
	}))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "hello",
		})
	})

	router.GET("/user", e.getUserProfile)
	router.PUT("/user", e.updateUserProfile)
	router.POST("/signin", e.signin)
	router.POST("/signout", e.signout)
	router.POST("/signup", e.signup)

	router.GET("/bookmarks", e.getUserBookmarks)
	router.POST("/bookmarks", e.addUserBookmark)
	router.DELETE("/bookmarks/:bookmarkId", e.removeUserBookmark)

	router.GET("/projects", getCurrentProjects)
  	router.GET("/project", getProjects)
  	router.GET("/project/{projectId}", getProject)
  	router.POST("/project", addProject)
	router.PUT("/project/{projectId}", updateProject)

	router.GET("/resume", getResumeDocs)

	router.GET("/section-1/{docId}", getSection1)
	router.GET("/section-1", getSection1Docs)
	router.POST("/section-1", addSection1)
	router.PUT("/section-1", updateSection1)
	router.DELETE("/section-1/{docId}", deleteSection1)

	router.GET("/animal/docs/{projectId}", getAnimalDocs)
	router.GET("/animal/{animalId}", getAnimal)
	router.POST("/animal", addAnimal)
	router.PUT("/animal", updateAnimal)
	router.PUT("/rate-of-gain", updateRateOfGain)

	router.GET("/feed/docs/{projectId}", getFeedDocs)
	router.GET("/feed/{feedId}", getFeed)
	router.POST("/feed", addFeed)
	//no endpoint for addFeedNoForm, all that's different is it adds a name. could probably just be added to POST /feed
	router.PUT("/feed", updateFeed)
	
	router.GET("/feed-purchase/docs/{projecId}", getFeedPurchaseDocs)
	router.GET("/feed-purchase/{feedPurchaseId}", getFeedPurchase)
	router.POST("/feed-purchase", addFeedPurchase)
	router.PUT("/feed-purchase", updateFeedPurchase)

	router.GET("/daily-feed/docs/{projectId}/{animalId}", getDailyFeedDocs)
	router.GET("/daily-feed/{dailyFeedId}", getDailyFeed)
	router.POST("/daily-feed", addDailyFeed)
	router.PUT("/daily-feed", updateDailyFeed)

	router.GET("/expenses/docs/{projectId}", getExpenseDocs)
	router.GET("/expenses/{expenseId}", getExpense)
	router.POST("/expenses/", addExpense)

	router.GET("/supply/docs/{projectId}", getSupplyDocs)
	router.GET("/supply/{supplyId}", getSupply)
	router.POST("/supply", addSupply)
	router.PUT("/supply", updateSupply)
	router.DELETE("/supply/{supplyId}", deleteSupply)

	e.api = router

	return e, nil

}
