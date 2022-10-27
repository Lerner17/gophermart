package config

import (
	"fmt"
	"log"

	"github.com/caarlos0/env"
)

type Config struct {
	ServerAddress        string `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	DatabaseDsn          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS" envDefault:"127.0.0.1:8081"`
}

var Instance *Config

func init() {
	log.Println("Load config...")
	log.Println("Successfully load config from env variables")
	Instance = new(Config)
	if err := env.Parse(Instance); err != nil {
		fmt.Printf("Cannot parse env vars %v\n", err)
	}
}
