package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/suryansh74/auth-package/internal/db"
	"github.com/suryansh74/auth-package/internal/dto"
	"github.com/suryansh74/auth-package/internal/services"
	"github.com/suryansh74/auth-package/internal/token"
)

type UserHandler struct {
	app *fiber.App
	// injecting service in handler
	srv        services.AuthService
	tokenMaker token.Maker
}

func NewUserHandler(app *fiber.App, db db.Auth, tokenMaker token.Maker) {
	userHandler := &UserHandler{
		app:        app,
		srv:        services.NewAuthenticator(db),
		tokenMaker: tokenMaker,
	}

	userHandler.SetupRoutes()
}

func (uh *UserHandler) SetupRoutes() {
	// public routes
	uh.app.Get("/check", uh.CheckHealth)
	uh.app.Post("/register", uh.Register)
	uh.app.Post("/login", uh.Login)

	// private routes
}

func (uh *UserHandler) CheckHealth(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"messgae": "working fine",
	})
}

func (uh *UserHandler) Register(ctx *fiber.Ctx) error {
	// get incoming req
	var req dto.UserRegisterRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	// call register func
	res, err := uh.srv.Register(req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}
	accessToken, err := uh.tokenMaker.CreateToken(res.UserID, req.Email, time.Minute*15)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": "unable to create token",
		})
	}
	res.AccessToken = accessToken

	return ctx.Status(fiber.StatusCreated).JSON(&res)
}

func (uh *UserHandler) Login(ctx *fiber.Ctx) error {
	// get incoming req
	var req dto.UserLoginRequest
	err := ctx.BodyParser(&req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	// call login func
	res, err := uh.srv.Login(req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}

	accessToken, err := uh.tokenMaker.CreateToken(res.UserID, req.Email, time.Minute*15)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": "unable to create token",
		})
	}
	res.AccessToken = accessToken
	return ctx.Status(fiber.StatusOK).JSON(&res)
}
