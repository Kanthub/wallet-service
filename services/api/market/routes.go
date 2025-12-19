package marketapi

import "github.com/go-chi/chi/v5"

func RegisterRoutes(r chi.Router, h *Handler) {
	r.Route("/market", func(r chi.Router) {
		r.Get("/price", h.GetPrice)
	})
}
