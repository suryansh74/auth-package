package internal

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suryansh74/auth-package/internal/db"
	"github.com/suryansh74/auth-package/internal/handlers"
)

type Server struct {
	app         *fiber.App
	auth        db.Auth
	userHandler *handlers.UserHandler
}

// StartServer starts actual api server for auth it takes fiber app
func NewAuthServer(app *fiber.App, dbObj *pgxpool.Pool) (*Server, error) {
	server := &Server{
		app:  app,
		auth: *db.NewAuth(dbObj),
	}

	// init user handler
	handlers.NewUserHandler(app, &server.auth)
	return server, nil
}
