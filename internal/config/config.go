package config

import (
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Env                   string        `default:"development" split_words:"true"`
	ServerAddr            string        `default:":8080" split_words:"true" validate:"required"`
	AllowedOrigins        []string      `default:"http://localhost:5173,http://localhost:4173,http://localhost:25342,https://baditaflorin.github.io" split_words:"true"`
	WorkDir               string        `default:"./audit-work" split_words:"true" validate:"required"`
	ToolTimeout           time.Duration `default:"90s" split_words:"true"`
	MaxUploadBytes        int64         `default:"1048576" split_words:"true"`
	MaxMaintainerPackages int           `default:"40" split_words:"true"`
	OllamaBaseURL         string        `split_words:"true"`
	OllamaModel           string        `default:"llama3.2" split_words:"true"`
	LLMCommand            string        `split_words:"true"`
	LogLevel              string        `default:"info" split_words:"true"`
}

func Load() (Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return cfg, err
	}

	if len(cfg.AllowedOrigins) == 1 && strings.Contains(cfg.AllowedOrigins[0], ",") {
		cfg.AllowedOrigins = strings.Split(cfg.AllowedOrigins[0], ",")
	}

	for i := range cfg.AllowedOrigins {
		cfg.AllowedOrigins[i] = strings.TrimSpace(cfg.AllowedOrigins[i])
	}

	if err := validator.New().Struct(cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
