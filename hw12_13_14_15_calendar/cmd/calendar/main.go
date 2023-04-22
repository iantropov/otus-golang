package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/app"
	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/config"
	internalhttp "github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/server/http"
	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/pkg/logger"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := config.NewConfig(configFile)
	if err != nil {
		log.Fatal("failed to get config", err)
	}

	logg, err := logger.New(config.Logger.Level)
	if err != nil {
		log.Fatal("failed to create logger", err)
	}

	var appStorage storage.Storage
	if config.Storage.Type == "memory" {
		appStorage = memorystorage.New()
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	if config.Storage.Type == "sql" {
		var db *sql.DB
		db, err = getSQLDb(config.Storage.DSN)
		if err != nil {
			logg.Error("failed to get DB: " + err.Error())
			cancel()
			os.Exit(1) //nolint:gocritic
		}

		sqlStorage := sqlstorage.New(logg, db)
		err = sqlStorage.Connect(ctx)
		if err != nil {
			logg.Error("failed to get sqlstorage: " + err.Error())
			cancel()
			os.Exit(1)
		}
		defer sqlStorage.Close(ctx)

		appStorage = sqlStorage
	}

	calendar := app.New(logg, appStorage)

	server := internalhttp.NewServer(config.HTTP.Host, config.HTTP.Port, logg, calendar)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1)
	}
}

func getSQLDb(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
