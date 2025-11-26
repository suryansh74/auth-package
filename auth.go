package auth

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/suryansh74/auth-package/internal/db"
	"github.com/suryansh74/auth-package/internal/handlers"
	"github.com/suryansh74/auth-package/internal/middleware"
	"github.com/suryansh74/auth-package/token"
)

type Server struct {
	app        *fiber.App
	auth       db.Auth
	tokenMaker token.Maker
	config     Config
}

func NewAuthServer(app *fiber.App, dbObj *pgxpool.Pool, config Config) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		app:        app,
		auth:       db.NewAuth(dbObj),
		tokenMaker: tokenMaker,
		config:     config,
	}
	return server, nil
}

// SetupRoutes registers authentication routes (register, login, check-auth-user)
//
// Public Routes:
//
//	POST /auth/register       → Register new user
//	POST /auth/login          → Login user
//
// Protected Routes:
//
//	GET  /auth/me             → Get current authenticated user info
func (s *Server) SetupRoutes() {
	userHandler := handlers.NewUserHandler(s.app, s.auth, s.tokenMaker, s.config.AccessTokenDuration)

	// Public auth routes at /auth prefix
	authGroup := s.app.Group("/auth")
	authGroup.Post("/register", userHandler.Register)
	authGroup.Post("/login", userHandler.Login)

	// Protected auth routes
	authGroup.Get("/me", s.AuthMiddleware(), userHandler.CheckAuthUser)
}

// AuthMiddleware returns the authentication middleware that can be used
// to protect custom routes in the application.
//
// Example usage:
//
//	server.SetupRoutes()
//	app.Get("/protected", server.AuthMiddleware(), myHandler)
func (s *Server) AuthMiddleware() fiber.Handler {
	return middleware.AuthMiddleware(s.tokenMaker)
}

// ProtectedGroup creates a new route group with authentication middleware applied.
// This is a convenience method for creating multiple protected routes.
//
// Example usage:
//
//	protected := server.ProtectedGroup("/api")
//	protected.Get("/users", getUsersHandler)
//	protected.Get("/posts", getPostsHandler)
func (s *Server) ProtectedGroup(prefix string) fiber.Router {
	return s.app.Group(prefix, s.AuthMiddleware())
}
