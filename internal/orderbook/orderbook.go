package orderbook

import (
	"exchange/internal/order"
	"exchange/internal/types"
	"fmt"
)

type OrderBook struct {
	Bids   *OrderHeap // max-heap (isMax = true)
	Asks   *OrderHeap // min-heap (isMax = false)
	Ticker string
}

func NewOrderBook(ticker string) *OrderBook {
	return &OrderBook{
		Bids:   NewOrderHeap(true),
		Asks:   NewOrderHeap(false),
		Ticker: ticker,
	}
}

func (ob *OrderBook) Submit(o *order.Order) bool {
	if !o.IsValid() {
		fmt.Println("Order not valid. WRONG_ORDER")
		return false
	}

	if o.Ticker != ob.Ticker {
		fmt.Println("Ticker does not match orderbook. WRONG_TICKER")
		return false
	}

	if !handleOrderType(o, ob) {
		fmt.Printf("%s Order was NOT filled %.2f @ %.2f ID: %s\n",
			o.Side, o.Amount, o.Price, o.ID)
		o.Status = types.Cancelled
		return false
	}

	fmt.Printf("%s %s Order submitted %.2f @ %.2f ID: %s\n",
		o.Side, o.Type, o.Amount, o.Price, o.ID)

	ob.MatchOrders()
	return true
}

func (ob *OrderBook) Withdraw(o *order.Order) bool {
	if o.Side == types.Buy {
		if ob.Bids.Len() > 0 {
			return ob.removeFromHeap(ob.Bids, o.ID)
		} else {
			return false
		}
	} else if o.Side == types.Sell {
		if ob.Asks.Len() > 0 {
			return ob.removeFromHeap(ob.Asks, o.ID)
		} else {
			return false
		}
	}
	return false
}

func (ob *OrderBook) MatchOrders() {
	for ob.Bids.Len() > 0 && ob.Asks.Len() > 0 {
		bid := ob.Bids.Peek()
		ask := ob.Asks.Peek()

		if bid.Type == types.Market {
			bid.Price = ask.Price
		} else if ask.Type == types.Market {
			ask.Price = bid.Price
		}

		if bid.Price >= ask.Price {
			tradeAmount := min(bid.Amount, ask.Amount)

			// Log the trade
			fmt.Printf("Matched %.2f @ %.2f (Buy ID: %s, Sell ID: %s)\n",
				tradeAmount, ask.Price, bid.ID, ask.ID)

			bid.Amount -= tradeAmount
			ask.Amount -= tradeAmount

			bid.Status = types.PartiallyFilled
			ask.Status = types.PartiallyFilled

			if bid.Amount == 0 {
				ob.Bids.Delete()
				bid.Status = types.Filled
			}
			if ask.Amount == 0 {
				ob.Asks.Delete()
				ask.Status = types.Filled
			}
		} else {
			break
		}
	}
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func handleOrderType(o *order.Order, ob *OrderBook) bool {
	switch o.Side {
	case types.Buy:
		switch o.Type {
		case types.Market:
			if ob.Asks.Len() > 0 {
				o.Price = ob.Asks.Peek().Price
				ob.Bids.Insert(o)
				return true
			}
			return false
		case types.Limit:
			ob.Bids.Insert(o)
			return true
		default:
			return false
		}
	case types.Sell:
		switch o.Type {
		case types.Market:
			if ob.Bids.Len() > 0 {
				o.Price = ob.Bids.Peek().Price
				ob.Asks.Insert(o)
				return true
			}
			return false
		case types.Limit:
			ob.Asks.Insert(o)
			return true
		default:
			return false
		}
	default:
		return false
	}
}
