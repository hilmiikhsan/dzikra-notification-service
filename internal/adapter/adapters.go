package adapter

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

var (
	Adapters *Adapter
)

type Option func(adapter *Adapter)

type Validator interface {
	Validate(i any) error
}

type Adapter struct {
	// Driving Adapters
	RestServer *fiber.App
	GRPCServer *grpc.Server

	//Driven Adapters
	DzikraPostgres *sqlx.DB
	Validator      Validator // *validator.Validator
}

func (a *Adapter) Sync(opts ...Option) error {
	var errs []string

	for _, opt := range opts {
		opt(a)
	}

	if a.DzikraPostgres == nil {
		errs = append(errs, "Dzikra Postgres not initialized")
	}

	if a.GRPCServer == nil && a.RestServer == nil {
		errs = append(errs, "No server initialized")
	}

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}

	return nil
}

func (a *Adapter) Unsync() error {
	var errs []string

	if a.RestServer != nil {
		if err := a.RestServer.Shutdown(); err != nil {
			errs = append(errs, err.Error())
		}
		log.Info().Msg("Rest server disconnected")
	}

	if a.GRPCServer != nil {
		if a.GRPCServer != nil {
			a.GRPCServer.GracefulStop()
		}
		log.Info().Msg("gRPC server disconnected")
	}

	if a.DzikraPostgres != nil {
		if err := a.DzikraPostgres.Close(); err != nil {
			errs = append(errs, err.Error())
		}
		log.Info().Msg("Dzikra Postgres disconnected")
	}

	if len(errs) > 0 {
		err := errors.New(strings.Join(errs, "\n"))
		log.Error().Msgf("Error while disconnecting adapters: %v", err)
		return err
	}

	return nil
}
