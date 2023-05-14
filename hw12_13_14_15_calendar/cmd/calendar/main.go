package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/app"
	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/config"
	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/server"
	internalgrpc "github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/server/http"
	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/internal/setup"
	"github.com/iantropov/otus-golang/hw12_13_14_15_calendar/pkg/logger"
	_ "github.com/lib/pq"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "../../configs/calendar.memory.toml", "Path to configuration file")
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

	storage, err := setup.Storage(ctx, config.Storage, logg)
	if err != nil {
		logg.Error(err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
	defer storage.Close(ctx)

	calendar := app.New(logg, storage)

	serversWg := sync.WaitGroup{}
	serversWg.Add(2)

	httpServer := internalhttp.NewServer(config.HTTP.Host, config.HTTP.Port, logg, calendar)
	grpcServer := internalgrpc.NewServer(config.GRPC.Host, config.GRPC.Port, logg, calendar)

	startServer(ctx, cancel, httpServer, "http", logg, &serversWg)
	startServer(ctx, cancel, grpcServer, "grpc", logg, &serversWg)

	serversWg.Wait()
}

func startServer(
	ctx context.Context,
	cancel func(),
	server server.Server,
	serverName string,
	logg *logger.Logger,
	wg *sync.WaitGroup,
) {
	go func() {
		defer wg.Done()

		<-ctx.Done()

		stopCtx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(stopCtx); err != nil {
			logg.Errorf("failed to stop %s server: %v\n", serverName, err)
		}
	}()
	go func() {
		logg.Infof("calendar %s api is running...\n", serverName)

		if err := server.Start(); err != nil {
			logg.Errorf("%s server: %v\n", serverName, err)
			cancel()
		}
	}()
}
