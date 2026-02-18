package main

import (
	"encore.app/internal/app"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

type structValidator struct {
	validate *validator.Validate
}

func (v *structValidator) Validate(out any) error {
	return v.validate.Struct(out)
}

func main() {
	fapp := fiber.New(
		fiber.Config{
			StructValidator: &structValidator{validate: validator.New()},
		},
	)
	app.InitApp(fapp)
	fapp.Listen(":5000", fiber.ListenConfig{
		// EnablePrefork: true,
	})
}
