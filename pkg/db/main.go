package db 

import (
	"4h-recordbook-backend/internal/config"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type Db interface {}

type env struct {
	logger	  *zap.SugaredLogger  `validate:"required"`
	validator *validator.Validate `validate:"required"`
	client	  *azcosmos.Client	  `validate:"required"`
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

	e := &env{
		logger:	   logger,
		validator: validate,
		client:    client,
	}

	return e, nil

}
