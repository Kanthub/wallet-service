package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum/log"
	"github.com/go-chi/chi/v5"

	"github.com/roothash-pay/wallet-services/services/api/service"
)

func (rs *Routes) NewsletterApi() {
	r := rs.router

	r.Route("/api/v1/newsletter", func(r chi.Router) {
		r.Post("/create", rs.createNewsletter)
		r.Post("/update", rs.updateNewsletter)
		r.Get("/info", rs.getNewsletter)
		r.Get("/list", rs.listNewsletter)
	})
}

func (rs *Routes) createNewsletter(w http.ResponseWriter, r *http.Request) {
	var req service.CreateNewsletterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	item, err := rs.svc.NewsletterService.Create(r.Context(), req)
	if err != nil {
		log.Error("create newsletter failed", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(item)
}

func (rs *Routes) updateNewsletter(w http.ResponseWriter, r *http.Request) {
	var req service.UpdateNewsletterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := rs.svc.NewsletterService.Update(r.Context(), req); err != nil {
		log.Error("update newsletter failed", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
	})
}

func (rs *Routes) getNewsletter(w http.ResponseWriter, r *http.Request) {
	guid := r.URL.Query().Get("guid")
	if guid == "" {
		http.Error(w, "guid required", http.StatusBadRequest)
		return
	}

	item, err := rs.svc.NewsletterService.GetByGuid(r.Context(), guid)
	if err != nil {
		log.Error("get newsletter failed", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(item)
}

func (rs *Routes) listNewsletter(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	req := service.ListNewsletterRequest{
		Page:     page,
		PageSize: pageSize,
		CatUUID:  r.URL.Query().Get("cat_uuid"),
		Title:    r.URL.Query().Get("title"),
	}

	list, total, err := rs.svc.NewsletterService.List(r.Context(), req)
	if err != nil {
		log.Error("list newsletter failed", "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"list":  list,
		"total": total,
	})
}
