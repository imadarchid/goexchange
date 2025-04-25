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

func TestSubmitWrongOrder(t *testing.T) {
	ob := NewOrderBook()

	new_bad_order := order.NewOrder(-10, 100, types.Sell, types.Market)
	ob.Submit(new_bad_order)

	if ob.Asks.Len() > 0 {
		t.Errorf("Expected order not to be submitted to the orderbook.")
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
		t.Error("BUY order was transmitted as a SELL order.")
	}

	if ob.Asks.orders == nil {
		t.Error("BUY order was not transmitted to the orderbook.")
	}
}

func TestWithdrawOrder(t *testing.T) {
	ob := NewOrderBook()
	test_order := order.NewOrder(12, 100, types.Buy, types.Market)
	buy_order_to_withdraw := order.NewOrder(60, 100, types.Buy, types.Market)
	sell_order_to_withdraw := order.NewOrder(120, 100, types.Sell, types.Market)

	ob.Submit(test_order)
	ob.Submit(buy_order_to_withdraw)
	ob.Submit(sell_order_to_withdraw)

	if ob.Bids.Len()+ob.Asks.Len() != 3 {
		t.Error("Not all orders were submitted.")
	}

	ob.Withdraw(buy_order_to_withdraw)
	ob.Withdraw(sell_order_to_withdraw)

	if ob.Bids.Len() != 1 {
		t.Error("Withdrawing BUY order failed.")
	}

	if ob.Asks.Len() != 0 {
		t.Error("Withdrawing SELL order failed.")
	}
}

func TestWithdrawEmptyOrderBook(t *testing.T) {
	ob := NewOrderBook()
	order_to_withdraw := order.NewOrder(60, 100, types.Buy, types.Market)

	result := ob.Withdraw(order_to_withdraw)

	if result == true {
		t.Errorf("Expected an error, operation yielded a positive result.")
	}

}

func TestWithdrawBadOrder(t *testing.T) {
	ob := NewOrderBook()
	order_to_withdraw := order.NewOrder(60, 100, "SIDEWAYS", types.Market)

	result := ob.Withdraw(order_to_withdraw)

	if result == true {
		t.Errorf("Expected an error, operation yielded a positive result.")
	}

}

func TestOrdersMatched(t *testing.T) {
	ob := NewOrderBook()

	// Case 1: similar orders
	order_1 := order.NewOrder(60, 100, types.Buy, types.Market)
	order_2 := order.NewOrder(60, 100, types.Sell, types.Market)

	ob.Submit(order_1)
	ob.Submit(order_2)

	if ob.Asks.Len()-ob.Bids.Len() != 0 {
		t.Errorf("Expected an empty order book.")
	}

}

func TestOrdersPartiallyMatched(t *testing.T) {
	ob := NewOrderBook()

	order_1 := order.NewOrder(60, 100, types.Buy, types.Market)
	order_2 := order.NewOrder(60, 50, types.Sell, types.Market)

	ob.Submit(order_1)
	ob.Submit(order_2)

	topBid := ob.Bids.Peek()
	topAsk := ob.Asks.Peek()

	if topBid == nil {
		t.Errorf("Expected remaining bid orders left.")
	}

	if topAsk != nil {
		t.Errorf("Expected no remaining ask orders left.")
	}

	// if topBid.Status != types.PartiallyFilled {
	// 	t.Errorf("Expected status to be partially filled.")
	// }

}
