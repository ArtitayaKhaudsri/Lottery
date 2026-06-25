package config

import (
	"github.com/joho/godotenv"
)

type Config struct {
	Type []string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		Type: []string{"FIRST", "SECOND", "THIRD", "FOURTH", "FOURTH", "FIFTH", "FIFTH", "SIXTH", "SIXTH", "FIFTH"},
	}
}
