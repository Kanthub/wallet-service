package service

import (
	"context"

	"github.com/ethereum/go-ethereum/rpc"

	"github.com/roothash-pay/wallet-services/config"
	"github.com/roothash-pay/wallet-services/database"
	"github.com/roothash-pay/wallet-services/database/backend"
	"github.com/roothash-pay/wallet-services/services/api/validator"
	"github.com/roothash-pay/wallet-services/services/common"
	"github.com/roothash-pay/wallet-services/services/market/cache"
)

type HandlerSvc struct {
	v   *validator.Validator
	db  *database.DB
	cfg *config.Config

	backendAdminDB       backend.AdminDB
	emailService         *common.EmailService
	smsService           *common.SMSService
	authenticatorService *common.AuthenticatorService
	verificationManager  *common.VerificationCodeManager
	siweVerifier         *common.SIWEVerifier
	kodoService          *common.KodoService
	s3Service            *common.S3Service
	jwtSecret            string

	AddressAssetService      AddressAssetService
	AdminService             AdminService
	AuthService              AuthService
	AuthenticatorService     AuthenticatorService
	RoleService              RoleService
	RoleAuthService          RoleAuthService
	SysLogService            SysLogService
	ChainService             ChainService
	TokenService             TokenService
	ChainTokenService        ChainTokenService
	WalletService            WalletService
	WalletAddressService     WalletAddressService
	WalletAssetService       WalletAssetService
	AssetAmountStatService   AssetAmountStatService
	WalletTxRecordService    WalletTxRecordService
	WalletAddressNoteService WalletAddressNoteService
	FiatCurrencyRateService  FiatCurrencyRateService
	MarketPriceService       MarketPriceService
	KlineService             KlineService
	NewsletterCatService     NewsletterCatService
	NewsletterService        NewsletterService
	WalletBalanceService     WalletBalanceService

	DappLinkService DappLinkService
	RpcService      RpcService
	Client          map[ChainType]*rpc.Client
}

func New(v *validator.Validator,
	db *database.DB,
	cfg *config.Config,
	backendUserDB backend.AdminDB,
	emailService *common.EmailService,
	smsService *common.SMSService,
	authenticatorService *common.AuthenticatorService,
	kodoService *common.KodoService,
	s3Service *common.S3Service,
	jwtSecret string,
	domain string,
	cache cache.Cache,
) *HandlerSvc {

	chains := make([]ChainType, 0, len(cfg.Chains))
	for _, c := range cfg.Chains {
		chains = append(chains, ChainType(c))
	}

	//dappLinkService, err := NewDappLinkService(cfg, chains...)
	//if err != nil {
	//	panic(err)
	//}

	clients := make(map[ChainType]*rpc.Client)

	for _, chain := range chains {
		rpcURL, err := cfg.RpcConfig.RPC(string(chain))
		if err != nil {
			panic(err)
		}

		client, err := rpc.DialContext(context.Background(), rpcURL)
		if err != nil {
			panic(err)
		}

		clients[chain] = client
	}

	return &HandlerSvc{
		v:                    v,
		db:                   db,
		cfg:                  cfg,
		backendAdminDB:       backendUserDB,
		emailService:         emailService,
		smsService:           smsService,
		authenticatorService: authenticatorService,
		verificationManager:  common.NewVerificationCodeManager(),
		kodoService:          kodoService,
		s3Service:            s3Service,
		siweVerifier:         common.NewSIWEVerifier(jwtSecret, domain),
		jwtSecret:            jwtSecret,

		AddressAssetService:      NewAddressAssetService(db),
		AdminService:             NewAdminService(db, common.NewSIWEVerifier(jwtSecret, domain)),
		AuthService:              NewAuthService(db),
		AuthenticatorService:     NewAuthenticatorService(&HandlerSvc{authenticatorService: authenticatorService}),
		RoleService:              NewRoleService(db),
		RoleAuthService:          NewRoleAuthService(db),
		SysLogService:            NewSysLogService(db),
		ChainService:             NewChainService(db),
		TokenService:             NewTokenService(db),
		ChainTokenService:        NewChainTokenService(db),
		WalletService:            NewWalletService(db),
		WalletAddressService:     NewWalletAddressService(db),
		WalletAssetService:       NewWalletAssetService(db),
		AssetAmountStatService:   NewAssetAmountStatService(db),
		WalletTxRecordService:    NewWalletTxRecordService(db),
		WalletAddressNoteService: NewWalletAddressNoteService(db),
		FiatCurrencyRateService:  NewFiatCurrencyRateService(db),
		MarketPriceService:       NewMarketPriceService(db, cache),
		KlineService:             NewKlineService(db),
		NewsletterCatService:     NewNewsletterCatService(db),
		NewsletterService:        NewNewsletterService(db),
		WalletBalanceService:     NewWalletBalanceService(db),
		//DappLinkService:          dappLinkService,
		RpcService: NewRpcService(cfg.RpcServer.RPCURL()),
		Client:     clients,
	}

}
