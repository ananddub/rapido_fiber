package captain_auth_handler

import (
	"encore.app/internal/middleware"
	"github.com/gofiber/fiber/v3"
)

func RegisterAuthRoutes(app fiber.Router) error {
	handler, err := initAuthService()

	if err != nil {
		return err
	}

	group := app.Group("/auth")
	group.Post("/login", middleware.ValidateMiddlewareWithParams(handler.Login))
	group.Post("/refresh", middleware.ValidateMiddlewareWithParams(handler.Refresh))
	group.Get("/verify", middleware.ValidateMiddlewareWithParams(handler.Verify))

	return nil
}
