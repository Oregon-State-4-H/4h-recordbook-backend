package router

import (
  	_ "fmt"
  	_ "4h-recordbook/backend/docs"
	"4h-recordbook/backend/internal/handlers"
  	"github.com/gin-gonic/gin"
  	ginSwagger "github.com/swaggo/gin-swagger"
  	swaggerFiles "github.com/swaggo/files"
)

func New() *gin.Engine {

	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

  	router.GET("/", func(c *gin.Context){
		c.JSON(200, gin.H{
			"message": "Welcome to my website",
		})
  	})

	router.GET("/user", handlers.GetUserProfile)
	router.PUT("/user", handlers.UpdateUserProfile)
	
	router.GET("/bookmarks", handlers.GetUserBookmarks)
	router.POST("/bookmarks", handlers.AddUserBookmark)
	router.DELETE("/bookmarks/{bookmarkId}", handlers.RemoveUserBookmark)

	router.GET("/projects", handlers.GetCurrentProjects)
  	router.GET("/project", handlers.GetProjects)
  	router.GET("/project/{projectId}", handlers.GetProject)
  	router.POST("/project", handlers.AddProject)
	router.PUT("/project/{projectId}", handlers.UpdateProject)

	router.GET("/resume", handlers.GetResumeDocs)

	router.GET("/section-1/{docId}", handlers.GetSection1)
	router.GET("/section-1", handlers.GetSection1Docs)
	router.POST("/section-1", handlers.AddSection1)
	router.PUT("/section-1", handlers.UpdateSection1)
	router.DELETE("/section-1/{docId}", handlers.DeleteSection1)

	router.GET("/animal/docs/{projectId}", handlers.GetAnimalDocs)
	router.GET("/animal/{animalId}", handlers.GetAnimal)
	router.POST("/animal", handlers.AddAnimal)
	router.PUT("/animal", handlers.UpdateAnimal)
	router.PUT("/rate-of-gain", handlers.UpdateRateOfGain)

	router.GET("/feed/docs/{projectId}", handlers.GetFeedDocs)
	router.GET("/feed/{feedId}", handlers.GetFeed)
	router.POST("/feed", handlers.AddFeed)
	//no endpoint for handlers.AddFeedNoForm, all that's different is it adds a name. could probably just be added to POST /feed
	router.PUT("/feed", handlers.UpdateFeed)
	
	router.GET("/feed-purchase/docs/{projecId}", handlers.GetFeedPurchaseDocs)
	router.GET("/feed-purchase/{feedPurchaseId}", handlers.GetFeedPurchase)
	router.POST("/feed-purchase", handlers.AddFeedPurchase)
	router.PUT("/feed-purchase", handlers.UpdateFeedPurchase)

	router.GET("/daily-feed/docs/{projectId}/{animalId}", handlers.GetDailyFeedDocs)
	router.GET("/daily-feed/{dailyFeedId}", handlers.GetDailyFeed)
	router.POST("/daily-feed", handlers.AddDailyFeed)
	router.PUT("/daily-feed", handlers.UpdateDailyFeed)

	router.GET("/expenses/docs/{projectId}", handlers.GetExpenseDocs)
	router.GET("/expenses/{expenseId}", handlers.GetExpense)
	router.POST("/expenses/", handlers.AddExpense)

	router.GET("/supply/docs/{projectId}", handlers.GetSupplyDocs)
	router.GET("/supply/{supplyId}", handlers.GetSupply)
	router.POST("/supply", handlers.AddSupply)
	router.PUT("/supply", handlers.UpdateSupply)
	router.DELETE("/supply/{supplyId}", handlers.DeleteSupply)
	
	return router;

}