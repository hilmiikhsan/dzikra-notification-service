package cmd

import (
	"flag"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-notification-service/cmd/proto/notification"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-notification-service/internal/adapter"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-notification-service/internal/infrastructure"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-notification-service/internal/infrastructure/config"
	notificationService "github.com/Digitalkeun-Creative/be-dzikra-ecommerce-notification-service/internal/module/notification/handler/grpc"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-notification-service/pkg/validator"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func RunServeGRPC(cmd *flag.FlagSet, args []string) {
	var (
		envs        = config.Envs
		flagAppPort = cmd.String("port", envs.App.GrpcPort, "Application port")
		SERVER_PORT string
	)

	logLevel, err := zerolog.ParseLevel(envs.App.LogLevel)
	if err != nil {
		logLevel = zerolog.InfoLevel
	}

	if err := cmd.Parse(args); err != nil {
		log.Fatal().Err(err).Msg("Error while parsing flags")
	}

	if envs.App.GrpcPort != "" {
		SERVER_PORT = envs.App.GrpcPort
	} else {
		SERVER_PORT = *flagAppPort
	}

	grpcServer := grpc.NewServer()

	// Sync adapters
	err = adapter.Adapters.Sync(
		adapter.WithGRPCServer(grpcServer),
		adapter.WithDzikraPostgres(),
		adapter.WithValidator(validator.NewValidator()),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to sync adapters")
	}

	notification.RegisterNotificationServiceServer(grpcServer, notificationService.NewNotificationEmailAPI())

	infrastructure.InitializeLogger(
		envs.App.Environtment,
		envs.App.LogFile,
		logLevel,
	)

	lis, err := net.Listen("tcp", ":"+SERVER_PORT)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to listen on gRPC port")
	}

	// Run gRPC server in a goroutine
	go func() {
		log.Info().Msgf("gRPC server is running on port %s", envs.App.GrpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal().Err(err).Msg("Failed to serve gRPC")
		}
	}()

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	shutdownSignals := []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGINT}
	if runtime.GOOS == "windows" {
		shutdownSignals = []os.Signal{os.Interrupt}
	}
	signal.Notify(quit, shutdownSignals...)
	<-quit

	log.Info().Msg("gRPC server is shutting down ...")
	grpcServer.GracefulStop()

	err = adapter.Adapters.Unsync()
	if err != nil {
		log.Error().Err(err).Msg("Error while closing adapters")
	}

	log.Info().Msg("gRPC server gracefully stopped")
}
