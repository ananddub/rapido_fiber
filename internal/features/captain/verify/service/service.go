package verify_service

import (
	"context"
	"encoding/base64"
	"fmt"

	activity_verfication "encore.app/background/temporal/captain/verification/activity"
	"encore.app/internal/pkg/errs"
)

type Repository interface {
	UploadToMinio(ctx context.Context, bucket, key string, data []byte, contentType string) error
	SignalWorkflow(ctx context.Context, workflowID, signalName string, data interface{}) error
	QueryWorkflowStatus(ctx context.Context, workflowID string) (*activity_verfication.BackgroundVerificationWorkflowInput, error)
	GenerateCaptainId(captainID int32) string
}

type VerifyService struct {
	repo Repository
}

func NewVerifyService(repo Repository) *VerifyService {
	return &VerifyService{repo: repo}
}

func (s *VerifyService) UploadDocument(ctx context.Context, captainID int32, docType, fileDataBase64, contentType string) error {
	fileData, err := base64.StdEncoding.DecodeString(fileDataBase64)
	if err != nil {
		return errs.BadRequest("invalid base64 file data")
	}

	bucket := "captain-docs"
	key := fmt.Sprintf("%d/%s", captainID, docType)

	if err := s.repo.UploadToMinio(ctx, bucket, key, fileData, contentType); err != nil {
		return err
	}

	signalName := fmt.Sprintf("user_upload_%s", docType)
	workflowID := s.repo.GenerateCaptainId(captainID)

	if err := s.repo.SignalWorkflow(ctx, workflowID, signalName, true); err != nil {
		return err
	}

	return nil
}

func (s *VerifyService) GetVerificationStatus(ctx context.Context, captainID int32) (*activity_verfication.BackgroundVerificationWorkflowInput, error) {
	workflowID := s.repo.GenerateCaptainId(captainID)
	status, err := s.repo.QueryWorkflowStatus(ctx, workflowID)
	if err != nil {
		return nil, err
	}
	return status, nil
}
