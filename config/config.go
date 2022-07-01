package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		App     `yaml:"app"`
		HTTP    `yaml:"http"`
		GRPC    `yaml:"grpc"`
		Mongo   `yaml:"mongo"`
		Log     `yaml:"logger"`
		Users   []User `yaml:"users"`
		Profile `yaml:"profile"`
	}

	// App -.
	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	// HTTP -.
	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	// GRPC -.
	GRPC struct {
		Port string `env-required:"true" yaml:"port" env:"GRPC_PORT"`
	}

	Mongo struct {
		Dsn           string `env-required:"true" yaml:"dsn" env:"MONGO_DSN"`
		DbName        string `env-required:"true" yaml:"db_name" env:"DB_NAME"`
		MigrationPath string `env-required:"true" yaml:"migration_path" env:"MIGRATION_PATH"`
		MigrationRun  bool   `yaml:"migration_run" env:"MIGRATION_RUN" env-default:"false"`
		MigrationMode string `env-required:"true" yaml:"migration_mode" env:"MIGRATION_MODE"`
	}

	// Log -.
	Log struct {
		Level string `env-required:"true" yaml:"log_level"   env:"LOG_LEVEL"`
	}

	// User -.
	User struct {
		Name     string `yaml:"name"`
		Password string `yaml:"password"`
	}

	// Profile -.
	Profile struct {
		Login    string `yaml:"login"`
		Password string `yaml:"password"`
		Enabled  bool   `yaml:"enabled"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
