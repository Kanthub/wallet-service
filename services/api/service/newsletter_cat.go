package service

import (
	"context"
	"fmt"

	"github.com/roothash-pay/wallet-services/database"
	"github.com/roothash-pay/wallet-services/database/backend"
)

type NewsletterCatService interface {
	Create(ctx context.Context, req CreateNewsletterCatRequest) (*backend.NewsletterCat, error)
	Update(ctx context.Context, req UpdateNewsletterCatRequest) error
	GetByGuid(ctx context.Context, guid string) (*backend.NewsletterCat, error)
	ListAll(ctx context.Context) ([]*backend.NewsletterCat, error)
}

type CreateNewsletterCatRequest struct {
	Guid    string `json:"guid"`
	CatName string `json:"cat_name"`
}

type UpdateNewsletterCatRequest struct {
	Guid    string `json:"guid"`
	CatName string `json:"cat_name"`
}

type newsletterCatService struct {
	db *database.DB
}

func NewNewsletterCatService(db *database.DB) NewsletterCatService {
	return &newsletterCatService{db: db}
}

func (s *newsletterCatService) Create(
	ctx context.Context,
	req CreateNewsletterCatRequest,
) (*backend.NewsletterCat, error) {

	if req.Guid == "" {
		return nil, fmt.Errorf("guid required")
	}
	if req.CatName == "" {
		return nil, fmt.Errorf("cat_name required")
	}

	item := &backend.NewsletterCat{
		Guid:    req.Guid,
		CatName: req.CatName,
	}

	if err := s.db.BackendNewsletterCat.StoreNewsletterCat(item); err != nil {
		return nil, err
	}

	return item, nil
}

func (s *newsletterCatService) Update(
	ctx context.Context,
	req UpdateNewsletterCatRequest,
) error {

	if req.Guid == "" {
		return fmt.Errorf("guid required")
	}
	if req.CatName == "" {
		return fmt.Errorf("cat_name required")
	}

	updates := map[string]interface{}{
		"cat_name": req.CatName,
	}

	return s.db.BackendNewsletterCat.UpdateNewsletterCat(req.Guid, updates)
}

func (s *newsletterCatService) GetByGuid(
	ctx context.Context,
	guid string,
) (*backend.NewsletterCat, error) {

	if guid == "" {
		return nil, fmt.Errorf("guid required")
	}
	return s.db.BackendNewsletterCat.GetByGuid(guid)
}

func (s *newsletterCatService) ListAll(
	ctx context.Context,
) ([]*backend.NewsletterCat, error) {

	return s.db.BackendNewsletterCat.ListAll()
}
