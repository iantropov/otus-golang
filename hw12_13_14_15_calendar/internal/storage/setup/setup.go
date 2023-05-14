package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/config"
	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/pkg/logger"
)

func Setup(ctx context.Context, config config.StorageConf, logg *logger.Logger) (storage.Storage, error) {
	if config.Type == "memory" {
		return memorystorage.New(), nil
	}

	if config.Type == "sql" {
		var db *sql.DB
		db, err := getSQLDb(config.DSN)
		if err != nil {
			return nil, fmt.Errorf("failed to get DB: %w", err)
		}

		sqlStorage := sqlstorage.New(logg, db)
		err = sqlStorage.Connect(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get sqlstorage: %w", err)
		}

		return sqlStorage, nil
	}

	return nil, errors.New("unsupported memory storage")
}

func getSQLDb(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
