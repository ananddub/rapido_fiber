package middleware

import "github.com/gofiber/fiber/v3"

type HandlerVlidatorFuncWithParams[T any, F any] func(fiber.Ctx, *T) (*F, error)

func ValidateMiddlewareWithParams[T any, F any](handler HandlerVlidatorFuncWithParams[T, F]) fiber.Handler {
	return func(c fiber.Ctx) error {
		req := new(T)
		if err := c.Bind().All(req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		data, err := handler(c, req)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		if data == nil && err == nil {
			return nil
		}
		return c.JSON(data)
	}
}

type HandlerVlidatorFunc[T any] func(fiber.Ctx) (*T, error)

func ValidateMiddleware[T any](handler HandlerVlidatorFunc[T]) fiber.Handler {
	return func(c fiber.Ctx) error {
		data, err := handler(c)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		if data == nil && err == nil {
			return nil
		}
		return c.JSON(data)
	}
}
