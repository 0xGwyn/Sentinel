package middleware

import (
	"crypto/sha256"
	"crypto/subtle"

	"github.com/0xgwyn/sentinel/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
)

func ValidateAPIKey(c *fiber.Ctx, key string) (bool, error) {
	apiKey, err := config.LoadEnv("API_KEY")
	if err != nil {
		return false, err
	}

	hashedAPIKey := sha256.Sum256([]byte(apiKey))
	hashedKey := sha256.Sum256([]byte(key))

	if subtle.ConstantTimeCompare(hashedAPIKey[:], hashedKey[:]) == 1 {
		return true, nil
	}
	return false, keyauth.ErrMissingOrMalformedAPIKey
}

func NewAuthMiddleware() fiber.Handler {
	return keyauth.New(keyauth.Config{
		KeyLookup:    "header:X-API-Key",
		Validator:    ValidateAPIKey,
		ErrorHandler: handleAuthError,
	})
}

func handleAuthError(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error":   "Unauthorized",
		"message": "Invalid or missing API key",
	})
}
