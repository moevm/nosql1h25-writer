package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type (
	Config struct {
		App    `yaml:"app"`
		HTTP   `yaml:"http"`
		Log    `yaml:"logger"`
		JWT    `yaml:"jwt"`
		Hasher `yaml:"hasher"`
	}

	App struct {
		Name    string `env-required:"true" yaml:"name" env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"level" env:"LOG_LEVEL"`
	}

	JWT struct {
		SignKey  string        `env-required:"true" env:"JWT_SIGN_KEY"`
		TokenTTL time.Duration `env-required:"true" yaml:"token_ttl" env:"JWT_TOKEN_TTL"`
	}

	Hasher struct {
		Salt string `env-required:"true" env:"HASHER_SALT"`
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
