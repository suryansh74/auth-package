package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/suryansh74/auth-package/internal/db"
	"github.com/suryansh74/auth-package/internal/dto"
	"github.com/suryansh74/auth-package/internal/middleware"
	"github.com/suryansh74/auth-package/internal/services"
	"github.com/suryansh74/auth-package/token"
)

type UserHandler interface {
	Register(ctx *fiber.Ctx) error
	Login(ctx *fiber.Ctx) error
	CheckAuthUser(ctx *fiber.Ctx) error
}

type userHandler struct {
	app *fiber.App
	// injecting service in handler
	srv                 services.AuthService
	tokenMaker          token.Maker
	accessTokenDuration time.Duration
}

func NewUserHandler(app *fiber.App, db db.Auth, tokenMaker token.Maker, accessTokenDuration time.Duration) UserHandler {
	return &userHandler{
		app:                 app,
		srv:                 services.NewAuthenticator(db),
		tokenMaker:          tokenMaker,
		accessTokenDuration: accessTokenDuration,
	}
}

func (uh *userHandler) Register(ctx *fiber.Ctx) error {
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
	accessToken, err := uh.tokenMaker.CreateToken(res.UserID, req.Email, time.Minute)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": "unable to create token",
		})
	}
	res.AccessToken = accessToken

	return ctx.Status(fiber.StatusCreated).JSON(&res)
}

func (uh *userHandler) Login(ctx *fiber.Ctx) error {
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

// CheckAuthUser verifies that the authentication middleware is working correctly.
//
// This endpoint is intended for testing purposes. If the request reaches this
// handler, it means the user has passed the authentication middleware and the
// token is valid.
//
// It returns a 200 OK response with a simple JSON message.
func (uh *userHandler) CheckAuthUser(ctx *fiber.Ctx) error {
	payload, err := middleware.GetAuthPayload(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}
	user, err := uh.srv.GetUserByID(ctx.Context(), payload.UserID)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(&fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(&dto.UserResponse{
		UserID:    user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
	})
}
