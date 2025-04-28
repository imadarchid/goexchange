package handler

import (
	"encoding/json"
	"exchange/internal/db"
	"exchange/internal/types"
	"fmt"
	"net/http"
)

type OrderRequest struct {
	Amount    int             `json:"amount"`
	Price     float64         `json:"price"`
	Side      types.Side      `json:"side"`
	OrderType types.OrderType `json:"order_type"`
	Ticker    string          `json:"ticker"`
}

type Handler struct {
	Queries *db.Queries
}

func (h *Handler) SubmitOrder(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) GetOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.Queries.GetAllOrders(r.Context())
	if err != nil {
		http.Error(w, "failed to list orders", http.StatusInternalServerError)
		fmt.Print(err)
		return
	}

	json.NewEncoder(w).Encode(orders)
}
