package app

import (
	"github.com/gofiber/fiber/v3"
)

func InitApp(app *fiber.App) error {

	err := InitCaptainApp(app)
	if err != nil {
		return err
	}

	err = InitCaptainApp(app)
	if err != nil {
		return err
	}

	return nil
}
