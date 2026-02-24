package app

import (
	"context"

	background "encore.app/internal/background/state"
	booking_common "encore.app/internal/common/booking/redis"
	"encore.app/internal/connection"
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
	conn, err := connection.InitConnection()
	if err != nil {
		return err
	}

	client := background.NewClient(conn.Mqtt, conn.Redis, conn.Query)
	ctx := context.Background()
	booking_common.NewRepo(conn.Redis).Init(ctx)
	go client.ListnerUser(ctx)
	return nil
}
