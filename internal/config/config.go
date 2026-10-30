package config

import (
	"errors"
	"os"
)

type Config struct {
	DBURL    string
	HTTPAddr string
}

func Load() (Config, error) {
	config := Config{
		DBURL:    os.Getenv("DATABASE_URL"),
		HTTPAddr: os.Getenv("HTTP_ADDR"),
	}

	if config.DBURL == "" {
		return Config{}, errors.New("DATABASE_URL is required")
	}

	if config.HTTPAddr == "" {
		config.HTTPAddr = ":8080"
	}
	return config, nil
}
