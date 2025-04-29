package router

import (
	"encoding/json"
	"net/http"

	"exchange/internal/api/handler"
	"exchange/internal/api/middlewares"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(h *handler.Handler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middlewares.JSONMiddleware)

	// Default route
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{
			"message": "it works",
		})
	})

	// Orders route
	r.Post("/orders", h.SubmitOrder)
	r.Get("/orders", h.GetOrders)

	// Assets route
	r.Get("/assets", h.GetAssets)

	return r
}
