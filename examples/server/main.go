package main

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/timeout"
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
	defer connPool.Close()

	if err := connPool.Ping(context.Background()); err != nil {
		log.Fatalf("Unable to ping database: %v\n", err)
	}

	server, err := auth.NewAuthServer(app, connPool, config)
	if err != nil {
		log.Fatal("failed to make server:", err)
	}

	// Setup auth routes
	server.SetupRoutes()

	// Public route - no timeout
	app.Get("/hi", sayHello)

	// Route with timeout wrapper - handler MUST respect context
	app.Get("/hii", timeout.NewWithContext(sayHiWithTimeout, 3*time.Second))

	// Protected routes
	protected := server.ProtectedGroup("/api")
	protected.Get("/users", sayHello)

	log.Printf("Starting server on %s", config.ServerAddress)
	app.Listen(config.ServerAddress)
}

func sayHello(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"message": "Hello",
	})
}

// ✅ CORRECT: Handler that respects the timeout context
func sayHiWithTimeout(c *fiber.Ctx) error {
	// Use the helper function that checks context
	if err := sleepWithContext(c.UserContext(), 9*time.Second); err != nil {
		return err // Will return timeout error after 3 seconds
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"message": "Hello after sleep",
	})
}

// Helper function from Fiber docs that respects context timeout
func sleepWithContext(ctx context.Context, duration time.Duration) error {
	timer := time.NewTimer(duration)
	select {
	case <-ctx.Done():
		if !timer.Stop() {
			<-timer.C
		}
		return context.DeadlineExceeded // Returns timeout error
	case <-timer.C:
	}
	return nil
}
