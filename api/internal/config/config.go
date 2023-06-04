package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"go.uber.org/fx"
)

const (
	EnvironmentDevelopment = "development"
	EnvironmentStaging     = "staging"
	EnvironmentProduction  = "production"
)

type UserDetails struct {
	Password string
	Role     string
}

type Config struct {
	// Generic
	Environment  string `envconfig:"ENV" default:"staging"`
	Port         string `envconfig:"APP_PORT" default:"8080"`
	AllowedHosts string `envconfig:"ALLOWED_HOSTS" default:"*"`

	// Logging
	LoggerType  string `envconfig:"LOGGER_TYPE" required:"false" default:"zap"`
	LoggerLevel int    `envconfig:"LOGGER_LEVEL" required:"false" default:"2"`

	// MySQL (Internal)
	MySQLHost     string `envconfig:"MYSQL_HOST" required:"false"`
	MySQLPort     string `envconfig:"MYSQL_PORT" required:"false"`
	MySQLUser     string `envconfig:"MYSQL_USER" required:"false" secret:"gd_mysql_user"`
	MySQLPassword string `envconfig:"MYSQL_PASSWORD" secret:"gd_mysql_password"`
	MySQLDBName   string `envconfig:"MYSQL_DB_NAME"`

	// Sentry
	SentryEnabled bool   `envconfig:"SENTRY_ENABLED" required:"false" default:"false"`
	SentryDSN     string `envconfig:"SENTRY_DSN" required:"false"`
}

func ProvideConfig(cfgFile string) fx.Option {
	return fx.Provide(func() (*Config, error) {
		cfg, err := loadConfig(cfgFile)
		return cfg, err
	})
}

func loadConfig(cfgFile string) (cfg *Config, err error) {
	if cfgFile != "" {
		if cfg, err = readCfgFromFile(cfgFile); err != nil {
			return nil, err
		}
		return cfg, nil
	}
	if cfg, err = readCfgFromEnv(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func readCfgFromFile(cfgFile string) (*Config, error) {
	if err := godotenv.Load(cfgFile); err != nil {
		return nil, errors.WithStack(err)
	}
	return readCfgFromEnv()
}

func readCfgFromEnv() (*Config, error) {
	cfg := Config{}
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, errors.WithStack(err)
	}
	return &cfg, nil
}

func (c *Config) MySQLUrl() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		c.MySQLUser,
		c.MySQLPassword,
		c.MySQLHost,
		c.MySQLPort,
		c.MySQLDBName,
	)
}
