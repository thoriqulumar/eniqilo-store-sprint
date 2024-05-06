package main

import (
	"eniqilo-store/config"
	"eniqilo-store/database"
	pkgLog "eniqilo-store/log"
	"eniqilo-store/pkg/version"
	"eniqilo-store/server"
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	config.LoadConfig(".env")

	logger, err := pkgLog.New(
		zapcore.DebugLevel,
		version.ServiceID,
		version.Version,
	)
	if err != nil {
		panic(err)
	}

	db, err := database.NewDatabase()
	if err != nil {
		logger.Error("error opening database", zap.Error(err))
		panic(err)
	}
	defer db.Close()
	db.SetMaxIdleConns(80)

	s := server.NewServer(db)
	s.RegisterRoute()

	log.Fatal(s.Run())
}
