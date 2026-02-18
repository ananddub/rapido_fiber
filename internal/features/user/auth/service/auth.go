package auth_service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"encore.app/internal/features/captain/auth/repo"
	"encore.app/internal/features/captain/auth/util"
	"encore.app/internal/pkg/errs"
)

func (s *AuthService) GenerateOTP(ctx context.Context, phone string) error {
	otp := util.GenerateOTP()

	if err := s.repo.StoreOTP(ctx, phone, otp, 5*time.Minute); err != nil {
		return err
	}

	message := fmt.Sprintf("%s", otp)
	if err := s.repo.SendSMS(ctx, phone, message); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) VerifyOTP(ctx context.Context, phone, otp string) (accessToken string, refreshToken string, err error) {
	storedOTP, err := s.repo.GetOTP(ctx, phone)
	if err != nil {
		return "", "", err
	}

	if storedOTP != otp {
		return "", "", errs.BadRequest("invalid OTP")
	}

	_ = s.repo.DeleteOTP(ctx, phone)

	captainID, err := s.findOrCreateUser(ctx, phone)
	if err != nil {
		return "", "", err
	}

	accessToken, err = util.GenerateAccessToken(captainID, phone, s.cfg)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = util.GenerateRefreshToken(captainID, phone, s.cfg)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) RefreshTokens(ctx context.Context, refreshTokenStr string) (string, string, error) {
	captainID, phone, err := util.VerifyRefreshToken(refreshTokenStr, s.cfg)
	if err != nil {
		return "", "", err
	}

	newAccessToken, err := util.GenerateAccessToken(captainID, phone, s.cfg)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := util.GenerateRefreshToken(captainID, phone, s.cfg)
	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}

func (s *AuthService) findOrCreateUser(ctx context.Context, phone string) (int32, error) {
	captain, err := s.repo.GetByPhone(ctx, phone)
	if err == nil {
		return captain.ID, nil
	}

	if !errors.Is(err, repo.ErrNotFound) {
		return 0, errs.Internal(err, "failed to check captain existence")
	}

	newCaptain, err := s.repo.Create(ctx, "", phone)
	if err != nil {
		return 0, err
	}

	return newCaptain.ID, nil
}
