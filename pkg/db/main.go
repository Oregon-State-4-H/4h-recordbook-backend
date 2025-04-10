package db

import (
	"4h-recordbook-backend/internal/config"
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type Db interface {
	GetUser(context.Context, string) (User, error)
	UpsertUser(context.Context, User) (interface{}, error)
	GetBookmarkByLink(context.Context, string, string) (Bookmark, error)
	GetBookmarks(context.Context, string) ([]Bookmark, error)
	AddBookmark(context.Context, Bookmark) (Bookmark, error)
	RemoveBookmark(context.Context, string, string) (interface{}, error)
	GetProjectByID(context.Context, string, string) (Project, error)
	GetCurrentProjects(context.Context, string) ([]Project, error)
	GetProjectsByUser(context.Context, string) ([]Project, error)
	UpsertProject(context.Context, Project) (Project, error)
	RemoveProject(context.Context, string, string) (interface{}, error)
	GetResume(context.Context, string) (Resume, error)
	GetSection1ByID(context.Context, string, string) (Section1, error)
	GetSection2ByID(context.Context, string, string) (Section2, error)
	GetSection3ByID(context.Context, string, string) (Section3, error)
	GetSection4ByID(context.Context, string, string) (Section4, error)
	GetSection5ByID(context.Context, string, string) (Section5, error)
	GetSection6ByID(context.Context, string, string) (Section6, error)
	GetSection7ByID(context.Context, string, string) (Section7, error)
	GetSection8ByID(context.Context, string, string) (Section8, error)
	GetSection9ByID(context.Context, string, string) (Section9, error)
	GetSection10ByID(context.Context, string, string) (Section10, error)
	GetSection11ByID(context.Context, string, string) (Section11, error)
	GetSection12ByID(context.Context, string, string) (Section12, error)
	GetSection13ByID(context.Context, string, string) (Section13, error)
	GetSection14ByID(context.Context, string, string) (Section14, error)
	GetSection1sByUser(context.Context, string) ([]Section1, error)
	GetSection2sByUser(context.Context, string) ([]Section2, error)
	GetSection3sByUser(context.Context, string) ([]Section3, error)
	GetSection4sByUser(context.Context, string) ([]Section4, error)
	GetSection5sByUser(context.Context, string) ([]Section5, error)
	GetSection6sByUser(context.Context, string) ([]Section6, error)
	GetSection7sByUser(context.Context, string) ([]Section7, error)
	GetSection8sByUser(context.Context, string) ([]Section8, error)
	GetSection9sByUser(context.Context, string) ([]Section9, error)
	GetSection10sByUser(context.Context, string) ([]Section10, error)
	GetSection11sByUser(context.Context, string) ([]Section11, error)
	GetSection12sByUser(context.Context, string) ([]Section12, error)
	GetSection13sByUser(context.Context, string) ([]Section13, error)
	GetSection14sByUser(context.Context, string) ([]Section14, error)
	UpsertSection1(context.Context, Section1) (Section1, error)
	UpsertSection2(context.Context, Section2) (Section2, error)
	UpsertSection3(context.Context, Section3) (Section3, error)
	UpsertSection4(context.Context, Section4) (Section4, error)
	UpsertSection5(context.Context, Section5) (Section5, error)
	UpsertSection6(context.Context, Section6) (Section6, error)
	UpsertSection7(context.Context, Section7) (Section7, error)
	UpsertSection8(context.Context, Section8) (Section8, error)
	UpsertSection9(context.Context, Section9) (Section9, error)
	UpsertSection10(context.Context, Section10) (Section10, error)
	UpsertSection11(context.Context, Section11) (Section11, error)
	UpsertSection12(context.Context, Section12) (Section12, error)
	UpsertSection13(context.Context, Section13) (Section13, error)
	UpsertSection14(context.Context, Section14) (Section14, error)
	RemoveSection(context.Context, string, string) (interface{}, error)
	GetEventsByUser(context.Context, string) ([]Event, error)
	GetEventByID(context.Context, string, string) (Event, error)
	UpsertEvent(context.Context, Event) (Event, error)
	RemoveEvent(context.Context, string, string) (interface{}, error)
	GetAnimalsByProject(context.Context, string, string) ([]Animal, error)
	GetAnimalByID(context.Context, string, string) (Animal, error)
	UpsertAnimal(context.Context, Animal) (Animal, error)
	RemoveAnimal(context.Context, string, string) (interface{}, error)
	GetFeedsByProject(context.Context, string, string) ([]Feed, error)
	GetFeedByID(context.Context, string, string) (Feed, error)
	UpsertFeed(context.Context, Feed) (Feed, error)
	RemoveFeed(context.Context, string, string) (interface{}, error)
	GetFeedPurchasesByProject(context.Context, string, string) ([]FeedPurchase, error)
	GetFeedPurchaseByID(context.Context, string, string) (FeedPurchase, error)
	UpsertFeedPurchase(context.Context, FeedPurchase) (FeedPurchase, error)
	RemoveFeedPurchase(context.Context, string, string) (interface{}, error)
	GetDailyFeedsByProjectAndAnimal(context.Context, string, string, string) ([]DailyFeed, error)
	GetDailyFeedByID(context.Context, string, string) (DailyFeed, error)
	UpsertDailyFeed(context.Context, DailyFeed) (DailyFeed, error)
	RemoveDailyFeed(context.Context, string, string) (interface{}, error)
	GetExpensesByProject(context.Context, string, string) ([]Expense, error)
	GetExpenseByID(context.Context, string, string) (Expense, error)
	UpsertExpense(context.Context, Expense) (Expense, error)
	RemoveExpense(context.Context, string, string) (interface{}, error)
	GetSuppliesByProject(context.Context, string, string) ([]Supply, error)
	GetSupplyByID(context.Context, string, string) (Supply, error)
	UpsertSupply(context.Context, Supply) (Supply, error)
	RemoveSupply(context.Context, string, string) (interface{}, error)
}

type GenericDatabaseInfo struct {
	Created string `json:"created"`
	Updated string `json:"updated"`
}

type env struct {
	logger    *zap.SugaredLogger       `validate:"required"`
	validator *validator.Validate      `validate:"required"`
	client    *azcosmos.DatabaseClient `validate:"required"`
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
		logger:    logger,
		validator: validate,
		client:    dbClient,
	}

	return e, nil

}
