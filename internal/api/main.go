package api

import (
	_ "4h-recordbook-backend/internal/api/docs"
	"4h-recordbook-backend/internal/config"
	"4h-recordbook-backend/internal/middleware"
	"4h-recordbook-backend/pkg/db"
	"4h-recordbook-backend/pkg/upc"
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

const (
	API_VERSION = "1.0"

	PAGE_DEFAULT     = 0
	PER_PAGE_DEFAULT = 100
	PER_PAGE_MIN     = 1
	PER_PAGE_MAX     = 200

	CONTEXT_KEY_PAGE           = "page"
	CONTEXT_KEY_PER_PAGE       = "per_page"
	CONTEXT_KEY_SORT_BY_NEWEST = "sort_by_newest"
)

type Api interface {
	RunLocal() error
	RunAzureFunctions(string) error
}

type env struct {
	logger    *zap.SugaredLogger  `validate:"required"`
	validator *validator.Validate `validate:"required"`
	config    *config.Config      `validate:"required"`
	db        db.Db               `validate:"required"`
	upc       upc.Upc             `validate:"required"`
	api       *gin.Engine         `validate:"required"`
}

func ternary(s1 string, s2 string) string {
	if s1 == "" {
		return s2
	}
	return s1
}

type UserInfo struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
}

type CustomClaims struct {
	jwt.StandardClaims
	UserInfo
}

func decodeJWT(c *gin.Context) (*CustomClaims, error) {

	var claims *CustomClaims

	auth := c.Request.Header.Get("Authorization")
	if auth == "" {
		return claims, errors.New(ErrNoToken)
	}

	splitToken := strings.Split(auth, "Bearer ")
	if len(splitToken) == 1 {
		return claims, errors.New(ErrBadToken)
	}
	auth = splitToken[1]

	token, err := jwt.ParseWithClaims(auth, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("AccessToken"), nil
	})

	if err != nil {
		return claims, errors.New(ErrBadToken)
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return claims, errors.New(ErrBadToken)

}

func generateJWT(userid string, firstName string) (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = CustomClaims{
		jwt.StandardClaims{
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

func (e *env) RunAzureFunctions(port string) error {
	return http.ListenAndServe(":"+port, e.api)
}

func PaginationMiddleware(defaultSortByNewest bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		pageStr := c.DefaultQuery(CONTEXT_KEY_PAGE, strconv.Itoa(PAGE_DEFAULT))
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 0 {
			page = PAGE_DEFAULT
		}

		perPageStr := c.DefaultQuery(CONTEXT_KEY_PER_PAGE, strconv.Itoa(PER_PAGE_DEFAULT))
		perPage, err := strconv.Atoi(perPageStr)
		if err != nil {
			perPage = PER_PAGE_DEFAULT
		}
		if perPage < PER_PAGE_MIN {
			perPage = PER_PAGE_MIN
		} else if perPage > PER_PAGE_MAX {
			perPage = PER_PAGE_MAX
		}

		sortByNewestStr := c.DefaultQuery(CONTEXT_KEY_SORT_BY_NEWEST, strconv.FormatBool(defaultSortByNewest))
		sortByNewest, err := strconv.ParseBool(sortByNewestStr)
		if err != nil {
			sortByNewest = defaultSortByNewest
		}

		c.Set(CONTEXT_KEY_PAGE, page)
		c.Set(CONTEXT_KEY_PER_PAGE, perPage)
		c.Set(CONTEXT_KEY_SORT_BY_NEWEST, sortByNewest)

		c.Next()
	}
}

func CtxMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func New(logger *zap.SugaredLogger, cfg *config.Config, dbInstance db.Db, upcInstance upc.Upc) (Api, error) {

	logger.Info("Setting up API")

	validator := validator.New()

	e := &env{
		validator: validator,
		logger:    logger,
		config:    cfg,
		db:        dbInstance,
		upc:       upcInstance,
	}

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
	}))

	router.Use(CtxMiddleware(5 * time.Second))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "hello",
		})
	})

	router.POST("/signup", e.signup)
	router.POST("/signin", e.signin)

	/*Require Authentication for all later calls.*/
	router.Use(middleware.EnsureValidToken(), middleware.GetToken());

	router.POST("/register", e.register)

	/*Require a valid user*/
	router.Use(middleware.GetUser(e.db));

	router.GET("/auth", func(c *gin.Context){
		c.String(http.StatusOK, "Authorized successfully, welcome, " + c.GetString("user_name") + "!")
	})

	router.GET("/user", e.getUserProfile)
	router.PUT("/user", e.updateUserProfile)
	

	router.GET("/bookmarks", PaginationMiddleware(false), e.getUserBookmarks)
	router.GET("/bookmarks/:link", e.getBookmarkByLink)
	router.POST("/bookmarks", e.addUserBookmark)
	router.DELETE("/bookmarks/:bookmarkID", e.deleteUserBookmark)

	router.GET("/projects", PaginationMiddleware(true), e.getCurrentProjects)
	router.GET("/project", PaginationMiddleware(true), e.getProjects)
	router.GET("/project/:projectID", e.getProject)
	router.POST("/project", e.addProject)
	router.PUT("/project/:projectID", e.updateProject)
	router.DELETE("/project/:projectID", e.deleteProject)

	router.GET("/resume", e.getResume)

	router.GET("/section1", PaginationMiddleware(false), e.getSection1s)
	router.GET("/section1/:sectionID", e.getSection1)
	router.POST("/section1", e.addSection1)
	router.PUT("/section1/:sectionID", e.updateSection1)

	router.GET("/section2", PaginationMiddleware(false), e.getSection2s)
	router.GET("/section2/:sectionID", e.getSection2)
	router.POST("/section2", e.addSection2)
	router.PUT("/section2/:sectionID", e.updateSection2)

	router.GET("/section3", PaginationMiddleware(false), e.getSection3s)
	router.GET("/section3/:sectionID", e.getSection3)
	router.POST("/section3", e.addSection3)
	router.PUT("/section3/:sectionID", e.updateSection3)

	router.GET("/section4", PaginationMiddleware(false), e.getSection4s)
	router.GET("/section4/:sectionID", e.getSection4)
	router.POST("/section4", e.addSection4)
	router.PUT("/section4/:sectionID", e.updateSection4)

	router.GET("/section5", PaginationMiddleware(false), e.getSection5s)
	router.GET("/section5/:sectionID", e.getSection5)
	router.POST("/section5", e.addSection5)
	router.PUT("/section5/:sectionID", e.updateSection5)

	router.GET("/section6", PaginationMiddleware(false), e.getSection6s)
	router.GET("/section6/:sectionID", e.getSection6)
	router.POST("/section6", e.addSection6)
	router.PUT("/section6/:sectionID", e.updateSection6)

	router.GET("/section7", PaginationMiddleware(false), e.getSection7s)
	router.GET("/section7/:sectionID", e.getSection7)
	router.POST("/section7", e.addSection7)
	router.PUT("/section7/:sectionID", e.updateSection7)

	router.GET("/section8", PaginationMiddleware(false), e.getSection8s)
	router.GET("/section8/:sectionID", e.getSection8)
	router.POST("/section8", e.addSection8)
	router.PUT("/section8/:sectionID", e.updateSection8)

	router.GET("/section9", PaginationMiddleware(false), e.getSection9s)
	router.GET("/section9/:sectionID", e.getSection9)
	router.POST("/section9", e.addSection9)
	router.PUT("/section9/:sectionID", e.updateSection9)

	router.GET("/section10", PaginationMiddleware(false), e.getSection10s)
	router.GET("/section10/:sectionID", e.getSection10)
	router.POST("/section10", e.addSection10)
	router.PUT("/section10/:sectionID", e.updateSection10)

	router.GET("/section11", PaginationMiddleware(false), e.getSection11s)
	router.GET("/section11/:sectionID", e.getSection11)
	router.POST("/section11", e.addSection11)
	router.PUT("/section11/:sectionID", e.updateSection11)

	router.GET("/section12", PaginationMiddleware(false), e.getSection12s)
	router.GET("/section12/:sectionID", e.getSection12)
	router.POST("/section12", e.addSection12)
	router.PUT("/section12/:sectionID", e.updateSection12)

	router.GET("/section13", PaginationMiddleware(false), e.getSection13s)
	router.GET("/section13/:sectionID", e.getSection13)
	router.POST("/section13", e.addSection13)
	router.PUT("/section13/:sectionID", e.updateSection13)

	router.GET("/section14", PaginationMiddleware(false), e.getSection14s)
	router.GET("/section14/:sectionID", e.getSection14)
	router.POST("/section14", e.addSection14)
	router.PUT("/section14/:sectionID", e.updateSection14)

	router.DELETE("/section/:sectionID", e.deleteSection)

	router.GET("/event", PaginationMiddleware(false), e.getEvents)
	router.POST("/event", e.addEvent)
	router.PUT("/event/:eventID", e.updateEvent)
	router.DELETE("/event/:eventID", e.deleteEvent)
	router.GET("/event/:eventID", e.getEventWithSections)
	router.POST("/event/:eventID", e.addEventSection)
	router.DELETE("event/:eventID/:sectionID", e.deleteEventSection)

	router.GET("/project/:projectID/animal", PaginationMiddleware(false), e.getAnimals)
	router.GET("/animal/:animalID", e.getAnimal)
	router.POST("/animal", e.addAnimal)
	router.PUT("/animal/:animalID", e.updateAnimal)
	router.PUT("/rate-of-gain/:animalID", e.updateRateOfGain)
	router.DELETE("/animal/:animalID", e.deleteAnimal)

	router.GET("/project/:projectID/feed", PaginationMiddleware(false), e.getFeeds)
	router.GET("/feed/:feedID", e.getFeed)
	router.POST("/feed", e.addFeed)
	router.PUT("/feed/:feedID", e.updateFeed)
	router.DELETE("/feed/:feedID", e.deleteFeed)

	router.GET("/project/:projectID/feed-purchase", PaginationMiddleware(false), e.getFeedPurchases)
	router.GET("/feed-purchase/:feedPurchaseID", e.getFeedPurchase)
	router.POST("/feed-purchase", e.addFeedPurchase)
	router.PUT("/feed-purchase/:feedPurchaseID", e.updateFeedPurchase)
	router.DELETE("/feed-purchase/:feedPurchaseID", e.deleteFeedPurchase)

	router.GET("/project/:projectID/animal/:animalID/daily-feed", PaginationMiddleware(false), e.getDailyFeeds)
	router.GET("/daily-feed/:dailyFeedID", e.getDailyFeed)
	router.POST("/daily-feed", e.addDailyFeed)
	router.PUT("/daily-feed/:dailyFeedID", e.updateDailyFeed)
	router.DELETE("/daily-feed/:dailyFeedID", e.deleteDailyFeed)

	router.GET("/project/:projectID/expense", PaginationMiddleware(false), e.getExpenses)
	router.GET("/expense/:expenseID", e.getExpense)
	router.POST("/expense", e.addExpense)
	router.PUT("/expense/:expenseID", e.updateExpense)
	router.DELETE("/expense/:expenseID", e.deleteExpense)

	router.GET("/project/:projectID/supply", PaginationMiddleware(false), e.getSupplies)
	router.GET("/supply/:supplyID", e.getSupply)
	router.POST("/supply", e.addSupply)
	router.PUT("/supply/:supplyID", e.updateSupply)
	router.DELETE("/supply/:supplyID", e.deleteSupply)

	router.GET("/upc/:code", e.getUpcProduct)

	e.api = router

	return e, nil

}

// @title	4H Record Books API
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {

}
