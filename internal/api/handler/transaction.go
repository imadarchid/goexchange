package handler

import (
	"context"
	"exchange/internal/db"
	"exchange/internal/events"
	"fmt"
)

func StartTransactionPersistenceWorker(queries *db.Queries) {
	for event := range events.TransactionEventChan {
		fmt.Print("Buyer: ", event.BuyerOrder, "\n")
		fmt.Print("Seller: ", event.SellerOrder, "\n")

		queries.CreateTransaction(context.Background(), db.CreateTransactionParams{
			BuyerOrder:  event.BuyerOrder,
			SellerOrder: event.SellerOrder,
			Price:       event.Price,
			Amount:      event.Amount,
			Asset:       event.Asset,
			CreatedAt:   event.Timestamp,
		})

	}
}
