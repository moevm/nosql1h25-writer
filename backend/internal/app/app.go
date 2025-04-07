package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/sv-tools/mongoifc"

	"github.com/moevm/nosql1h25-writer/backend/config"
	"github.com/moevm/nosql1h25-writer/backend/internal/api"
	"github.com/moevm/nosql1h25-writer/backend/pkg/httpserver"
	"github.com/moevm/nosql1h25-writer/backend/pkg/mongodb"
)

type App struct {
	cfg       *config.Config
	interrupt <-chan os.Signal

	mongoClient mongoifc.Client

	echoHandler *echo.Echo

	mainDb mongoifc.Database

	ordersCollection mongoifc.Collection

	getHealthHandler api.Handler
	getUsersHandler  api.Handler
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

//	@title			Writer
//	@version		1.0.0
//	@description	API for freelancer's site
//	@host		localhost:80
//	@BasePath	/api
//	@securityDefinitions.apikey	JWT
//	@in							header
//	@name						Authorization
//	@description				JSON Web Token

func (app *App) Start() {
	log.Info("Connecting to mongo...")
	mongoClient, err := mongodb.New(app.cfg.Mongo.Uri)
	if err != nil {
		log.Fatalf("app - Start - mongodb.New: %v", err)
	}
	app.mongoClient = mongoClient

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), app.cfg.Mongo.ShutdownTimeout)
		defer cancel()

		if err := mongoClient.Disconnect(ctx); err != nil {
			log.Errorf("app - Start - mongoClient.Disconnect: %v", err)
		}
	}()

	log.Info("Start server...")
	httpServer := httpserver.New(app.EchoHandler(), httpserver.Port(app.cfg.HTTP.Port))
	httpServer.Start()

	defer func() {
		if err := httpServer.Shutdown(); err != nil {
			log.Errorf("app - Start - httpServer.Shutdown: %v", err)
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
