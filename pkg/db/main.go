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
