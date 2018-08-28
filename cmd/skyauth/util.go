package main

import (
	"context"
	"time"

	"github.com/skygeario/skygear-server/pkg/server/logging"
	"github.com/skygeario/skygear-server/pkg/server/skyconfig"
	"github.com/skygeario/skygear-server/pkg/server/skydb"
)

func ensureDB(config skyconfig.Configuration) func() (skydb.Conn, error) {
	logger := logging.LoggerEntryWithTag("main", "skydb")
	connOpener := func() (skydb.Conn, error) {
		return skydb.Open(
			context.Background(),
			config.DB.ImplName,
			config.App.Name,
			config.App.AccessControl,
			config.DB.Option,
			baseDBConfig(config),
		)
	}

	// Attempt to open connection to database. Retry for a number of
	// times before giving up.
	attempt := 0
	for {
		conn, connError := connOpener()
		if connError == nil {
			conn.Close()
			return connOpener
		}

		attempt++
		logger.Errorf("Failed to start skygear: %v", connError)
		if attempt >= 5 {
			logger.Fatalf("Failed to start skygear server because connection to database cannot be opened.")
		}

		logger.Info("Retrying in 1 second...")
		time.Sleep(time.Second * time.Duration(1))
	}
}

func initUserAuthRecordKeys(connOpener func() (skydb.Conn, error), authRecordKeys [][]string) {
	logger := logging.LoggerEntryWithTag("main", "auth")
	conn, err := connOpener()
	if err != nil {
		logger.Warnf("Failed to init user auth record keys: %v", err)
	}

	defer conn.Close()

	if err := conn.EnsureAuthRecordKeysExist(authRecordKeys); err != nil {
		panic(err)
	}

	if err := conn.EnsureAuthRecordKeysIndexesMatch(authRecordKeys); err != nil {
		panic(err)
	}
}

func baseDBConfig(config skyconfig.Configuration) skydb.DBConfig {
	passwordHistoryEnabled := config.UserAudit.PwHistorySize > 0 ||
		config.UserAudit.PwHistoryDays > 0

	return skydb.DBConfig{
		CanMigrate:             config.App.DevMode,
		PasswordHistoryEnabled: passwordHistoryEnabled,
	}
}
