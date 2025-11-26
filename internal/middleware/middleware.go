package middleware

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/suryansh74/auth-package/token"
)

const (
	AuthorizationPayloadKey = "authorization_payload"
)

func AuthMiddleware(maker token.Maker) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authorization header",
			})
		}

		// Expecting: Bearer <token>
		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization format",
			})
		}

		// Extract token string
		accessToken := authHeader[7:]

		// Validate token
		payload, err := maker.VerifyToken(accessToken)
		if err != nil {
			if err == token.ErrExpiredToken {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error":   "Token expired",
					"details": "login again",
				})
			}

			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Incorrect token",
				"details": err.Error(),
			})
		}

		// Save payload for further handlers
		c.Locals(AuthorizationPayloadKey, payload)

		return c.Next()
	}
}

// GetAuthPayload retrieves the authenticated user's payload from context
func GetAuthPayload(c *fiber.Ctx) (*token.Payload, error) {
	value := c.Locals(AuthorizationPayloadKey)
	if value == nil {
		return nil, errors.New("authorization payload not found")
	}

	payload, ok := value.(*token.Payload)
	if !ok {
		return nil, errors.New("invalid payload type")
	}

	return payload, nil
}
