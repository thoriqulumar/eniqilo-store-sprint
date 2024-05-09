package main

import (
	"context"
	"eniqilo-store/config"
	"eniqilo-store/database"
	logger "eniqilo-store/pkg/log"
	"eniqilo-store/server"
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	ctx := context.Background()

	cfg, err := config.LoadConfig(ctx)
	if err != nil {
		panic(err)
	}

	logger, err := logger.NewLogger(
		zapcore.DebugLevel,
		"eniqilo_store",
		"1",
	)
	if err != nil {
		panic(err)
	}

	db, err := database.NewDatabase(cfg)
	if err != nil {
		logger.Error("error opening database", zap.Error(err))
		panic(err)
	}
	defer db.Close()
	db.SetMaxIdleConns(80)

	s := server.NewServer(db)
	s.RegisterRoute(cfg)

	log.Fatal(s.Run())
}
