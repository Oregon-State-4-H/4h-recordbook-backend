package api

import (
	"errors"
	"time"
	"strings"
	"net/http"
	_ "4h-recordbook-backend/internal/api/docs"
	"4h-recordbook-backend/internal/config"
	"4h-recordbook-backend/pkg/db"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
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
	logger *zap.SugaredLogger `validate:"required"`
	validator *validator.Validate `validate:"required"`
	config *config.Config `validate:"required"`
	db db.Db `validate:"required"`
	api *gin.Engine `validate:"required"`
}

func ternary(s1 string, s2 string) (string){
	if s1 == "" {
		return s2
	}
	return s1
}

type UserInfo struct {
	ID string `json:"id"`
	FirstName string `json:"first_name"`
}

type CustomClaims struct {
	jwt.StandardClaims
	UserInfo
}

func decodeJWT(c *gin.Context) (*CustomClaims, error){

	var claims *CustomClaims

	auth := c.Request.Header.Get("Authorization")
	if auth == "" {
		return claims, errors.New("Unauthorized 1") 
	}

	splitToken := strings.Split(auth, "Bearer ")
	auth = splitToken[1]

	token, err := jwt.ParseWithClaims(auth, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("AccessToken"), nil
	})

	if err != nil {
		return claims, errors.New("Unauthorized 2") 
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return claims, errors.New("Unauthorized 3")

}

func generateJWT(userid string, firstName string) (string, error){

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = CustomClaims {
		jwt.StandardClaims {
			ExpiresAt: time.Now().Add(time.Minute * 60).Unix(),
		},
		UserInfo{
			userid,
			firstName,
		},
	}

	return token.SignedString([]byte("AccessToken"))

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

	router.GET("/section1", e.getSection1s)
	router.GET("/section1/:sectionId", e.getSection1)
	router.POST("/section1", e.addSection1)
	router.PUT("/section1/:sectionId", e.updateSection1)

	router.GET("/section2", e.getSection2s)
	router.GET("/section2/:sectionId", e.getSection2)
	router.POST("/section2", e.addSection2)
	router.PUT("/section2/:sectionId", e.updateSection2)

	router.GET("/section3", e.getSection3s)
	router.GET("/section3/:sectionId", e.getSection3)
	router.POST("/section3", e.addSection3)
	router.PUT("/section3/:sectionId", e.updateSection3)

	router.GET("/section4", e.getSection4s)
	router.GET("/section4/:sectionId", e.getSection4)
	router.POST("/section4", e.addSection4)
	router.PUT("/section4/:sectionId", e.updateSection4)

	router.GET("/section5", e.getSection5s)
	router.GET("/section5/:sectionId", e.getSection5)
	router.POST("/section5", e.addSection5)
	router.PUT("/section5/:sectionId", e.updateSection5)

	router.GET("/section6", e.getSection6s)
	router.GET("/section6/:sectionId", e.getSection6)
	router.POST("/section6", e.addSection6)
	router.PUT("/section6/:sectionId", e.updateSection6)

	router.GET("/section7", e.getSection7s)
	router.GET("/section7/:sectionId", e.getSection7)
	router.POST("/section7", e.addSection7)
	router.PUT("/section7/:sectionId", e.updateSection7)

	router.GET("/section8", e.getSection8s)
	router.GET("/section8/:sectionId", e.getSection8)
	router.POST("/section8", e.addSection8)
	router.PUT("/section8/:sectionId", e.updateSection8)

	router.GET("/section9", e.getSection9s)
	router.GET("/section9/:sectionId", e.getSection9)
	router.POST("/section9", e.addSection9)
	router.PUT("/section9/:sectionId", e.updateSection9)

	router.GET("/section10", e.getSection10s)
	router.GET("/section10/:sectionId", e.getSection10)
	router.POST("/section10", e.addSection10)
	router.PUT("/section10/:sectionId", e.updateSection10)

	router.GET("/section11", e.getSection11s)
	router.GET("/section11/:sectionId", e.getSection11)
	router.POST("/section11", e.addSection11)
	router.PUT("/section11/:sectionId", e.updateSection11)

	router.GET("/section12", e.getSection12s)
	router.GET("/section12/:sectionId", e.getSection12)
	router.POST("/section12", e.addSection12)
	router.PUT("/section12/:sectionId", e.updateSection12)

	router.GET("/section13", e.getSection13s)
	router.GET("/section13/:sectionId", e.getSection13)
	router.POST("/section13", e.addSection13)
	router.PUT("/section13/:sectionId", e.updateSection13)

	router.GET("/section14", e.getSection14s)
	router.GET("/section14/:sectionId", e.getSection14)
	router.POST("/section14", e.addSection14)
	router.PUT("/section14/:sectionId", e.updateSection14)

	router.DELETE("/section/:sectionId", e.deleteSection)

	router.GET("/animal", e.getAnimals)
	router.GET("/animal/:animalId", e.getAnimal)
	router.POST("/animal", e.addAnimal)
	router.PUT("/animal/:animalId", e.updateAnimal)
	router.PUT("/rate-of-gain/:animalId", e.updateRateOfGain)

	router.GET("/feed", e.getFeeds)
	router.GET("/feed/:feedId", e.getFeed)
	router.POST("/feed", e.addFeed)
	router.PUT("/feed/:feedId", e.updateFeed)
	
	router.GET("/feed-purchase", e.getFeedPurchases)
	router.GET("/feed-purchase/:feedPurchaseId", e.getFeedPurchase)
	router.POST("/feed-purchase", e.addFeedPurchase)
	router.PUT("/feed-purchase/:feedPurchaseId", e.updateFeedPurchase)

	router.GET("/daily-feed", e.getDailyFeeds)
	router.GET("/daily-feed/:dailyFeedId", e.getDailyFeed)
	router.POST("/daily-feed", e.addDailyFeed)
	router.PUT("/daily-feed/:dailyFeedId", e.updateDailyFeed)

	router.GET("/expense", e.getExpenses)
	router.GET("/expense/:expenseId", e.getExpense)
	router.POST("/expense", e.addExpense)

	router.GET("/supply/", e.getSupplies)
	router.GET("/supply/:supplyId", e.getSupply)
	router.POST("/supply", e.addSupply)
	router.PUT("/supply/:supplyId", e.updateSupply)
	router.DELETE("/supply/:supplyId", e.deleteSupply)

	e.api = router

	return e, nil

}

// @title	4H Record Books API
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {

}