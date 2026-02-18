package captain_verification

import (
	"log"

	activity_verfication "encore.app/background/temporal/captain/verification/activity"
	"encore.app/internal/connection"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

//encore:service
type Service struct {
	client *client.Client
	worker *worker.Worker
}

func initService() (*Service, error) {
	conn, err := connection.InitConnection()
	c := conn.Temporal
	acts, err := activity_verfication.Init()
	w := worker.New(c, "captain-verification-task-queue", worker.Options{})
	w.RegisterWorkflow(BackgroundVerficationWorkflow)
	w.RegisterActivity(acts.VerifyCriminalRecord)
	w.RegisterActivity(acts.VerifyAdhar)
	w.RegisterActivity(acts.VerifyLicense)
	w.RegisterActivity(acts.VerifyVehicle)
	go func() {
		err = w.Run(worker.InterruptCh())
		if err != nil {
			log.Fatal("Unable to start worker", err)
		}
	}()
	return &Service{}, nil
}

func IinitService() (*Service, error) {
	return initService()
}
