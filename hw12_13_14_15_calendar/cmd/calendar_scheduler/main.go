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
	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/scheduler"
	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/setup"
	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/pkg/logger"
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

	storage, err := setup.SetupStorage(ctx, config, logg)
	if err != nil {
		logg.Error(err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
	defer storage.Close(ctx)

	period := time.Duration(config.Scheduler.PeriodSeconds) * time.Second
	calendarScheduler := scheduler.New(logg, storage, period)
	calendarScheduler.Schedule(ctx)
}
