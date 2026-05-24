package main

import (
	"app/internal/app"
	"app/internal/store"
	"app/internal/ws"
	"database/sql"
	"fmt"
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

type Config struct {
	Port        int    `env:"SRV_PORT" envDefault:"8000"`
	PostgresDSN string `env:"POSTGRES_DSN,required"`
}

func loadConfig() *Config {
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	return cfg
}

func open(dsn string) (*bun.DB, error) {
	pgdb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(pgdb, pgdialect.New())

	// Log queries in dev
	db.WithQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	cfg := loadConfig()
	log.Println(cfg)

	db, err := open(cfg.PostgresDSN)
	if err != nil {
		log.Fatal(err)
	}

	msgStore := store.NewMessages(db)

	e := echo.New()
	e.Use(middleware.Recover(), middleware.RequestLogger())

	hub := ws.StartHub(msgStore)
	app.SetupRouter(e, hub, msgStore)

	addr := fmt.Sprintf(":%v", cfg.Port)

	if err := e.Start(addr); err != nil {
		e.Logger.Error("faild to start server", "error", err)
	}
}
