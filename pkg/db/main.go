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
	GetBookmarkByLink(ctx context.Context, userID string, link string) (Bookmark, error)
	GetBookmarks(ctx context.Context, userID string) ([]Bookmark, error)
	AddBookmark(ctx context.Context, bookmark Bookmark) (Bookmark, error)
	RemoveBookmark(ctx context.Context, userID string, bookmarkID string) (interface{}, error)
	GetProjectByID(ctx context.Context, userID string, projectID string) (Project, error)
	GetCurrentProjects(ctx context.Context, userID string) ([]Project, error)
	GetProjectsByUser(ctx context.Context, userID string) ([]Project, error)
	UpsertProject(ctx context.Context, project Project) (interface{}, error)
	RemoveProject(ctx context.Context, userID string, projectID string) (interface{}, error)
	GetResume(ctx context.Context, userID string) (Resume, error)
	GetSection1ByID(ctx context.Context, userID string, sectionID string) (Section1, error)
	GetSection2ByID(ctx context.Context, userID string, sectionID string) (Section2, error)
	GetSection3ByID(ctx context.Context, userID string, sectionID string) (Section3, error)
	GetSection4ByID(ctx context.Context, userID string, sectionID string) (Section4, error)
	GetSection5ByID(ctx context.Context, userID string, sectionID string) (Section5, error)
	GetSection6ByID(ctx context.Context, userID string, sectionID string) (Section6, error)
	GetSection7ByID(ctx context.Context, userID string, sectionID string) (Section7, error)
	GetSection8ByID(ctx context.Context, userID string, sectionID string) (Section8, error)
	GetSection9ByID(ctx context.Context, userID string, sectionID string) (Section9, error)
	GetSection10ByID(ctx context.Context, userID string, sectionID string) (Section10, error)
	GetSection11ByID(ctx context.Context, userID string, sectionID string) (Section11, error)
	GetSection12ByID(ctx context.Context, userID string, sectionID string) (Section12, error)
	GetSection13ByID(ctx context.Context, userID string, sectionID string) (Section13, error)
	GetSection14ByID(ctx context.Context, userID string, sectionID string) (Section14, error)
	GetSection1sByUser(ctx context.Context, userID string) ([]Section1, error)
	GetSection2sByUser(ctx context.Context, userID string) ([]Section2, error)
	GetSection3sByUser(ctx context.Context, userID string) ([]Section3, error)
	GetSection4sByUser(ctx context.Context, userID string) ([]Section4, error)
	GetSection5sByUser(ctx context.Context, userID string) ([]Section5, error)
	GetSection6sByUser(ctx context.Context, userID string) ([]Section6, error)
	GetSection7sByUser(ctx context.Context, userID string) ([]Section7, error)
	GetSection8sByUser(ctx context.Context, userID string) ([]Section8, error)
	GetSection9sByUser(ctx context.Context, userID string) ([]Section9, error)
	GetSection10sByUser(ctx context.Context, userID string) ([]Section10, error)
	GetSection11sByUser(ctx context.Context, userID string) ([]Section11, error)
	GetSection12sByUser(ctx context.Context, userID string) ([]Section12, error)
	GetSection13sByUser(ctx context.Context, userID string) ([]Section13, error)
	GetSection14sByUser(ctx context.Context, userID string) ([]Section14, error)
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
	RemoveSection(ctx context.Context, userID string, sectionID string) (interface{}, error)
	GetAnimalsByProject(ctx context.Context, userID string, projectID string) ([]Animal, error)
	GetAnimalByID(ctx context.Context, userID string, animalID string) (Animal, error)
	UpsertAnimal(ctx context.Context, animal Animal) (Animal, error)
	RemoveAnimal(ctx context.Context, userID string, animalID string) (interface{}, error)
	GetFeedsByProject(ctx context.Context, userID string, projectID string) ([]Feed, error)
	GetFeedByID(ctx context.Context, userID string, feedID string) (Feed, error)
	UpsertFeed(ctx context.Context, feed Feed) (Feed, error)
	RemoveFeed(ctx context.Context, userID string, feedID string) (interface{}, error)
	GetFeedPurchasesByProject(ctx context.Context, userID string, projectID string) ([]FeedPurchase, error)
	GetFeedPurchaseByID(ctx context.Context, userID string, feedPurchaseID string) (FeedPurchase, error)
	UpsertFeedPurchase(ctx context.Context, feedPurchase FeedPurchase) (FeedPurchase, error)
	RemoveFeedPurchase(ctx context.Context, userID string, feedpurchaseID string) (interface{}, error)
	GetDailyFeedsByProjectAndAnimal(ctx context.Context, userID string, projectID string, animalID string) ([]DailyFeed, error)
	GetDailyFeedByID(ctx context.Context, userID string, dailyFeedID string) (DailyFeed, error)
	UpsertDailyFeed(ctx context.Context, dailyFeed DailyFeed) (DailyFeed, error)
	RemoveDailyFeed(ctx context.Context, userID string, dailyfeedID string) (interface{}, error)
	GetExpensesByProject(ctx context.Context, userID string, projectID string) ([]Expense, error)
	GetExpenseByID(ctx context.Context, userID string, expenseID string) (Expense, error)
	UpsertExpense(ctx context.Context, expense Expense) (Expense, error)
	RemoveExpense(ctx context.Context, userID string, expenseID string) (interface{}, error)
	GetSuppliesByProject(ctx context.Context, userID string, projectID string) ([]Supply, error)
	GetSupplyByID(ctx context.Context, userID string, supplyID string) (Supply, error)
	UpsertSupply(ctx context.Context, supply Supply) (interface{}, error)
	RemoveSupply(ctx context.Context, userID string, supplyID string) (interface{}, error)
}

type GenericDatabaseInfo struct {
	Created string `json:"created"`
	Updated string `json:"updated"`
}

type env struct {
	logger *zap.SugaredLogger `validate:"required"`
	validator *validator.Validate `validate:"required"`
	client *azcosmos.DatabaseClient `validate:"required"`
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
