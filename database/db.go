package database

import (
	"context"
	"fmt"
	"os"

	"path/filepath"

	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/roothash-pay/wallet-services/common/retry"
	"github.com/roothash-pay/wallet-services/config"
	"github.com/roothash-pay/wallet-services/database/backend"
)

type DB struct {
	gorm                     *gorm.DB
	BackendAdmin             backend.AdminDB
	BackendAuth              backend.AuthDB
	BackendRole              backend.RoleDB
	BackendRoleAuth          backend.RoleAuthDB
	BackendSysLog            backend.SysLogDB
	BackendAddressAsset      backend.AddressAssetDB
	BackendAssetAmountStat   backend.AssetAmountStatDB
	BackendChain             backend.ChainDB
	BackendChainToken        backend.ChainTokenDB
	BackendFiatCurrencyRate  backend.FiatCurrencyRateDB
	BackendKline             backend.KlineDB
	BackendMarketPrice       backend.MarketPriceDB
	BackendNewsletter        backend.NewsletterDB
	BackendNewsletterCat     backend.NewsletterCatDB
	BackendToken             backend.TokenDB
	BackendWallet            backend.WalletDB
	BackendWalletAddress     backend.WalletAddressDB
	BackendWalletAddressNote backend.WalletAddressNoteDB
	BackendWalletAsset       backend.WalletAssetDB
	BackendWalletTxRecord    backend.WalletTxRecordDB
	QueneTxDB                backend.QueueTxDB
}

func NewDB(ctx context.Context, dbConfig config.DBConfig) (*DB, error) {
	dsn := fmt.Sprintf("host=%s dbname=%s sslmode=disable", dbConfig.Host, dbConfig.Name)
	if dbConfig.Port != 0 {
		dsn += fmt.Sprintf(" port=%d", dbConfig.Port)
	}
	if dbConfig.User != "" {
		dsn += fmt.Sprintf(" user=%s", dbConfig.User)
	}
	if dbConfig.Password != "" {
		dsn += fmt.Sprintf(" password=%s", dbConfig.Password)
	}

	gormConfig := gorm.Config{
		SkipDefaultTransaction: true,
		CreateBatchSize:        3_000,
	}
	retryStrategy := &retry.ExponentialStrategy{Min: 1000, Max: 20_000, MaxJitter: 250}
	gorms, err := retry.Do[*gorm.DB](context.Background(), 10, retryStrategy, func() (*gorm.DB, error) {
		gorms, err := gorm.Open(postgres.Open(dsn), &gormConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to database: %w", err)
		}
		return gorms, nil
	})

	if err != nil {
		return nil, err
	}

	db := &DB{
		gorm:                     gorms,
		BackendAdmin:             backend.NewAdminDB(gorms),
		BackendAuth:              backend.NewAuthDB(gorms),
		BackendRole:              backend.NewRoleDB(gorms),
		BackendRoleAuth:          backend.NewRoleAuthDB(gorms),
		BackendSysLog:            backend.NewSysLogDB(gorms),
		BackendAddressAsset:      backend.NewAddressAssetDB(gorms),
		BackendAssetAmountStat:   backend.NewAssetAmountStatDB(gorms),
		BackendChain:             backend.NewChainDB(gorms),
		BackendChainToken:        backend.NewChainTokenDB(gorms),
		BackendFiatCurrencyRate:  backend.NewFiatCurrencyRateDB(gorms),
		BackendKline:             backend.NewKlineDB(gorms),
		BackendMarketPrice:       backend.NewMarketPriceDB(gorms),
		BackendNewsletter:        backend.NewNewsletterDB(gorms),
		BackendNewsletterCat:     backend.NewNewsletterCatDB(gorms),
		BackendToken:             backend.NewTokenDB(gorms),
		BackendWallet:            backend.NewWalletDB(gorms),
		BackendWalletAddress:     backend.NewWalletAddressDB(gorms),
		BackendWalletAddressNote: backend.NewWalletAddressNoteDB(gorms),
		BackendWalletAsset:       backend.NewWalletAssetDB(gorms),
		BackendWalletTxRecord:    backend.NewWalletTxRecordDB(gorms),
		QueneTxDB:                backend.NewQueueTxDB(gorms),
	}
	return db, nil
}

func (db *DB) Transaction(fn func(db *DB) error) error {
	return db.gorm.Transaction(func(tx *gorm.DB) error {
		txDB := &DB{
			gorm:                     tx,
			BackendAdmin:             backend.NewAdminDB(tx),
			BackendAuth:              backend.NewAuthDB(tx),
			BackendRole:              backend.NewRoleDB(tx),
			BackendRoleAuth:          backend.NewRoleAuthDB(tx),
			BackendSysLog:            backend.NewSysLogDB(tx),
			BackendAddressAsset:      backend.NewAddressAssetDB(tx),
			BackendAssetAmountStat:   backend.NewAssetAmountStatDB(tx),
			BackendChain:             backend.NewChainDB(tx),
			BackendChainToken:        backend.NewChainTokenDB(tx),
			BackendFiatCurrencyRate:  backend.NewFiatCurrencyRateDB(tx),
			BackendKline:             backend.NewKlineDB(tx),
			BackendMarketPrice:       backend.NewMarketPriceDB(tx),
			BackendNewsletter:        backend.NewNewsletterDB(tx),
			BackendNewsletterCat:     backend.NewNewsletterCatDB(tx),
			BackendToken:             backend.NewTokenDB(tx),
			BackendWallet:            backend.NewWalletDB(tx),
			BackendWalletAddress:     backend.NewWalletAddressDB(tx),
			BackendWalletAddressNote: backend.NewWalletAddressNoteDB(tx),
			BackendWalletAsset:       backend.NewWalletAssetDB(tx),
			BackendWalletTxRecord:    backend.NewWalletTxRecordDB(tx),
			QueneTxDB:                backend.NewQueueTxDB(tx),
		}
		return fn(txDB)
	})
}

func (db *DB) Exec(sql string, values ...interface{}) *gorm.DB {
	return db.gorm.Exec(sql, values...)
}

func (db *DB) Close() error {
	sql, err := db.gorm.DB()
	if err != nil {
		return err
	}
	return sql.Close()
}

func (db *DB) ExecuteSQLMigration(migrationsFolder string) error {
	err := filepath.Walk(migrationsFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Failed to process migration file: %s", path))
		}
		if info.IsDir() {
			return nil
		}
		fileContent, readErr := os.ReadFile(path)
		if readErr != nil {
			return errors.Wrap(readErr, fmt.Sprintf("Error reading SQL file: %s", path))
		}

		execErr := db.gorm.Exec(string(fileContent)).Error
		if execErr != nil {
			return errors.Wrap(execErr, fmt.Sprintf("Error executing SQL script: %s", path))
		}
		return nil
	})
	return err
}
