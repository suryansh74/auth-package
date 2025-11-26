package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suryansh74/auth-package/internal"
	"github.com/suryansh74/auth-package/internal/utils"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	app := fiber.New()
	// Create connection pool
	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}

	// Close closes the database connection pool
	defer connPool.Close()

	// Verify connection
	if err := connPool.Ping(context.Background()); err != nil {
		log.Fatalf("Unable to ping database: %v\n", err)
	}

	// for package it only server should start it should have config and database object
	// config will have secret key for paseto, address and port
	// and database object of pgx ONLY for now
	internal.NewAuthServer(app, connPool, config)

	app.Listen(config.ServerAddress)
}
