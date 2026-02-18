package repo

import (
	"bytes"
	"context"
	"fmt"

	captain_verification "encore.app/background/temporal/captain/verification"
	activity_verfication "encore.app/background/temporal/captain/verification/activity"
	"encore.app/internal/pkg/errs"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.temporal.io/sdk/client"
)

type VerifyRepo struct {
	s3       *s3.Client
	temporal client.Client
}

func NewVerifyRepo(s3 *s3.Client, temporal client.Client) *VerifyRepo {
	return &VerifyRepo{
		s3:       s3,
		temporal: temporal,
	}
}

func (r *VerifyRepo) GenerateCaptainId(captainID int32) string {
	return fmt.Sprintf("captain-%d", captainID)
}

func (r *VerifyRepo) UploadToMinio(ctx context.Context, bucket, key string, data []byte, contentType string) error {
	_, err := r.s3.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return errs.Internal(err, "failed to upload file to minio")
	}
	return nil
}

func (r *VerifyRepo) SignalWorkflow(ctx context.Context, workflowID, signalName string, data interface{}) error {
	err := r.temporal.SignalWorkflow(ctx, workflowID, "", signalName, data)
	if err != nil {
		return errs.Internal(err, fmt.Sprintf("failed to signal workflow: %s", signalName))
	}
	return nil
}

func (r *VerifyRepo) QueryWorkflowStatus(ctx context.Context, workflowID string) (*activity_verfication.BackgroundVerificationWorkflowInput, error) {
	resp, err := r.temporal.QueryWorkflow(ctx, workflowID, "", "status", true)
	if err != nil {
		return nil, errs.Internal(err, "failed to query workflow status")
	}

	var status activity_verfication.BackgroundVerificationWorkflowInput
	if err := resp.Get(&status); err != nil {
		return nil, errs.Internal(err, "failed to decode workflow status")
	}

	return &status, nil
}

func (r *VerifyRepo) StartVerificationWorkflow(ctx context.Context, captainID int32) (string, error) {
	workflowOptions := client.StartWorkflowOptions{
		ID:        r.GenerateCaptainId(captainID),
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

	return we.GetID(), nil
}
