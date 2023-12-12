package config

import (
	"github.com/caarlos0/env/v6"
	dotenv "github.com/joho/godotenv"
	"log"
)

type App struct {
	AppName string `env:"APP_NAME"`
	Port    int    `env:"PORT"`
}

func AppConfig() (*App, error) {
	if err := dotenv.Load(".env"); err != nil {
		log.Fatalln(err)
	}
	cfg := &App{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
