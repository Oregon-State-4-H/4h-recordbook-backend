package upc 

import (
	"4h-recordbook-backend/internal/config"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type Upc interface {
	GetProductByCode(code string) (Product, error)
}

type UpcClient struct {
	Endpoint string
	Key string
}

type env struct {
	logger *zap.SugaredLogger `validate:"required"`
	validator *validator.Validate `validate:"required"`
	client UpcClient `validate: "required"`
}

func New(logger *zap.SugaredLogger, cfg *config.Config) (Upc, error){

	logger.Info("Creating new UPC client")

	validate := validator.New()

	upcClient := UpcClient{
		Endpoint: cfg.Upc.Endpoint,
		Key: cfg.Upc.Current.Key,
	}

	e := &env{
		logger: logger,
		validator: validate,
		client: upcClient,
	}

	return e, nil

}