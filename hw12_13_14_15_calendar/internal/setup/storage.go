package setup

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/config"
	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/pkg/logger"
)

func SetupStorage(ctx context.Context, config config.Config, logg *logger.Logger) (storage.Storage, error) {
	var appStorage storage.Storage
	if config.Storage.Type == "memory" {
		appStorage = memorystorage.New()
	}

	if config.Storage.Type == "sql" {
		var db *sql.DB
		db, err := getSQLDb(config.Storage.DSN)
		if err != nil {
			return nil, fmt.Errorf("failed to get DB: %w", err)
		}

		sqlStorage := sqlstorage.New(logg, db)
		err = sqlStorage.Connect(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get sqlstorage: %w", err)
		}
		defer sqlStorage.Close(ctx)

		appStorage = sqlStorage
	}

	return appStorage, nil
}

func getSQLDb(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
