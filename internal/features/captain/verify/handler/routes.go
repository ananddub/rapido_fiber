package captain_verify_handler

import (
	"encore.app/internal/middleware"
	"github.com/gofiber/fiber/v3"
)

func RegisterVerifyRoutes(app fiber.Router) error {
	handler, err := initVerifyService()
	if err != nil {
		return err
	}

	group := app.Group("/verify")

	group.Post("/upload/aadhaar", middleware.ValidateMiddlewareWithParams(handler.UploadAadhaar))
	group.Post("/upload/license", middleware.ValidateMiddlewareWithParams(handler.UploadLicense))
	group.Post("/upload/vehicle", middleware.ValidateMiddlewareWithParams(handler.UploadVehicle))
	group.Get("/status", middleware.ValidateMiddleware(handler.GetVerificationStatus))

	return nil
}
