package sqlstorage

import (
	"context"
	"fmt"
	"time"

	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/logger"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
)

func (s *Storage) Connect(
	ctx context.Context,
	logg logger.Logger,
	dbHost string,
	dbPort int,
	dbUser, dbPassword, dbName string,
) (err error) {
	s.pgxPool, err = getPgxpool(ctx, logg, dbHost, dbPort, dbUser, dbPassword, dbName)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	s.pgxPool.Close()
	return nil
}

func getPgxpool(
	ctx context.Context,
	logg logger.Logger,
	dbHost string,
	dbPort int,
	dbUser, dbPassword, dbName string,
) (*pgxpool.Pool, error) {
	psqlConn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost,
		dbPort,
		dbUser,
		dbPassword,
		dbName,
	)
	pgxConfig, err := pgxpool.ParseConfig(psqlConn)
	if err != nil {
		return nil, err
	}

	pgxConfig.MaxConnIdleTime = time.Minute
	pgxConfig.MaxConnLifetime = time.Hour
	pgxConfig.MinConns = 2
	pgxConfig.MaxConns = 10

	pgxConfig.ConnConfig.Logger = logg
	pgxConfig.ConnConfig.LogLevel = pgx.LogLevelDebug

	pool, err := pgxpool.ConnectConfig(ctx, pgxConfig)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return pool, nil
}
