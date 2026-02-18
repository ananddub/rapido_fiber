package app

import (
	captain_auth_handler "encore.app/internal/features/captain/auth/handler"
	captain_verify_handler "encore.app/internal/features/captain/verify/handler"
	"github.com/gofiber/fiber/v3"
)

func InitCaptainApp(app *fiber.App) error {
	captain := app.Group("/captain")
	err := captain_auth_handler.RegisterAuthRoutes(captain)
	if err != nil {
		return err
	}
	err = captain_verify_handler.RegisterVerifyRoutes(captain)
	if err != nil {
		return err
	}
	return nil
}
