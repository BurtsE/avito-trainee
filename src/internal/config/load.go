package config

import (
	"encoding/json"
	"os"

	"github.com/caarlos0/env/v11"
)

func InitConfig() (*Config, error) {
	cfg := &Config{}

	data, err := os.ReadFile("./configs/config.json")
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	err = env.Parse(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
