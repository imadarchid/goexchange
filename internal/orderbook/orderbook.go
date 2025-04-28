package orderbook

import (
	"exchange/internal/db"
	"exchange/internal/events"
	"exchange/internal/order"
	"fmt"
	"time"
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

// Submit implements the OrderBookInterface
func (ob *OrderBook) Submit(o *order.Order) bool {
	if !o.IsValid() {
		fmt.Println("Order not valid. WRONG_ORDER")
		return false
	}

	if ob == nil {
		fmt.Println("Ticker does not exist. BAD_TICKER")
		return false
	}

	if o.Ticker != ob.Ticker {
		fmt.Println("Ticker does not match orderbook. WRONG_TICKER")
		return false
	}

	if !handleOrderType(o, ob) {
		fmt.Printf("%s Order was NOT filled %d @ %.2f ID: %s\n",
			o.Side, o.Amount, o.Price, o.ID)
		o.Status = db.OrderStatusTypeCANCELED
		return false
	}

	fmt.Print("NEW ORDER EVENT", o.ID, "\n")

	event := events.OrderEvent{
		Price:     o.Price,
		Amount:    o.Amount,
		OrderType: o.Type,
		Side:      o.Side,
		Ticker:    o.Ticker,
		CreatedBy: o.CreatedBy,
	}
	events.NewOrderChan <- event

	fmt.Printf("%s %s Order submitted %d @ %.2f ID: %s\n",
		o.Side, o.Type, o.Amount, o.Price, o.ID)

	ob.MatchOrders()
	return true
}

// Withdraw implements the OrderBookInterface
func (ob *OrderBook) Withdraw(o *order.Order) bool {
	if o.Side == db.OrderSideTypeBUY {
		if ob.Bids.Len() > 0 {
			return ob.removeFromHeap(ob.Bids, o.ID)
		} else {
			return false
		}
	} else if o.Side == db.OrderSideTypeSELL {
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

		if bid.Type == db.OrderTypeMARKET {
			bid.Price = ask.Price
		} else if ask.Type == db.OrderTypeMARKET {
			ask.Price = bid.Price
		}

		if bid.Price >= ask.Price {
			tradeAmount := bid.Amount
			if ask.Amount < bid.Amount {
				tradeAmount = ask.Amount
			}

			bid.Amount -= tradeAmount
			ask.Amount -= tradeAmount

			bid.Status = db.OrderStatusTypePARTIALLYFILLED
			ask.Status = db.OrderStatusTypePARTIALLYFILLED

			if bid.Amount == 0 {
				ob.Bids.Delete()
				bid.Status = db.OrderStatusTypeFILLED
			}
			if ask.Amount == 0 {
				ob.Asks.Delete()
				ask.Status = db.OrderStatusTypeFILLED
			}

			event := events.TransactionEvent{
				Price:       ask.Price,
				Amount:      tradeAmount,
				BuyerOrder:  bid.ID,
				SellerOrder: ask.ID,
				Asset:       bid.Ticker,
				Timestamp:   time.Now(),
			}
			events.MatchEventChan <- event

			// Log the trade
			fmt.Printf("Matched %d @ %.2f (Buy ID: %s, Sell ID: %s)\n",
				tradeAmount, ask.Price, bid.ID, ask.ID)

		} else {
			break
		}
	}
}

func handleOrderType(o *order.Order, ob *OrderBook) bool {
	switch o.Side {
	case db.OrderSideTypeBUY:
		switch o.Type {
		case db.OrderTypeMARKET:
			if ob.Asks.Len() > 0 {
				o.Price = ob.Asks.Peek().Price
				ob.Bids.Insert(o)
				return true
			}
			return false
		case db.OrderTypeLIMIT:
			ob.Bids.Insert(o)
			return true
		default:
			return false
		}
	case db.OrderSideTypeSELL:
		switch o.Type {
		case db.OrderTypeMARKET:
			if ob.Bids.Len() > 0 {
				o.Price = ob.Bids.Peek().Price
				ob.Asks.Insert(o)
				return true
			}
			return false
		case db.OrderTypeLIMIT:
			ob.Asks.Insert(o)
			return true
		default:
			return false
		}
	default:
		return false
	}
}
