package service

import (
	"context"
	"fmt"

	"github.com/roothash-pay/wallet-services/database"
	"github.com/roothash-pay/wallet-services/database/backend"
)

type NewsletterService interface {
	Create(ctx context.Context, req CreateNewsletterRequest) (*backend.Newsletter, error)
	Update(ctx context.Context, req UpdateNewsletterRequest) error
	GetByGuid(ctx context.Context, guid string) (*backend.Newsletter, error)
	List(ctx context.Context, req ListNewsletterRequest) ([]*backend.Newsletter, int64, error)
}

type CreateNewsletterRequest struct {
	Guid    string `json:"guid"`
	CatUUID string `json:"cat_uuid"`
	Title   string `json:"title"`
	Image   string `json:"image"`
	Detail  string `json:"detail"`
	LinkURL string `json:"link_url"`
}

type ListNewsletterRequest struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	CatUUID  string `json:"cat_uuid"`
	Title    string `json:"title"`
}

type UpdateNewsletterRequest struct {
	Guid    string                 `json:"guid"`
	Updates map[string]interface{} `json:"updates"`
}

type newsletterService struct {
	db *database.DB
}

func NewNewsletterService(db *database.DB) NewsletterService {
	return &newsletterService{db: db}
}

func (s *newsletterService) Create(
	ctx context.Context,
	req CreateNewsletterRequest,
) (*backend.Newsletter, error) {

	if req.Guid == "" {
		return nil, fmt.Errorf("guid required")
	}
	if req.CatUUID == "" || req.Title == "" || req.Image == "" || req.LinkURL == "" {
		return nil, fmt.Errorf("missing required fields")
	}

	item := &backend.Newsletter{
		Guid:    req.Guid,
		CatUUID: req.CatUUID,
		Title:   req.Title,
		Image:   req.Image,
		Detail:  req.Detail,
		LinkURL: req.LinkURL,
	}

	if err := s.db.BackendNewsletter.StoreNewsletter(item); err != nil {
		return nil, err
	}

	return item, nil
}

func (s *newsletterService) Update(
	ctx context.Context,
	req UpdateNewsletterRequest,
) error {

	if req.Guid == "" {
		return fmt.Errorf("guid required")
	}
	if len(req.Updates) == 0 {
		return fmt.Errorf("updates empty")
	}

	return s.db.BackendNewsletter.UpdateNewsletter(req.Guid, req.Updates)
}

func (s *newsletterService) GetByGuid(
	ctx context.Context,
	guid string,
) (*backend.Newsletter, error) {

	if guid == "" {
		return nil, fmt.Errorf("guid required")
	}
	return s.db.BackendNewsletter.GetByGuid(guid)
}

func (s *newsletterService) List(
	ctx context.Context,
	req ListNewsletterRequest,
) ([]*backend.Newsletter, int64, error) {

	filters := map[string]interface{}{}

	if req.CatUUID != "" {
		filters["cat_uuid"] = req.CatUUID
	}
	if req.Title != "" {
		filters["title"] = req.Title
	}

	return s.db.BackendNewsletter.GetNewsletterList(
		req.Page,
		req.PageSize,
		filters,
	)
}
