package captain_auth_handler

import (
	"fmt"

	"encore.app/internal/config"
	"encore.app/internal/connection"
	"encore.app/internal/features/captain/auth/dto"
	"encore.app/internal/features/captain/auth/repo"
	auth_service "encore.app/internal/features/captain/auth/service"
	"github.com/gofiber/fiber/v3"
)

type AuthService struct {
	svc *auth_service.AuthService
}

func initAuthService() (*AuthService, error) {

	cfg, err := config.InitConfig()
	if err != nil {
		return nil, err
	}

	conn, err := connection.InitConnection()
	if err != nil {
		return nil, err
	}

	captainRepo := repo.NewCaptainRepo(
		conn.Query,
		conn.Redis,
		conn.KafkaProducer,
		conn.Temporal,
	)

	authSvc := auth_service.NewAuthService(captainRepo, cfg)
	fmt.Println("Auth Service Initialized")
	return &AuthService{svc: authSvc}, nil
}

func (s *AuthService) Login(ctx fiber.Ctx, req *dto.LoginRequest) (*dto.LoginResponse, error) {

	if err := s.svc.GenerateOTP(ctx, req.Phone); err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Message: "OTP sent successfully",
	}, nil
}

func (s *AuthService) Verify(ctx fiber.Ctx, req *dto.VerifyRequest) (*dto.VerifyResponse, error) {

	accessToken, refreshToken, workflowID, err := s.svc.VerifyOTP(ctx, req.Phone, req.OTP)
	if err != nil {
		return nil, err
	}

	return &dto.VerifyResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		WorkflowID:   workflowID,
		Message:      "Verification successful",
	}, nil
}

func (s *AuthService) Refresh(ctx fiber.Ctx, req *dto.RefreshRequest) (*dto.RefreshResponse, error) {

	accessToken, refreshToken, err := s.svc.RefreshTokens(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &dto.RefreshResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Message:      "Token refreshed successfully",
	}, nil
}
