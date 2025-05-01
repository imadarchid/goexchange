package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"exchange/internal/api/handler"
	"exchange/internal/api/router"
	"exchange/internal/db"
	"exchange/internal/order"
	"exchange/internal/orderbook"
)

func main() {
	// @TODO: add connection string to .env
	connStr := "postgresql://postgres:postgres@localhost:5432/testdb?sslmode=disable"

	database, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Could not connect to DB: %v", err)
	}
	queries := db.New(database)

	// Load available assets on database

	assets, err := queries.GetAllAssets(context.Background())
	if err != nil {
		log.Fatal("Could not load tradable assets:", err)
	}

	validTickers := make(map[string]struct{})
	for _, asset := range assets {
		if asset.IsTradable {
			validTickers[asset.Ticker] = struct{}{}
		}
	}

	// Build orderbooks based on available assets
	orderbooks := make(map[string]*orderbook.OrderBook)
	for ticker := range validTickers {
		orderbooks[ticker] = orderbook.NewOrderBook(ticker)
	}

	// Set up API handler and router
	h := &handler.Handler{Queries: queries, OrderBooks: orderbooks, ValidTickers: validTickers}
	r := router.NewRouter(h)

	// Get all submitted orders not filled and add them to the orderbook
	orders, err := queries.GetSubmittedOrders(context.Background())
	if err != nil {
		panic(err)
	}

	for _, o := range orders {
		persisted_order := order.NewOrder(o.Price, o.Amount, o.Side, o.OrderType, o.Asset, o.CreatedBy)
		persisted_order.ID = o.ID
		if o.Side == "BUY" {
			h.OrderBooks[o.Asset].Bids.Insert(persisted_order)
		} else {
			h.OrderBooks[o.Asset].Asks.Insert(persisted_order)
		}
	}

	fmt.Printf("Loaded %d orders in memory \n", len(orders))

	// TX processing workers
	go handler.StartTransactionPersistenceWorker(queries)

	http.ListenAndServe(":3000", r)
}
