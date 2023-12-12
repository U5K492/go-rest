package main

import (
	"github.com/U5K492/go-rest/config"
	"github.com/U5K492/go-rest/framework/middleware"
	"github.com/U5K492/go-rest/framework/router"
	"github.com/U5K492/go-rest/handler"
	"log"
)

func main() {
	cfg, err := config.AppConfig()
	if err != nil {
		log.Fatalln(err)
	}
	engine := router.NewEngine()
	r := engine.Router

	r.Use(middleware.Logger)
	r.Use(middleware.TimeOut)
	r.Use(middleware.TimeCost)
	r.Get("/health", handler.HealthHandler())
	r.Get("/name/:name", handler.NameHandler())
	r.Post("/test", handler.RequestBodyHandler())

	engine.Run(uint64(cfg.Port))
}
