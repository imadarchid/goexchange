package orderbook

import (
	"exchange/internal/order"
	"exchange/internal/types"
	"testing"
)

func TestNewOrderBook(t *testing.T) {
	ob := NewOrderBook()

	if ob.Bids == nil {
		t.Error("Bids heap should not be nil")
	}
	if ob.Asks == nil {
		t.Error("Asks heap should not be nil")
	}

	if ob.Bids.Len() != 0 {
		t.Errorf("Expected Bids heap to be empty, got length %d", ob.Bids.Len())
	}
	if ob.Asks.Len() != 0 {
		t.Errorf("Expected Asks heap to be empty, got length %d", ob.Asks.Len())
	}
}

func TestSubmitSellOrder(t *testing.T) {
	ob := NewOrderBook()
	test_order := order.NewOrder(10, 100, types.Sell, types.Market)

	ob.Submit(test_order)
	if ob.Bids.Len() > 0 {
		t.Error("SELL order was transmitted as a BUY order.")
	}

	if ob.Asks.orders == nil {
		t.Error("SELL order was not transmitted to the orderbook.")
	}
}

func TestSubmitBuyOrder(t *testing.T) {
	ob := NewOrderBook()
	test_order := order.NewOrder(10, 100, types.Buy, types.Market)

	ob.Submit(test_order)
	if ob.Asks.Len() > 0 {
		t.Error("SELL order was transmitted as a BUY order.")
	}

	if ob.Asks.orders == nil {
		t.Error("SELL order was not transmitted to the orderbook.")
	}
}
