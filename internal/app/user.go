package app

import (
	user_auth_handler "encore.app/internal/features/user/auth/handler"
	"github.com/gofiber/fiber/v3"
)

func RegisterUserRoutes(app *fiber.App) error {
	user := app.Group("/user")
	err := user_auth_handler.RegisterAuthRoutes(user)
	if err != nil {
		return err
	}
	return nil
}
