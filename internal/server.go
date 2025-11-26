package internal

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suryansh74/auth-package/internal/db"
	"github.com/suryansh74/auth-package/internal/handlers"
	"github.com/suryansh74/auth-package/internal/token"
)

type Server struct {
	app         *fiber.App
	auth        db.Auth
	userHandler *handlers.UserHandler
	tokenMaker  token.Maker
}

const (
	symmetricKey = "GhR8pJHc2K3dN6mB4R7fj5G8Wol5hEHu"
)

// StartServer starts actual api server for auth it takes fiber app
func NewAuthServer(app *fiber.App, dbObj *pgxpool.Pool) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(symmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w ", err)
	}
	server := &Server{
		app:        app,
		auth:       db.NewAuth(dbObj),
		tokenMaker: tokenMaker,
	}

	// init user handler
	handlers.NewUserHandler(app, server.auth, server.tokenMaker)
	return server, nil
}
