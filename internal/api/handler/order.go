package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"exchange/internal/db"
	"exchange/internal/events"
	"exchange/internal/order"
	"exchange/internal/orderbook"
	"exchange/internal/types"

	"github.com/google/uuid"
)

type OrderRequest struct {
	Amount    int32           `json:"amount"`
	Price     float64         `json:"price"`
	Side      types.Side      `json:"side"`
	OrderType types.OrderType `json:"order_type"`
	Ticker    string          `json:"ticker"`
}

type ErrorResponse struct {
	Code    string `json:"error"`
	Message string `json:"message"`
}

type Handler struct {
	Queries      *db.Queries
	OrderBooks   map[string]*orderbook.OrderBook
	ValidTickers map[string]struct{}
}

func (h *Handler) SubmitOrder(w http.ResponseWriter, r *http.Request) {
	var req OrderRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Code:    "INVALID_REQUEST_BODY",
			Message: "Failed to decode request body",
		})
		return
	}

	s := "fca8a8c9-d7fe-4a69-b8dd-79ef94c52863"
	id, err := uuid.Parse(s)
	if err != nil {
		log.Fatal("Invalid UUID:", err)
	}

	order := order.NewOrder(req.Price, req.Amount, db.OrderSideType(req.Side), db.OrderType(req.OrderType), req.Ticker, id)
	if order.IsValid() {
		status := h.OrderBooks[req.Ticker].Submit(order)
		fmt.Print("NEW ORDER SUBMITTED", order.ID, "\n")
		if status {
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(order)
		} else {
			json.NewEncoder(w).Encode(ErrorResponse{
				Code:    "ORDERBOOK_SUBMISSION_FAILED",
				Message: "Failed to submit order to orderbook",
			})
		}

	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Code:    "ORDER_NOT_VALID",
			Message: "Order is not valid",
		})
	}
}

func (h *Handler) GetOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.Queries.GetAllOrders(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Code:    "ORDER_RETRIEVAL_FAILED",
			Message: "Failed to retrieve orders",
		})
		fmt.Print(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
}

func StartOrderPersistenceWorker(queries *db.Queries) {
	for event := range events.NewOrderChan {
		fmt.Print((event.CreatedBy))

		queries.CreateOrder(context.Background(), db.CreateOrderParams{
			Price:     event.Price,
			Amount:    event.Amount,
			Side:      (db.OrderSideType(event.Side)),
			OrderType: (db.OrderType(event.OrderType)),
			Asset:     event.Ticker,
			CreatedBy: event.CreatedBy,
		})
	}
}
