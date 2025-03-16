package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/moevm/nosql1h25-writer/backend/config"
	"github.com/moevm/nosql1h25-writer/backend/pkg/httpserver"
	"github.com/moevm/nosql1h25-writer/backend/pkg/mongodb"
	log "github.com/sirupsen/logrus"
)

type App struct {
	cfg       *config.Config
	interrupt <-chan os.Signal

	echoHandler *echo.Echo
}

func New(configPath string) *App {
	cfg, err := config.New(configPath)
	if err != nil {
		log.Fatalf("app - New - config.New: %v", err)
	}

	initLogger(cfg.Log.Level)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	return &App{
		cfg:       cfg,
		interrupt: interrupt,
	}
}

func (app *App) Start() {
	log.Info("Connecting to mongo...")
	mongoClient, err := mongodb.New(app.cfg.Mongo.Uri)
	if err != nil {
		log.Fatalf("app - Start - mongodb.New: %v", err)
	}

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(app.cfg.Mongo.ShutdownTimeout))
		defer cancel()

		if err := mongoClient.Disconnect(ctx); err != nil {
			log.Errorf("app - Start - mongoClient.Disconnect: %v", err)
		}
	}()

	
	log.Info("Start server...")
	httpServer := httpserver.New(app.getEchoHandler(), httpserver.Port(app.cfg.HTTP.Port))
	httpServer.Start()

	defer func() {
		if err := httpServer.Shutdown(); err != nil {
			log.Errorf("app - Run - httpServer.Shutdown: %v", err)
		}
	}()

	select {
	case s := <-app.interrupt:
		log.Infof("app - Start - signal: %v", s)
	case err := <-httpServer.Notify():
		log.Errorf("app - Start - server error: %v", err)
	}

	log.Info("Shutting down...")
}
