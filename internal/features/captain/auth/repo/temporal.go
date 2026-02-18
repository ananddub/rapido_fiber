package repo

import (
	"context"
	"fmt"

	captain_verification "encore.app/background/temporal/captain/verification"
	"encore.app/internal/pkg/errs"

	"go.temporal.io/sdk/client"
)

func (r *CaptainRepo) StartVerificationWorkflow(ctx context.Context, captainID int32) (string, error) {
	workflowOptions := client.StartWorkflowOptions{
		ID:        fmt.Sprintf("captain-%d", captainID),
		TaskQueue: "captain-verification-task-queue",
	}

	we, err := r.temporal.ExecuteWorkflow(
		ctx,
		workflowOptions,
		captain_verification.BackgroundVerficationWorkflow,
		captainID,
	)
	if err != nil {
		return "", errs.Internal(err, "failed to start verification workflow")
	}
	fmt.Println("Verification workflow started with ID:", we.GetID())
	return we.GetID(), nil
}
