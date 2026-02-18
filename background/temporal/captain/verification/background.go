package captain_verification

import (
	activity_verfication "encore.app/background/temporal/captain/verification/activity"

	"go.temporal.io/sdk/workflow"
)

func BackgroundVerficationWorkflow(ctx workflow.Context, captainID int32) (bool, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("BackgroundVerficationWorkflow started", "captainID", captainID)
	status := activity_verfication.InitBackgroundVerficationState(captainID, "PENDING_UPLOAD_CAPTAIN")

	err := workflow.SetQueryHandler(ctx, "status", func() (activity_verfication.BackgroundVerificationWorkflowInput, error) {
		return status, nil
	})

	if err != nil {
		return false, err
	}

	activites, err := activity_verfication.Init()
	if err != nil {
		return false, err
	}

	_, err = activity_verfication.UploadCaptainDocumentsWorkflowAndWait(ctx, &status)
	if err != nil {
		return false, err
	}

	success, err := activity_verfication.VerifyAdminDocumentsWorkflow(ctx, activites, &status)
	if err != nil {
		return false, err
	}

	if success {
		status.Status = "VERIFIED"
		return true, nil
	}

	status.Status = "REJECTED"
	return true, nil
}
