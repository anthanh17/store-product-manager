package main

import (
	"log"

	"store-product-manager/configs"
	"store-product-manager/internal/dataaccess/cache"
	db "store-product-manager/internal/dataaccess/database/sqlc"
	"store-product-manager/internal/handler/http"
	"store-product-manager/internal/utils"

	"go.uber.org/zap"
)

func main() {
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	// Run server HTTP
	if err = runGinServer(config); err != nil {
		panic(err)
	}
}

func runGinServer(config configs.Config) error {
	// Logger
	logger, cleanup, err := utils.InitializeLogger(config.Log.Level)
	if err != nil {
		cleanup()
		logger.With(zap.Error(err)).Error("cannot initialize logger")
		return err
	}
	defer cleanup()

	// Database Accessor
	store, cleanupFunc, err := db.InitializeUpDB(config.Database, logger)
	if err != nil {
		cleanupFunc()
		logger.Info("error InitializeUpDB")
		return err
	}
	defer cleanupFunc()

	// Caching: in case using redis caching
	cacheMaker, err := cache.NewCachierClient(config.Cache, logger)
	if err != nil {
		logger.Info("error NewCachierClient")
		return err
	}

	// Gin Server
	server, err := http.NewServer(config, store, cacheMaker, logger)
	if err != nil {
		logger.Info("cannot create serve")
		return err
	}

	err = server.Start(config.HTTP.Address)
	if err != nil {
		logger.Info("cannot start serve")
		return err
	}

	return nil
}
