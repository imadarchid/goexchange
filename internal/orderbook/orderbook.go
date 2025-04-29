package orderbook

import (
	"context"
	"exchange/internal/db"
	"exchange/internal/events"
	"exchange/internal/order"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// OrderCreator defines the interface required for creating an order persistence mechanism.
// This allows mocking the database interaction during tests.
type OrderCreator interface {
	CreateOrder(ctx context.Context, arg db.CreateOrderParams) (uuid.UUID, error)
}

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
func (ob *OrderBook) Submit(o *order.Order, creator OrderCreator) bool {
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

	// Use the creator interface to persist the order
	o_id, err := creator.CreateOrder(context.Background(), db.CreateOrderParams{
		Price:     o.Price,
		Amount:    o.Amount,
		Side:      o.Side,
		OrderType: o.Type,
		Asset:     o.Ticker,
		CreatedBy: o.CreatedBy,
	})

	// Handle potential error from CreateOrder
	if err != nil {
		fmt.Printf("Failed to create order in DB: %v\n", err)
		// Decide how to handle DB error - perhaps mark order as failed?
		// For now, we'll return false as the submission wasn't fully successful.
		return false
	}

	o.ID = o_id

	fmt.Printf("%s %s %s Order submitted %d @ %.2f ID: %s\n",
		o.Side, o.Ticker, o.Type, o.Amount, o.Price, o_id.String())

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

			fmt.Print(bid.Ticker)

			event := events.TransactionEvent{
				Price:       ask.Price,
				Amount:      tradeAmount,
				BuyerOrder:  bid,
				SellerOrder: ask,
				Asset:       bid.Ticker,
				Timestamp:   time.Now(),
			}
			events.TransactionEventChan <- event

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
