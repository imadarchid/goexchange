package order

import (
	"exchange/internal/db"
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID        uuid.UUID
	Price     float64
	Amount    int32
	Side      db.OrderSideType
	Time      time.Time
	Type      db.OrderType
	Status    db.OrderStatusType
	Ticker    string
	CreatedBy uuid.UUID
}

func NewOrder(price float64, amount int32, side db.OrderSideType, orderType db.OrderType, ticker string, createdBy uuid.UUID) *Order {
	return &Order{
		ID:        uuid.New(),
		Price:     price,
		Amount:    amount,
		Side:      side,
		Time:      time.Now().UTC(),
		Type:      orderType,
		Status:    "SUBMITTED",
		Ticker:    ticker,
		CreatedBy: createdBy,
	}
}

func (o *Order) IsValid() bool {
	if o.Price <= 0 || o.Amount <= 0 {
		return false
	}
	if o.Side != db.OrderSideTypeBUY && o.Side != db.OrderSideTypeSELL {
		return false
	}

	return true
}
