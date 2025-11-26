package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suryansh74/auth-package"
)

func main() {
	config, err := auth.LoadConfig(".")
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
	server, err := auth.NewAuthServer(app, connPool, config)
	if err != nil {
		log.Println("failed to make server")
	}

	// Setup auth routes
	server.SetupRoutes()

	// Public routes
	app.Get("/hi", sayHello)

	// Create protected group for all your authenticated routes
	protected := server.ProtectedGroup("/api") // ✅ All routes under /api require auth
	protected.Get("/users", sayHello)

	app.Get("/hii", sayHello)

	app.Get("/hiii", server.AuthMiddleware(), sayHello)

	app.Listen(config.ServerAddress)
}

func sayHello(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"message": "Hello",
	})
}
