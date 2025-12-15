package service

import (
	"fmt"

	"github.com/roothash-pay/wallet-services/services/api/models"
)

type AuthenticatorService interface {
	/*
	* ============== authenticator ==============
	 */
	GenerateTOTP(req models.GenerateTOTPRequest) (*models.GenerateTOTPResponse, error)
	VerifyTOTP(req models.VerifyTOTPRequest) (*models.VerifyTOTPResponse, error)
}

type authenticatorService struct {
	h *HandlerSvc
}

func NewAuthenticatorService(h *HandlerSvc) AuthenticatorService {
	return &authenticatorService{h: h}
}

func (au *authenticatorService) GenerateTOTP(req models.GenerateTOTPRequest) (*models.GenerateTOTPResponse, error) {
	if req.UserId == "" {
		return nil, fmt.Errorf("user id is required")
	}

	totpKey, err := au.h.authenticatorService.GenerateSecret(req.UserId)
	if err != nil {
		return nil, fmt.Errorf("failed to generate TOTP: %w", err)
	}

	currentCode, err := au.h.authenticatorService.GenerateCurrentCode(totpKey.Secret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate current code: %w", err)
	}

	remainingTime := au.h.authenticatorService.GetRemainingTime()

	return &models.GenerateTOTPResponse{
		Secret:        totpKey.Secret,
		QRCodeURL:     totpKey.QRCodeURL,
		QRCodeImage:   totpKey.QRCodeImage,
		CurrentCode:   currentCode,
		RemainingTime: remainingTime,
		Message:       "TOTP secret generated successfully",
	}, nil
}

func (au *authenticatorService) VerifyTOTP(req models.VerifyTOTPRequest) (*models.VerifyTOTPResponse, error) {
	if req.Secret == "" {
		return nil, fmt.Errorf("secret is required")
	}
	if req.Code == "" {
		return nil, fmt.Errorf("code is required")
	}

	valid, err := au.h.authenticatorService.VerifyCode(req.Secret, req.Code)
	if err != nil {
		return nil, fmt.Errorf("failed to verify TOTP: %w", err)
	}

	remainingTime := au.h.authenticatorService.GetRemainingTime()

	message := "TOTP code is valid"
	if !valid {
		message = "TOTP code is invalid or expired"
	}

	return &models.VerifyTOTPResponse{
		Valid:         valid,
		RemainingTime: remainingTime,
		Message:       message,
	}, nil
}
