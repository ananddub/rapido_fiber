package captain_verify_handler

import (
	"strconv"

	activity_verfication "encore.app/background/temporal/captain/verification/activity"
	"encore.app/internal/connection"
	"encore.app/internal/features/captain/verify/dto"
	"encore.app/internal/features/captain/verify/repo"
	verify_service "encore.app/internal/features/captain/verify/service"
	"github.com/gofiber/fiber/v3"
)

type VerifyService struct {
	svc *verify_service.VerifyService
}

func initVerifyService() (*VerifyService, error) {
	conn, err := connection.InitConnection()
	if err != nil {
		return nil, err
	}

	verifyRepo := repo.NewVerifyRepo(conn.S3Client, conn.Temporal)
	verifySvc := verify_service.NewVerifyService(verifyRepo)

	return &VerifyService{svc: verifySvc}, nil
}

func getCaptainID(ctx fiber.Ctx) (int32, error) {
	user_id := ctx.Get("user_id")
	id, _ := strconv.Atoi(user_id)

	return int32(id), nil
}

func (s *VerifyService) UploadAadhaar(ctx fiber.Ctx, req *dto.UploadDocumentRequest) (*dto.UploadDocumentResponse, error) {
	captainID, err := getCaptainID(ctx)
	if err != nil {
		return nil, err
	}

	if err := s.svc.UploadDocument(ctx, captainID, "aadhaar", req.FileData, req.ContentType); err != nil {
		return nil, err
	}

	return &dto.UploadDocumentResponse{Message: "Aadhaar uploaded successfully"}, nil
}

func (s *VerifyService) UploadLicense(ctx fiber.Ctx, req *dto.UploadDocumentRequest) (*dto.UploadDocumentResponse, error) {

	captainID, err := getCaptainID(ctx)
	if err != nil {
		return nil, err
	}

	if err := s.svc.UploadDocument(ctx, captainID, "license", req.FileData, req.ContentType); err != nil {
		return nil, err
	}

	return &dto.UploadDocumentResponse{Message: "License uploaded successfully"}, nil
}

func (s *VerifyService) UploadVehicle(ctx fiber.Ctx, req *dto.UploadDocumentRequest) (*dto.UploadDocumentResponse, error) {

	captainID, err := getCaptainID(ctx)
	if err != nil {
		return nil, err
	}

	if err := s.svc.UploadDocument(ctx, captainID, "vehicle", req.FileData, req.ContentType); err != nil {
		return nil, err
	}

	return &dto.UploadDocumentResponse{Message: "Vehicle document uploaded successfully"}, nil
}

func (s *VerifyService) GetVerificationStatus(ctx fiber.Ctx) (*activity_verfication.BackgroundVerificationWorkflowInput, error) {

	captainID, err := getCaptainID(ctx)
	if err != nil {
		return nil, err
	}

	status, err := s.svc.GetVerificationStatus(ctx, captainID)
	if err != nil {
		return nil, err
	}

	return status, nil
}
