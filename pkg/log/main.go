package log

import (
	"go.uber.org/zap"
)

type LoggerOptions struct {
	Level		zap.AtomicLevel
	OutputFile 	string
}

func New(params LoggerOptions) (*zap.SugaredLogger, error) {

	cfg := zap.NewDevelopmentConfig()

	cfg.Level = params.Level
	if params.OutputFile != "" {
		cfg.OutputPaths = append(cfg.OutputPaths, params.OutputFile)
	}

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	sugar := logger.Sugar()

	return sugar, nil

}