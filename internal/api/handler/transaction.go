package handler

import (
	"context"
	"exchange/internal/db"
	"exchange/internal/events"
	"fmt"
)

func StartTransactionPersistenceWorker(queries *db.Queries) {
	for event := range events.TransactionEventChan {

		_, err := queries.CreateTransaction(context.Background(), db.CreateTransactionParams{
			BuyerOrder:  event.BuyerOrder.ID,
			SellerOrder: event.SellerOrder.ID,
			Price:       event.Price,
			Amount:      event.Amount,
			Asset:       event.Asset,
			CreatedAt:   event.Timestamp,
		})

		if err == nil {
			fmt.Print("Buyer: ", event.BuyerOrder, "\n")
			queries.UpdateOrderStatus(context.Background(), db.UpdateOrderStatusParams{
				ID:          event.BuyerOrder.ID,
				OrderStatus: event.BuyerOrder.Status,
			})
			fmt.Print("Seller: ", event.SellerOrder, "\n")

			queries.UpdateOrderStatus(context.Background(), db.UpdateOrderStatusParams{
				ID:          event.SellerOrder.ID,
				OrderStatus: event.SellerOrder.Status,
			})
		} else {
			fmt.Print(err)
		}

	}
}
