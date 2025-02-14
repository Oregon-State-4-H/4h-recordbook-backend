package db 

import (
	"context"
	"4h-recordbook-backend/internal/config"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type Db interface {
	GetUser(ctx context.Context, id string) (User, error)
	UpsertUser(ctx context.Context, user User) (interface{}, error)
	GetBookmarkByLink(ctx context.Context, userid string, link string) (Bookmark, error)
	GetBookmarks(ctx context.Context, userid string) ([]Bookmark, error)
	AddBookmark(ctx context.Context, bookmark Bookmark) (interface{}, error)
	RemoveBookmark(ctx context.Context, userid string, bookmarkid string) (interface{}, error)
	GetProjectByID(ctx context.Context, userid string, projectid string) (Project, error)
	GetCurrentProjects(ctx context.Context, userid string) ([]Project, error)
	GetProjectsByUser(ctx context.Context, userid string) ([]Project, error)
	UpsertProject(ctx context.Context, project Project) (interface{}, error)
	GetResume(ctx context.Context, userid string) (Resume, error)
	GetSection1ByID(ctx context.Context, userid string, sectionid string) (Section1, error)
	GetSection2ByID(ctx context.Context, userid string, sectionid string) (Section2, error)
	GetSection3ByID(ctx context.Context, userid string, sectionid string) (Section3, error)
	GetSection4ByID(ctx context.Context, userid string, sectionid string) (Section4, error)
	GetSection5ByID(ctx context.Context, userid string, sectionid string) (Section5, error)
	GetSection6ByID(ctx context.Context, userid string, sectionid string) (Section6, error)
	GetSection7ByID(ctx context.Context, userid string, sectionid string) (Section7, error)
	GetSection8ByID(ctx context.Context, userid string, sectionid string) (Section8, error)
	GetSection9ByID(ctx context.Context, userid string, sectionid string) (Section9, error)
	GetSection10ByID(ctx context.Context, userid string, sectionid string) (Section10, error)
	GetSection11ByID(ctx context.Context, userid string, sectionid string) (Section11, error)
	GetSection12ByID(ctx context.Context, userid string, sectionid string) (Section12, error)
	GetSection13ByID(ctx context.Context, userid string, sectionid string) (Section13, error)
	GetSection14ByID(ctx context.Context, userid string, sectionid string) (Section14, error)
	GetSection1(ctx context.Context, userid string) ([]Section1, error)
	GetSection2(ctx context.Context, userid string) ([]Section2, error)
	GetSection3(ctx context.Context, userid string) ([]Section3, error)
	GetSection4(ctx context.Context, userid string) ([]Section4, error)
	GetSection5(ctx context.Context, userid string) ([]Section5, error)
	GetSection6(ctx context.Context, userid string) ([]Section6, error)
	GetSection7(ctx context.Context, userid string) ([]Section7, error)
	GetSection8(ctx context.Context, userid string) ([]Section8, error)
	GetSection9(ctx context.Context, userid string) ([]Section9, error)
	GetSection10(ctx context.Context, userid string) ([]Section10, error)
	GetSection11(ctx context.Context, userid string) ([]Section11, error)
	GetSection12(ctx context.Context, userid string) ([]Section12, error)
	GetSection13(ctx context.Context, userid string) ([]Section13, error)
	GetSection14(ctx context.Context, userid string) ([]Section14, error)
	UpsertSection1(ctx context.Context, section Section1) (interface{}, error)
	UpsertSection2(ctx context.Context, section Section2) (interface{}, error)
	UpsertSection3(ctx context.Context, section Section3) (interface{}, error)
	UpsertSection4(ctx context.Context, section Section4) (interface{}, error)
	UpsertSection5(ctx context.Context, section Section5) (interface{}, error)
	UpsertSection6(ctx context.Context, section Section6) (interface{}, error)
	UpsertSection7(ctx context.Context, section Section7) (interface{}, error)
	UpsertSection8(ctx context.Context, section Section8) (interface{}, error)
	UpsertSection9(ctx context.Context, section Section9) (interface{}, error)
	UpsertSection10(ctx context.Context, section Section10) (interface{}, error)
	UpsertSection11(ctx context.Context, section Section11) (interface{}, error)
	UpsertSection12(ctx context.Context, section Section12) (interface{}, error)
	UpsertSection13(ctx context.Context, section Section13) (interface{}, error)
	UpsertSection14(ctx context.Context, section Section14) (interface{}, error)
	RemoveSection(ctx context.Context, userid string, sectionid string) (interface{}, error)
	GetAnimalsByProject(ctx context.Context, userid string, projectid string) ([]Animal, error)
	GetAnimalByID(ctx context.Context, userid string, animalid string) (Animal, error)
	UpsertAnimal(ctx context.Context, animal Animal) (interface{}, error)
	GetFeedsByProject(ctx context.Context, userid string, projectid string) ([]Feed, error)
	GetFeedByID(ctx context.Context, userid string, feedid string) (Feed, error)
	UpsertFeed(ctx context.Context, feed Feed) (interface{}, error)
	GetFeedPurchasesByProject(ctx context.Context, userid string, projectid string) ([]FeedPurchase, error)
	GetFeedPurchaseByID(ctx context.Context, userid string, feedPurchaseID string) (FeedPurchase, error)
	UpsertFeedPurchase(ctx context.Context, feedPurchase FeedPurchase) (interface{}, error)
	GetDailyFeedsByProjectAndAnimal(ctx context.Context, userid string, projectid string, animalid string) ([]DailyFeed, error)
	GetDailyFeedByID(ctx context.Context, userid string, dailyFeedID string) (DailyFeed, error)
	UpsertDailyFeed(ctx context.Context, dailyFeed DailyFeed) (interface{}, error)
	GetExpensesByProject(ctx context.Context, userid string, projectid string) ([]Expense, error)
	GetExpenseByID(ctx context.Context, userid string, expenseid string) (Expense, error)
	UpsertExpense(ctx context.Context, expense Expense) (interface{}, error)
	GetSuppliesByProject(ctx context.Context, userid string, projectid string) ([]Supply, error)
	GetSupplyByID(ctx context.Context, userid string, supplyid string) (Supply, error)
	UpsertSupply(ctx context.Context, supply Supply) (interface{}, error)
	RemoveSupply(ctx context.Context, userid string, supplyid string) (interface{}, error)
}

type GenericDatabaseInfo struct {
	Created string `json:"created"`
	Updated string `json:"updated"`
}

type env struct {
	logger	  *zap.SugaredLogger  	   `validate:"required"`
	validator *validator.Validate 	   `validate:"required"`
	client	  *azcosmos.DatabaseClient `validate:"required"`
}

func New(logger *zap.SugaredLogger, cfg *config.Config) (Db, error) {

	logger.Info("Creating new database client")

	validate := validator.New()

	cred, err := azcosmos.NewKeyCredential(cfg.Database.Current.Key)
	if err != nil {
		return nil, err
	}

	client, err := azcosmos.NewClientWithKey(cfg.Database.Current.Endpoint, cred, nil)
	if err != nil {
		return nil, err
	}

	dbClient, err := client.NewDatabase("recordbooks-db")
	if err != nil {
		return nil, err
	}

	e := &env{
		logger:	   logger,
		validator: validate,
		client:    dbClient,
	}

	return e, nil

}
