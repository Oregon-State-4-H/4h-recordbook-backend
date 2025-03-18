package config

import (
	_ "embed"
	"encoding/json"
	"os"

	"go.uber.org/zap"
)

//go:embed config.json
var configJSON []byte

const (
	MAX_PAGE_SIZE  = 500
	PRODUCTION_ENV = "PRODUCTION"
)

type Config struct {
	MaxPageSize int      `json:"max_page_size"`
	Database    Database `json:"cosmos"`
	Upc         Upc      `json:"upc"`
}

type Database struct {
	Current     DatabaseParams
	Production  DatabaseParams `json:"production"`
	Development DatabaseParams `json:"development"`
}

type DatabaseParams struct {
	Endpoint string `json:"endpoint"`
	Key      string `json:"key"`
}

type Upc struct {
	Endpoint    string `json:"endpoint"`
	Current     UpcParams
	Production  UpcParams `json:"production"`
	Development UpcParams `json:"development"`
}

type UpcParams struct {
	Key string `json:"key"`
}

func New(logger *zap.SugaredLogger) (*Config, error) {

	logger.Info("Setting up config")

	var c Config

	err := json.Unmarshal(configJSON, &c)
	if err != nil {
		logger.Errorf("Failed to unmarshal config: %v", err)
		return nil, err
	}

	c.MaxPageSize = MAX_PAGE_SIZE

	env := os.Getenv("APP_ENV")
	if env == PRODUCTION_ENV {
		logger.Debug("Running with production config")
		c.Database.Current = c.Database.Production
		c.Upc.Current = c.Upc.Production
	} else {
		logger.Debug("Running with development config")
		c.Database.Current = c.Database.Development
		c.Upc.Current = c.Upc.Development
	}

	return &c, nil

}
