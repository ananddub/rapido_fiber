package captain_verify_handler

import (
	"context"
	"fmt"
	"testing"
	"time"

	captain_verification "encore.app/background/temporal/captain/verification"
	activity_verfication "encore.app/background/temporal/captain/verification/activity"
	"encore.app/internal/connection"
	"github.com/stretchr/testify/assert"
)

func TestTemporalStatus(t *testing.T) {
	conn, err := connection.InitConnection()
	if err != nil {
		t.Fatal(err)
	}
	resp, err := conn.Temporal.QueryWorkflow(
		context.Background(),
		"captain-1",
		"",
		"status",
	)
	if err != nil {
		t.Fatal(err)
	}

	var result activity_verfication.BackgroundVerificationWorkflowInput
	err = resp.Get(&result)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(result)

}
func TestTemporalStatusWithHandler(t *testing.T) {
	// verfiyhandler, err := initVerifyService()
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// ctx := auth.WithContext(context.Background(), auth.UID("1"), &auth_middleware.AuthData{CaptainID: 1, Phone: "1234567890"})
	// // handler, err := verfiyhandler.GetVerificationStatus(ctx)

	// fmt.Println(handler)
	// assert.NoError(t, err)

}

func TestSendSignal(t *testing.T) {
	conn, err := connection.InitConnection()
	go captain_verification.IinitService()

	captain_name := "captain-1"

	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	err = conn.Temporal.SignalWorkflow(
		ctx,
		captain_name,
		"",
		"user_upload_aadhaar",
		true,
	)
	assert.NoError(t, err)
	err = conn.Temporal.SignalWorkflow(
		ctx,
		captain_name,
		"",
		"user_upload_license",
		true,
	)
	assert.NoError(t, err)
	err = conn.Temporal.SignalWorkflow(
		ctx,
		captain_name,
		"",
		"user_upload_vehicle",
		true,
	)
	assert.NoError(t, err)
	err = conn.Temporal.SignalWorkflow(
		ctx,
		captain_name,
		"",
		"user_upload_background",
		true,
	)
	assert.NoError(t, err)
	TestAdminSendSignal(t)
}

func TestAdminSendSignal(t *testing.T) {
	conn, err := connection.InitConnection()
	captain_name := "captain-1"
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	err = conn.Temporal.SignalWorkflow(
		ctx,
		captain_name,
		"",
		"admin_verify_aadhar",
		true,
	)
	assert.NoError(t, err)
	err = conn.Temporal.SignalWorkflow(
		ctx,
		captain_name,
		"",
		"admin_verify_license",
		true,
	)
	assert.NoError(t, err)
	err = conn.Temporal.SignalWorkflow(
		ctx,
		captain_name,
		"",
		"admin_verify_vehicle",
		true,
	)
	assert.NoError(t, err)
	err = conn.Temporal.SignalWorkflow(
		ctx,
		captain_name,
		"",
		"admin_verify_background",
		true,
	)
	assert.NoError(t, err)
	time.Sleep(10 * time.Minute)
}
