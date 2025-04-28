package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"exchange/internal/db"
	"exchange/internal/order"
	"exchange/internal/types"
)

type OrderRequest struct {
	Amount    int             `json:"amount"`
	Price     float64         `json:"price"`
	Side      types.Side      `json:"side"`
	OrderType types.OrderType `json:"order_type"`
	Ticker    string          `json:"ticker"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type Handler struct {
	Queries *db.Queries
}

func (h *Handler) SubmitOrder(w http.ResponseWriter, r *http.Request) {
	var req OrderRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "Invalid Request",
			Message: "Failed to decode request body",
		})
		return
	}

	order := order.NewOrder(req.Price, float64(req.Amount), req.Side, req.OrderType, req.Ticker)
	if order.IsValid() {
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(order)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "Invalid Order",
			Message: "Failed to submit order",
		})
	}
}

func (h *Handler) GetOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.Queries.GetAllOrders(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error:   "Internal Server Error",
			Message: "Failed to retrieve orders",
		})
		fmt.Print(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
}
