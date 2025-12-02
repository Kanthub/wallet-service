// newsletter_cat.go
package backend

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"gorm.io/gorm"
)

type NewsletterCat struct {
	Guid       string    `gorm:"primaryKey;column:guid;type:varchar" json:"guid"`
	CatName    string    `gorm:"column:cat_name;type:varchar;not null" json:"cat_name"`
	CreateTime time.Time `gorm:"column:created_at;autoCreateTime" json:"create_time"`
	UpdateTime time.Time `gorm:"column:updated_at;autoUpdateTime" json:"update_time"`
}

func (NewsletterCat) TableName() string {
	return "newsletter_cat"
}

type NewsletterCatView interface {
	GetByGuid(guid string) (*NewsletterCat, error)
}

type NewsletterCatDB interface {
	NewsletterCatView

	StoreNewsletterCat(c *NewsletterCat) error
	StoreNewsletterCats(list []*NewsletterCat) error
	UpdateNewsletterCat(guid string, updates map[string]interface{}) error
}

type newsletterCatDB struct {
	gorm *gorm.DB
}

func NewNewsletterCatDB(db *gorm.DB) NewsletterCatDB {
	return &newsletterCatDB{gorm: db}
}

func (db *newsletterCatDB) StoreNewsletterCat(c *NewsletterCat) error {
	if err := db.gorm.Create(c).Error; err != nil {
		log.Error("StoreNewsletterCat error", "err", err)
		return err
	}
	return nil
}

func (db *newsletterCatDB) StoreNewsletterCats(list []*NewsletterCat) error {
	if err := db.gorm.CreateInBatches(list, len(list)).Error; err != nil {
		log.Error("StoreNewsletterCats error", "err", err)
		return err
	}
	return nil
}

func (db *newsletterCatDB) GetByGuid(guid string) (*NewsletterCat, error) {
	var c NewsletterCat
	if err := db.gorm.Where("guid = ?", guid).First(&c).Error; err != nil {
		log.Error("GetByGuid NewsletterCat error", "err", err)
		return nil, err
	}
	return &c, nil
}

func (db *newsletterCatDB) UpdateNewsletterCat(guid string, updates map[string]interface{}) error {
	if guid == "" {
		return fmt.Errorf("invalid guid")
	}
	if len(updates) == 0 {
		return fmt.Errorf("updates is empty")
	}

	updates["updated_at"] = time.Now()

	if err := db.gorm.Model(&NewsletterCat{}).Where("guid = ?", guid).Updates(updates).Error; err != nil {
		log.Error("UpdateNewsletterCat error", "err", err)
		return err
	}
	return nil
}
