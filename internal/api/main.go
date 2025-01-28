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

func ternary(s1 string, s2 string) (string){
	if s1 == "" {
		return s2
	}
	return s1
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

	router.GET("/projects", e.getCurrentProjects)
  	router.GET("/project", e.getProjects)
  	router.GET("/project/:projectId", e.getProject)
  	router.POST("/project", e.addProject)
	router.PUT("/project/:projectId", e.updateProject)

	router.GET("/resume", e.getResume)

	router.GET("/section1", e.getSection1)
	router.POST("/section1", e.addSection1)
	router.PUT("/section1/:sectionId", e.updateSection1)

	router.GET("/section2", e.getSection2)
	router.POST("/section2", e.addSection2)
	router.PUT("/section2/:sectionId", e.updateSection2)

	router.GET("/section3", e.getSection3)
	router.POST("/section3", e.addSection3)
	router.PUT("/section3/:sectionId", e.updateSection3)

	router.GET("/section4", e.getSection4)
	router.POST("/section4", e.addSection4)
	router.PUT("/section4/:sectionId", e.updateSection4)

	router.GET("/section5", e.getSection5)
	router.POST("/section5", e.addSection5)
	router.PUT("/section5/:sectionId", e.updateSection5)

	router.GET("/section6", e.getSection6)
	router.POST("/section6", e.addSection6)
	router.PUT("/section6/:sectionId", e.updateSection6)

	router.GET("/section7", e.getSection7)
	router.POST("/section7", e.addSection7)
	router.PUT("/section7/:sectionId", e.updateSection7)

	router.GET("/section8", e.getSection8)
	router.POST("/section8", e.addSection8)
	router.PUT("/section8/:sectionId", e.updateSection8)

	router.GET("/section9", e.getSection9)
	router.POST("/section9", e.addSection9)
	router.PUT("/section9/:sectionId", e.updateSection9)

	router.GET("/section10", e.getSection10)
	router.POST("/section10", e.addSection10)
	router.PUT("/section10/:sectionId", e.updateSection10)

	router.GET("/section11", e.getSection11)
	router.POST("/section11", e.addSection11)
	router.PUT("/section11/:sectionId", e.updateSection11)

	router.GET("/section12", e.getSection12)
	router.POST("/section12", e.addSection12)
	router.PUT("/section12/:sectionId", e.updateSection12)

	router.GET("/section13", e.getSection13)
	router.POST("/section13", e.addSection13)
	router.PUT("/section13/:sectionId", e.updateSection13)

	router.GET("/section14", e.getSection14)
	router.POST("/section14", e.addSection14)
	router.PUT("/section14/:sectionId", e.updateSection14)

	router.DELETE("/section/:sectionId", e.deleteSection)

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
