package router

import (
	"net/http"

	"exchange/internal/api/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(h *handler.Handler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	// Orders route
	r.Post("/orders", h.SubmitOrder)
	r.Get("/orders", h.GetOrders)

	return r
}
