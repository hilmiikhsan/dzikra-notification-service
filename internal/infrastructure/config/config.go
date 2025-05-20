package config

import (
	"sync"

	"github.com/rs/zerolog/log"

	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-notification-service/pkg/config"
	"github.com/Digitalkeun-Creative/be-dzikra-ecommerce-notification-service/pkg/utils"
)

var (
	Envs *Config // Envs is global vars Config.
	once sync.Once
)

type Config struct {
	App struct {
		Name                    string `env:"APP_NAME" env-default:"dzikra-notification-service"`
		Environtment            string `env:"APP_ENV" env-default:"development"`
		BaseURL                 string `env:"APP_BASE_URL" env-default:"http://localhost:9091"`
		Port                    string `env:"APP_PORT" env-default:"9091"`
		GrpcPort                string `env:"APP_GRPC_PORT" env-default:"7001"`
		LogLevel                string `env:"APP_LOG_LEVEL" env-default:"debug"`
		LogFile                 string `env:"APP_LOG_FILE" env-default:"./logs/app.log"`
		LogFileWs               string `env:"APP_LOG_FILE_WS" env-default:"./logs/ws.log"`
		LocalStoragePublicPath  string `env:"LOCAL_STORAGE_PUBLIC_PATH" env-default:"./storage/public"`
		LocalStoragePrivatePath string `env:"LOCAL_STORAGE_PRIVATE_PATH" env-default:"./storage/private"`
	}
	DB struct {
		ConnectionTimeout int `env:"DB_CONN_TIMEOUT" env-default:"30" env-description:"database timeout in seconds"`
		MaxOpenCons       int `env:"DB_MAX_OPEN_CONS" env-default:"20" env-description:"database max open conn in seconds"`
		MaxIdleCons       int `env:"DB_MAX_IdLE_CONS" env-default:"20" env-description:"database max idle conn in seconds"`
		ConnMaxLifetime   int `env:"DB_CONN_MAX_LIFETIME" env-default:"0" env-description:"database conn max lifetime in seconds"`
	}
	Notification struct {
		MailHost     string `env:"MAIL_HOST" env-default:"smtp.gmail.com"`
		MailPort     string `env:"MAIL_PORT" env-default:"587"`
		MailUser     string `env:"MAIL_USER" env-default:""`
		MailPassword string `env:"MAIL_PASSWORD" env-default:""`
	}
	DzikraPostgres struct {
		Host     string `env:"DZIKRA_POSTGRES_HOST" env-default:"localhost"`
		Port     string `env:"DZIKRA_POSTGRES_PORT" env-default:"5432"`
		Username string `env:"DZIKRA_POSTGRES_USER" env-default:"postgres"`
		Password string `env:"DZIKRA_POSTGRES_PASSWORD" env-default:"postgres"`
		Database string `env:"DZIKRA_POSTGRES_DB" env-default:"dzikra_notifications"`
		SslMode  string `env:"DZIKRA_POSTGRES_SSL_MODE" env-default:"disable"`
	}
	RedisDB struct {
		Host     string `env:"DZIKRA_REDIS_HOST" env-default:"redis"`
		Port     string `env:"DZIKRA_REDIS_PORT" env-default:"6379"`
		Password string `env:"DZIKRA_REDIS_PASSWORD" env-default:"password"`
		Database int    `env:"DZIKRA_REDIS_DB" env-default:"0"`
	}
	FirebaseMessaging struct {
		CredentialServiceAccount string `env:"FIREBASE_MESSAGING_CREDENTIAL_SERVICE_ACCOUNT" env-default:""`
		ProjectID                string `env:"FIREBASE_MESSAGING_PROJECT_ID" env-default:""`
	}
}

// Option is Configure type return func.
type Option = func(c *Configure) error

// Configure is the data struct.
type Configure struct {
	path     string
	filename string
}

// Configuration create instance.
func Configuration(opts ...Option) *Configure {
	c := &Configure{}

	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			panic(err)
		}
	}
	return c
}

// Initialize will create instance of Configure.
func (c *Configure) Initialize() {
	once.Do(func() {
		Envs = &Config{}
		if err := config.Load(config.Opts{
			Config:    Envs,
			Paths:     []string{c.path},
			Filenames: []string{c.filename},
		}); err != nil {
			log.Fatal().Err(err).Msg("get config error")
		}

		Envs.App.Name = utils.GetEnv("APP_NAME", Envs.App.Name)
		Envs.App.Port = utils.GetEnv("APP_PORT", Envs.App.Port)
		Envs.App.GrpcPort = utils.GetEnv("APP_GRPC_PORT", Envs.App.GrpcPort)
		Envs.App.LogLevel = utils.GetEnv("APP_LOG_LEVEL", Envs.App.LogLevel)
		Envs.App.LogFile = utils.GetEnv("APP_LOG_FILE", Envs.App.LogFile)
		Envs.App.LogFileWs = utils.GetEnv("APP_LOG_FILE_WS", Envs.App.LogFileWs)
		Envs.App.LocalStoragePublicPath = utils.GetEnv("LOCAL_STORAGE_PUBLIC_PATH", Envs.App.LocalStoragePublicPath)
		Envs.App.LocalStoragePrivatePath = utils.GetEnv("LOCAL_STORAGE_PRIVATE_PATH", Envs.App.LocalStoragePrivatePath)
		Envs.DB.ConnectionTimeout = utils.GetIntEnv("DB_CONN_TIMEOUT", Envs.DB.ConnectionTimeout)
		Envs.DB.MaxOpenCons = utils.GetIntEnv("DB_MAX_OPEN_CONS", Envs.DB.MaxOpenCons)
		Envs.DB.MaxIdleCons = utils.GetIntEnv("DB_MAX_IdLE_CONS", Envs.DB.MaxIdleCons)
		Envs.DB.ConnMaxLifetime = utils.GetIntEnv("DB_CONN_MAX_LIFETIME", Envs.DB.ConnMaxLifetime)
		Envs.Notification.MailHost = utils.GetEnv("MAIL_HOST", Envs.Notification.MailHost)
		Envs.Notification.MailPort = utils.GetEnv("MAIL_PORT", Envs.Notification.MailPort)
		Envs.Notification.MailUser = utils.GetEnv("MAIL_USER", Envs.Notification.MailUser)
		Envs.Notification.MailPassword = utils.GetEnv("MAIL_PASSWORD", Envs.Notification.MailPassword)
		Envs.DzikraPostgres.Host = utils.GetEnv("DZIKRA_POSTGRES_HOST", Envs.DzikraPostgres.Host)
		Envs.DzikraPostgres.Port = utils.GetEnv("DZIKRA_POSTGRES_PORT", Envs.DzikraPostgres.Port)
		Envs.DzikraPostgres.Username = utils.GetEnv("DZIKRA_POSTGRES_USER", Envs.DzikraPostgres.Username)
		Envs.DzikraPostgres.Password = utils.GetEnv("DZIKRA_POSTGRES_PASSWORD", Envs.DzikraPostgres.Password)
		Envs.DzikraPostgres.Database = utils.GetEnv("DZIKRA_POSTGRES_DB", Envs.DzikraPostgres.Database)
		Envs.DzikraPostgres.SslMode = utils.GetEnv("DZIKRA_POSTGRES_SSL_MODE", Envs.DzikraPostgres.SslMode)
		Envs.RedisDB.Host = utils.GetEnv("DZIKRA_REDIS_HOST", Envs.RedisDB.Host)
		Envs.RedisDB.Port = utils.GetEnv("DZIKRA_REDIS_PORT", Envs.RedisDB.Port)
		Envs.RedisDB.Password = utils.GetEnv("DZIKRA_REDIS_PASSWORD", Envs.RedisDB.Password)
		Envs.RedisDB.Database = utils.GetIntEnv("DZIKRA_REDIS_DB", Envs.RedisDB.Database)
		Envs.FirebaseMessaging.CredentialServiceAccount = utils.GetEnv("FIREBASE_MESSAGING_CREDENTIAL_SERVICE_ACCOUNT", Envs.FirebaseMessaging.CredentialServiceAccount)
	})
}

// WithPath will assign to field path Configure.
func WithPath(path string) Option {
	return func(c *Configure) error {
		c.path = path
		return nil
	}
}

// WithFilename will assign to field name Configure.
func WithFilename(name string) Option {
	return func(c *Configure) error {
		c.filename = name
		return nil
	}
}
