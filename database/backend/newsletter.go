// newsletter.go
package backend

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"gorm.io/gorm"
)

type Newsletter struct {
	Guid       string    `gorm:"primaryKey;column:guid;type:varchar" json:"guid"`
	CatUUID    string    `gorm:"column:cat_uuid;type:varchar(255);not null" json:"cat_uuid"`
	Title      string    `gorm:"column:title;type:varchar(255);not null" json:"title"`
	Image      string    `gorm:"column:image;type:varchar(700);not null" json:"image"`
	Detail     string    `gorm:"column:detail;type:text;default:''" json:"detail"`
	LinkURL    string    `gorm:"column:link_url;type:varchar(255);not null" json:"link_url"`
	CreateTime time.Time `gorm:"column:created_at;autoCreateTime" json:"create_time"`
	UpdateTime time.Time `gorm:"column:updated_at;autoUpdateTime" json:"update_time"`
}

func (Newsletter) TableName() string {
	return "newsletter"
}

type NewsletterView interface {
	GetByGuid(guid string) (*Newsletter, error)
	GetNewsletterList(page, pageSize int, filters map[string]interface{}) ([]*Newsletter, int64, error)
}

type NewsletterDB interface {
	NewsletterView

	StoreNewsletter(n *Newsletter) error
	StoreNewsletters(list []*Newsletter) error
	UpdateNewsletter(guid string, updates map[string]interface{}) error
}

type newsletterDB struct {
	gorm *gorm.DB
}

func NewNewsletterDB(db *gorm.DB) NewsletterDB {
	return &newsletterDB{gorm: db}
}

func (db *newsletterDB) StoreNewsletter(n *Newsletter) error {
	if err := db.gorm.Create(n).Error; err != nil {
		log.Error("StoreNewsletter error", "err", err)
		return err
	}
	return nil
}

func (db *newsletterDB) StoreNewsletters(list []*Newsletter) error {
	if err := db.gorm.CreateInBatches(list, len(list)).Error; err != nil {
		log.Error("StoreNewsletters error", "err", err)
		return err
	}
	return nil
}

func (db *newsletterDB) GetByGuid(guid string) (*Newsletter, error) {
	var n Newsletter
	if err := db.gorm.Where("guid = ?", guid).First(&n).Error; err != nil {
		log.Error("GetByGuid Newsletter error", "err", err)
		return nil, err
	}
	return &n, nil
}

func (db *newsletterDB) GetNewsletterList(page, pageSize int, filters map[string]interface{}) ([]*Newsletter, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	var list []*Newsletter
	query := db.gorm.Model(&Newsletter{})

	for key, value := range filters {
		if value == nil || value == "" {
			continue
		}
		switch key {
		case "title":
			query = query.Where("title LIKE ?", "%"+value.(string)+"%")
		default:
			query = query.Where(key+" = ?", value)
		}
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		log.Error("GetNewsletterList count error", "err", err)
		return nil, 0, err
	}

	if err := query.Order("created_at DESC").Limit(pageSize).Offset(offset).Find(&list).Error; err != nil {
		log.Error("GetNewsletterList list error", "err", err)
		return nil, 0, err
	}

	return list, total, nil
}

func (db *newsletterDB) UpdateNewsletter(guid string, updates map[string]interface{}) error {
	if guid == "" {
		return fmt.Errorf("invalid guid")
	}
	if len(updates) == 0 {
		return fmt.Errorf("updates is empty")
	}

	updates["updated_at"] = time.Now()

	if err := db.gorm.Model(&Newsletter{}).Where("guid = ?", guid).Updates(updates).Error; err != nil {
		log.Error("UpdateNewsletter error", "err", err)
		return err
	}
	return nil
}
