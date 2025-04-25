package main

import (
	"exchange/internal/order"
	"exchange/internal/orderbook"
	"exchange/internal/types"
)

func main() {
	testOrderbook := orderbook.NewOrderBook()

	new_order := order.NewOrder(12, 100, types.Buy, types.Market)
	testOrderbook.Submit(new_order)
	testOrderbook.Submit(new_order)

	new_order = order.NewOrder(13, 100, types.Sell, types.Market)
	testOrderbook.Submit(new_order)

	new_order = order.NewOrder(14, 100, types.Buy, types.Market)
	testOrderbook.Submit(new_order)

	new_order = order.NewOrder(12, 50, types.Sell, types.Market)
	testOrderbook.Submit(new_order)
}
