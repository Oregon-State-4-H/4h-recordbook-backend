package main

import (

    "flag"
    "4h-recordbook-backend/internal/api"
    "4h-recordbook-backend/internal/config"
    "4h-recordbook-backend/pkg/db"
    "4h-recordbook-backend/pkg/log"
    "go.uber.org/zap"

)

func main() {

    debug := flag.Bool("d", false, "enable debug mode")
    logFile := flag.String("l", "", "log file")
    flag.Parse()

    logOptions := log.LoggerOptions{
        Level:      zap.NewAtomicLevelAt(zap.InfoLevel),
        OutputFile: *logFile,
    }
    if *debug {
        logOptions.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
    }

    logger, err := log.New(logOptions)
    if err != nil {
        panic(err)
    }
    defer logger.Sync()
    
    cfg, err := config.New(logger)
    if err != nil {
        panic(err)
    }

    dbInstance, err := db.New(logger, cfg)
    if err != nil {
        panic(err)
    }

    apiInstance, err := api.New(logger, cfg, dbInstance)
    if err != nil {
        panic(err)
    }

    err = apiInstance.RunLocal()
    if err != nil {
        panic(err)
    }

}