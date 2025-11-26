package internal

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suryansh74/auth-package/internal/db"
	"github.com/suryansh74/auth-package/internal/handlers"
	"github.com/suryansh74/auth-package/internal/token"
	"github.com/suryansh74/auth-package/internal/utils"
)

type Server struct {
	app         *fiber.App
	auth        db.Auth
	userHandler *handlers.UserHandler
	tokenMaker  token.Maker
}

// NewAuthServer starts actual api server for auth it takes fiber app
func NewAuthServer(app *fiber.App, dbObj *pgxpool.Pool, config utils.Config) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w ", err)
	}
	server := &Server{
		app:        app,
		auth:       db.NewAuth(dbObj),
		tokenMaker: tokenMaker,
	}

	// init user handler
	handlers.NewUserHandler(app, server.auth, server.tokenMaker, config.AccessTokenDuration)
	return server, nil
}
