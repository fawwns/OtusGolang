package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os/signal"
	"syscall"
	"time"

	pb "github.com/fawwns/OtusGolang/hw12_13_14_15_calendar/api/proto"
	"github.com/fawwns/OtusGolang/hw12_13_14_15_calendar/internal/app"
	"github.com/fawwns/OtusGolang/hw12_13_14_15_calendar/internal/logger"
	"github.com/fawwns/OtusGolang/hw12_13_14_15_calendar/internal/server/grpcserver"
	internalhttp "github.com/fawwns/OtusGolang/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/fawwns/OtusGolang/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/fawwns/OtusGolang/hw12_13_14_15_calendar/internal/storage/sql"
	"google.golang.org/grpc"
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
	unaryInterceptor := grpcserver.UnaryLoggingInterceptor(logg)

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP,
	)
	defer cancel()

	var storage app.Storage

	switch config.Storage.Type {
	case "sql":
		sqlStorage := sqlstorage.New(config.Storage.DSN)
		if err := sqlStorage.Connect(ctx); err != nil {
			logg.Error("failed to connect to SQL DB: " + err.Error())
			return
		}
		storage = sqlStorage

	case "inmemory":
		storage = memorystorage.New()

	default:
		logg.Error("unknown storage type in config: " + config.Storage.Type)
		return
	}

	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(logg, calendar, config.Server.Host, config.Server.Port)

	go func() {
		if err := server.Start(ctx); err != nil {
			logg.Error("HTTP server error: " + err.Error())
			cancel()
		}
	}()

	addrGRPC := fmt.Sprintf("%s:%d", config.GRPC.Host, config.GRPC.Port)
	logg.Info("Starting gRPC server at " + addrGRPC)

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(unaryInterceptor))
	grpcHandler := grpcserver.New(calendar)

	grpcListener, err := net.Listen("tcp", addrGRPC)
	if err != nil {
		logg.Error("failed to listen grpc: " + err.Error())
		return
	}

	pb.RegisterCalendarServiceServer(grpcServer, grpcHandler)

	go func() {
		logg.Info("gRPC server on " + addrGRPC)
		if err := grpcServer.Serve(grpcListener); err != nil {
			logg.Error("gRPC error: " + err.Error())
			cancel()
		}
	}()

	logg.Info("calendar is running...")

	<-ctx.Done()
	logg.Info("shutdown...")

	httpCtx, httpCancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer httpCancel()
	server.Stop(httpCtx)

	grpcServer.GracefulStop()
}
