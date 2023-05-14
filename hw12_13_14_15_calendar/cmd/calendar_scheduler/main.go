package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/config"
	setupQueue "github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/queue/setup"
	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/scheduler"
	setupStorage "github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/storage/setup"
	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/pkg/logger"
	_ "github.com/lib/pq"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "../../configs/config.memory.toml", "Path to configuration file")
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

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	storage, err := setupStorage.Setup(ctx, config.Storage, logg)
	if err != nil {
		logg.Error(err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
	defer storage.Close(ctx)

	queueConn, err := setupQueue.Setup(config.Queue)
	if err != nil {
		logg.Error(err.Error())
		cancel()
		os.Exit(1)
	}
	defer queueConn.Close()

	producer := queueConn.NewProducer(logg)

	calendarScheduler := scheduler.New(logg, storage, producer, getSchedulerPeriod(config))
	calendarScheduler.Schedule(ctx)
}

func getSchedulerPeriod(config config.Config) time.Duration {
	return time.Duration(config.Scheduler.PeriodSeconds) * time.Second
}
