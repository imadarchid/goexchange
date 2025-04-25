package orderbook

import (
	"exchange/internal/order"
	"exchange/internal/types"
	"fmt"
)

type OrderBook struct {
	Bids *OrderHeap // max-heap (isMax = true)
	Asks *OrderHeap // min-heap (isMax = false)
}

func NewOrderBook() *OrderBook {
	return &OrderBook{
		Bids: NewOrderHeap(true),
		Asks: NewOrderHeap(false),
	}
}

func (ob *OrderBook) Submit(o *order.Order) {
	if !o.IsValid() {
		return
	}
	if o.Side == types.Buy {
		ob.Bids.Insert(o)
	} else if o.Side == types.Sell {
		ob.Asks.Insert(o)
	}

	fmt.Printf("%s Order submitted %.2f @ %.2f ID: %s\n",
		o.Side, o.Amount, o.Price, o.ID)

	ob.MatchOrders()
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
