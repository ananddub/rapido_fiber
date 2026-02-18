package activity_verfication

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

type Status struct {
	IsUploaded bool `json:"is_uploaded"`
	IsPending  bool `json:"is_pending"`
	IsAccepted bool `json:"is_accepted"`
	IsRejected bool `json:"is_rejected"`
}

type BackgroundVerificationWorkflowInput struct {
	CaptainID       int32  `json:"captain_id"`
	AadhaarUploaded Status `json:"aadhaar_uploaded"`
	LicenseUploaded Status `json:"license_uploaded"`
	VehicleUploaded Status `json:"vehicle_uploaded"`

	AadhaarVerified    Status `json:"aadhaar_verified"`
	LicenseVerified    Status `json:"license_verified"`
	VehicleVerified    Status `json:"vehicle_verified"`
	BackgroundVerified Status `json:"background_verified"`

	Status string `json:"status"`
}

func InitBackgroundVerficationState(captainID int32, status string) BackgroundVerificationWorkflowInput {
	c := Status{
		IsUploaded: false,
		IsPending:  false,
		IsAccepted: false,
		IsRejected: false,
	}
	return BackgroundVerificationWorkflowInput{
		CaptainID:          captainID,
		AadhaarUploaded:    c,
		LicenseUploaded:    c,
		VehicleUploaded:    c,
		AadhaarVerified:    c,
		LicenseVerified:    c,
		VehicleVerified:    c,
		BackgroundVerified: c,
		Status:             status,
	}
}

func UploadCaptainDocumentsWorkflowAndWait(ctx workflow.Context, input *BackgroundVerificationWorkflowInput) (bool, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("UploadCaptainDocumentsWorkflow started", "captainID", input.CaptainID)

	sigUserAadhaar := workflow.GetSignalChannel(ctx, "user_upload_aadhaar")
	sigUserLicense := workflow.GetSignalChannel(ctx, "user_upload_license")
	sigUserVehicle := workflow.GetSignalChannel(ctx, "user_upload_vehicle")

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 1 * time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	uploadSelector := workflow.NewSelector(ctx)

	uploadSelector.AddReceive(sigUserAadhaar, func(c workflow.ReceiveChannel, more bool) {
		input.AadhaarUploaded.IsUploaded = true
		input.AadhaarUploaded.IsPending = true
		input.AadhaarUploaded.IsAccepted = false
		input.AadhaarUploaded.IsRejected = false
		logger.Info("Aadhaar uploaded by user")

		c.Receive(ctx, nil)
	})

	uploadSelector.AddReceive(sigUserLicense, func(c workflow.ReceiveChannel, more bool) {
		input.LicenseUploaded.IsUploaded = true
		input.LicenseUploaded.IsPending = true
		input.LicenseUploaded.IsAccepted = false
		input.LicenseUploaded.IsRejected = false
		logger.Info("License uploaded by user")

		c.Receive(ctx, nil)
	})

	uploadSelector.AddReceive(sigUserVehicle, func(c workflow.ReceiveChannel, more bool) {
		input.VehicleUploaded.IsUploaded = true
		input.VehicleUploaded.IsPending = true
		input.VehicleUploaded.IsAccepted = false
		input.VehicleUploaded.IsRejected = false
		logger.Info("Vehicle uploaded by user")

		c.Receive(ctx, nil)
	})

	workflow.Go(ctx, func(ctx workflow.Context) {
		for {
			uploadSelector.Select(ctx)
		}
	})
	workflow.Await(ctx, func() bool {
		return input.AadhaarUploaded.IsUploaded && input.LicenseUploaded.IsUploaded && input.VehicleUploaded.IsUploaded
	})
	input.Status = "PENDING_ADMIN_VERIFICATION"
	return true, nil
}
