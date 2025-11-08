package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fawwns/OtusGolang/hw12_13_14_15_calendar/internal/app"
	"github.com/fawwns/OtusGolang/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/fawwns/OtusGolang/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/fawwns/OtusGolang/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/fawwns/OtusGolang/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := LoadConfig(configFile)
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	logg := logger.New(config.Logger.Level)
	ctx := context.Background()
	var storage app.Storage

	switch config.Storage.Active {
	case "sql":
		sqlStorage := sqlstorage.New(config.Storage.SQL.DSN)
		if err := sqlStorage.Connect(ctx); err != nil {
			logg.Error("failed to connect to SQL DB: " + err.Error())
			return
		}
		storage = sqlStorage

	case "inmemory":
		storage = memorystorage.New()

	default:
		logg.Error("unknown storage type in config: " + config.Storage.Active)
		return
	}

	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(logg, calendar, config.Server.Host, config.Server.Port)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

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
