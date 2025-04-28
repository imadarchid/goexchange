package handler

import (
	"context"
	"exchange/internal/db"
	"exchange/internal/events"
	"fmt"
)

func StartTransactionPersistenceWorker(queries *db.Queries) {
	for event := range events.MatchEventChan {
		fmt.Print(event.BuyerOrder)
		queries.CreateTransaction(context.Background(), db.CreateTransactionParams{
			BuyerOrder:  event.BuyerOrder,
			SellerOrder: event.SellerOrder,
			Price:       event.Price,
			Amount:      event.Amount,
			CreatedAt:   event.Timestamp,
		})
	}
}
