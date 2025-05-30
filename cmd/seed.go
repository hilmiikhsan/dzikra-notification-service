package cmd

import (
	"flag"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-notification-service/db/seeds"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-notification-service/internal/adapter"
	"github.com/rs/zerolog/log"
)

func RunSeed(cmd *flag.FlagSet, args []string) {
	var (
		table = cmd.String("table", "", "seed to run")
	)

	if err := cmd.Parse(args); err != nil {
		log.Fatal().Err(err).Msg("Error while parsing flags")
	}

	adapter.Adapters.Sync(
		adapter.WithDzikraPostgres(),
	)

	defer func() {
		if err := adapter.Adapters.Unsync(); err != nil {
			log.Fatal().Err(err).Msg("Error while closing database connection")
		}
	}()

	seeds.Execute(adapter.Adapters.DzikraPostgres, *table)
}
