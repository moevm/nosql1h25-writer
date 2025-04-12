package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App   App   `yaml:"app"`
		HTTP  HTTP  `yaml:"http"`
		Log   Log   `yaml:"logger"`
		Mongo Mongo `yaml:"mongo"`
		Auth  Auth  `yaml:"auth"`
	}

	App struct {
		Name    string `env-required:"true" yaml:"name" env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	Mongo struct {
		Uri             string        `env-required:"true" yaml:"uri" env:"MONGO_URI"`
		ShutdownTimeout time.Duration `env-required:"true" yaml:"shutdown_timeout" env:"MONGO_SHUTDOWN_TIMEOUT"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"level" env:"LOG_LEVEL"`
	}

	Auth struct {
		JWTSecretKey    string        `env-required:"true" yaml:"jwt_secret_key" env:"AUTH_JWT_SECRET_KEY"`
		AccessTokenTTL  time.Duration `env_required:"true" yaml:"access_token_ttl" env:"AUTH_ACCESS_TOKEN_TTL"`
		RefreshTokenTTL time.Duration `env_required:"true" yaml:"refresh_token_ttl" env:"AUTH_REFRESH_TOKEN_TTL"`
	}
)

func New(configPath string) (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig(configPath, cfg)
	if err != nil {
		return nil, fmt.Errorf("config - NewConfig - cleanenv.ReadConfig: %w", err)
	}

	err = cleanenv.UpdateEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("config - NewConfig - cleanenv.UpdateEnv: %w", err)
	}

	return cfg, nil
}
