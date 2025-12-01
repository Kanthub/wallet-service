package service

import (
	"github.com/roothash-pay/wallet-services/database"
	"github.com/roothash-pay/wallet-services/database/backend"
	"github.com/roothash-pay/wallet-services/services/api/models"
	backend2 "github.com/roothash-pay/wallet-services/services/api/models/backend"
	"github.com/roothash-pay/wallet-services/services/api/validator"
	"github.com/roothash-pay/wallet-services/services/common"
)

type Service interface {
	/*
	* ============== authenticator ==============
	 */
	GenerateTOTP(req models.GenerateTOTPRequest) (*models.GenerateTOTPResponse, error)
	VerifyTOTP(req models.VerifyTOTPRequest) (*models.VerifyTOTPResponse, error)

	/*
	* ============== backend user ==============
	 */
	AdminUserLogin(req backend2.AdminLoginRequest) (*backend2.AdminLoginResponse, error)
	AdminLogout(req backend2.AdminLogoutRequest) (*backend2.AdminLogoutResponse, error)
}

type HandlerSvc struct {
	v                    *validator.Validator
	db                   *database.DB
	backendAdminDB       backend.AdminDB
	emailService         *common.EmailService
	smsService           *common.SMSService
	authenticatorService *common.AuthenticatorService
	verificationManager  *common.VerificationCodeManager
	siweVerifier         *common.SIWEVerifier
	kodoService          *common.KodoService
	s3Service            *common.S3Service
	jwtSecret            string
}

func New(v *validator.Validator,
	db *database.DB,
	backendUserDB backend.AdminDB,
	emailService *common.EmailService,
	smsService *common.SMSService,
	authenticatorService *common.AuthenticatorService,
	kodoService *common.KodoService,
	s3Service *common.S3Service,
	jwtSecret string,
	domain string,
) Service {
	return &HandlerSvc{
		v:                    v,
		db:                   db,
		backendAdminDB:       backendUserDB,
		emailService:         emailService,
		smsService:           smsService,
		authenticatorService: authenticatorService,
		verificationManager:  common.NewVerificationCodeManager(),
		kodoService:          kodoService,
		s3Service:            s3Service,
		siweVerifier:         common.NewSIWEVerifier(jwtSecret, domain),
		jwtSecret:            jwtSecret,
	}
}
