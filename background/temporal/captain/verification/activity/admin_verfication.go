package activity_verfication

import "go.temporal.io/sdk/workflow"

func VerifyAdminDocumentsWorkflow(ctx workflow.Context, activities *ActivityVerfication, input *BackgroundVerificationWorkflowInput) (bool, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("VerifyAdminDocumentsWorkflow started", "captainID", input.CaptainID)
	sigAdminAadhaar := workflow.GetSignalChannel(ctx, "admin_verify_aadhaar")
	sigAdminLicense := workflow.GetSignalChannel(ctx, "admin_verify_license")
	sigAdminVehicle := workflow.GetSignalChannel(ctx, "admin_verify_vehicle")
	sigAdminBackground := workflow.GetSignalChannel(ctx, "admin_verify_background")

	selector := workflow.NewSelector(ctx)
	selector.AddReceive(sigAdminAadhaar, func(c workflow.ReceiveChannel, more bool) {
		var verified bool
		c.Receive(ctx, &verified)
		state := input.AadhaarUploaded
		if verified {
			state.IsAccepted = true
			state.IsRejected = false
			workflow.ExecuteActivity(ctx, activities.VerifyAdhar, input.CaptainID).Get(ctx, nil)
		} else {
			state.IsAccepted = false
			state.IsRejected = true
		}
		logger.Info("Aadhaar verified by admin")
	})
	selector.AddReceive(sigAdminLicense, func(c workflow.ReceiveChannel, more bool) {
		var verified bool
		c.Receive(ctx, &verified)
		state := input.LicenseUploaded
		if verified {
			state.IsAccepted = true
			state.IsRejected = false
			workflow.ExecuteActivity(ctx, activities.VerifyLicense, input.CaptainID).Get(ctx, nil)
		} else {
			state.IsAccepted = false
			state.IsRejected = true
		}
	})
	selector.AddReceive(sigAdminVehicle, func(c workflow.ReceiveChannel, more bool) {
		var verified bool
		c.Receive(ctx, &verified)
		state := input.VehicleUploaded
		if verified {
			state.IsAccepted = true
			state.IsRejected = false
			workflow.ExecuteActivity(ctx, activities.VerifyVehicle, input.CaptainID).Get(ctx, nil)
		} else {
			state.IsAccepted = false
			state.IsRejected = true
		}
	})
	selector.AddReceive(sigAdminBackground, func(c workflow.ReceiveChannel, more bool) {
		var verified bool
		c.Receive(ctx, &verified)
		state := input.BackgroundVerified
		if verified {
			state.IsAccepted = true
			state.IsRejected = false
			workflow.ExecuteActivity(ctx, activities.VerifyCriminalRecord, input.CaptainID).Get(ctx, nil)
		} else {
			state.IsAccepted = false
			state.IsRejected = true
		}
	})
	workflow.Go(ctx, func(ctx workflow.Context) {
		for {
			selector.Select(ctx)
		}
	})
	current_status := false
	workflow.Await(ctx, func() bool {
		if input.AadhaarVerified.IsRejected || input.LicenseVerified.IsRejected || input.VehicleVerified.IsRejected || input.BackgroundVerified.IsRejected {
			current_status = false
			return true
		}
		current_status = true
		return input.AadhaarVerified.IsAccepted && input.LicenseVerified.IsAccepted && input.VehicleVerified.IsAccepted && input.BackgroundVerified.IsAccepted
	})
	return current_status, nil
}
