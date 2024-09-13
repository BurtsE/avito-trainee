package app

import (
	"context"
	"log"
)

type App struct {
	serviceProvider *serviceProvider
}

func NewApp(ctx context.Context) *App {
	a := &App{}
	a.serviceProvider = NewSericeProvider()
	return a
}

func (a *App) Run() error {
	return a.runServer()
}

func (a *App) runServer() error {
	log.Printf("server is running on %s", a.serviceProvider.Config().Service.Host)
	return a.serviceProvider.router.Start()
}