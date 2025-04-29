package events

import (
	"exchange/internal/db"
	"exchange/internal/order"
	"time"

	"github.com/google/uuid"
)

type OrderEvent struct {
	Price     float64
	Amount    int32
	Side      db.OrderSideType
	OrderType db.OrderType
	Ticker    string
	CreatedBy uuid.UUID
}

type TransactionEvent struct {
	Price       float64
	Amount      int32
	BuyerOrder  *order.Order
	SellerOrder *order.Order
	Asset       string
	Timestamp   time.Time
}

var TransactionEventChan = make(chan TransactionEvent, 10000)
var NewOrderChan = make(chan OrderEvent, 10000)
