package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/config"
	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/rabbit"
	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/sender"
	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/setup"
	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/pkg/logger"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "../../configs/calendar_sender.toml", "Path to configuration file")
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

	rabbitConn, err := setup.Rabbit(config.Rabbit)
	if err != nil {
		logg.Error(err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
	defer rabbitConn.Close()

	rabbitConsumer := rabbit.NewConsumer(logg, rabbitConn)

	calendarSender := sender.New(logg, rabbitConsumer)
	calendarSender.Start(ctx)
}
