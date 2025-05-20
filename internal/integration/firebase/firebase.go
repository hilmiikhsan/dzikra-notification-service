package firebase

import (
	"context"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-notification-service/internal/infrastructure/config"
	"google.golang.org/api/option"
)

func InitFirebaseMessaging() *messaging.Client {
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatal().Err(err).Msg("Error getting current working directory")
	}

	credentialFilePath := filepath.Join(workingDir, config.Envs.FirebaseMessaging.CredentialServiceAccount)

	cfg := &firebase.Config{
		ProjectID: config.Envs.FirebaseMessaging.ProjectID,
	}

	opt := option.WithCredentialsFile(credentialFilePath)
	app, err := firebase.NewApp(context.Background(), cfg, opt)
	if err != nil {
		log.Fatal().Err(err).Msg("Error connecting to Firebase")
	}

	client, err := app.Messaging(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get firebase messaging client")
	}

	return client
}
