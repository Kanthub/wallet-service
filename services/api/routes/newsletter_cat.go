package routes

import (
	"encoding/json"
	"net/http"

	"github.com/ethereum/go-ethereum/log"
	"github.com/go-chi/chi/v5"

	"github.com/roothash-pay/wallet-services/services/api/service"
)

func (rs *Routes) NewsletterCatApi() {
	r := rs.router

	r.Route("/api/v1/newsletter-cat", func(r chi.Router) {
		r.Post("/create", rs.createNewsletterCat)
		r.Post("/update", rs.updateNewsletterCat)
		r.Get("/info", rs.getNewsletterCat)
		r.Get("/list", rs.listNewsletterCats)
	})
}

func (rs *Routes) createNewsletterCat(w http.ResponseWriter, r *http.Request) {
	var req service.CreateNewsletterCatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	item, err := rs.svc.NewsletterCatService.Create(r.Context(), req)
	if err != nil {
		log.Error("create newsletter_cat failed", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(item)
}

func (rs *Routes) updateNewsletterCat(w http.ResponseWriter, r *http.Request) {
	var req service.UpdateNewsletterCatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := rs.svc.NewsletterCatService.Update(r.Context(), req); err != nil {
		log.Error("update newsletter_cat failed", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "update success",
	})
}

func (rs *Routes) getNewsletterCat(w http.ResponseWriter, r *http.Request) {
	guid := r.URL.Query().Get("guid")
	if guid == "" {
		http.Error(w, "guid required", http.StatusBadRequest)
		return
	}

	item, err := rs.svc.NewsletterCatService.GetByGuid(r.Context(), guid)
	if err != nil {
		log.Error("get newsletter_cat failed", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(item)
}

func (rs *Routes) listNewsletterCats(w http.ResponseWriter, r *http.Request) {
	list, err := rs.svc.NewsletterCatService.ListAll(r.Context())
	if err != nil {
		log.Error("list newsletter_cat failed", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(list)
}
