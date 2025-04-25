package order

import (
	"exchange/internal/types"
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID     string
	Price  float64
	Amount float64
	Side   types.Side
	Time   time.Time
	Type   types.OrderType
	Status types.Status
}

// Order Actions
func NewOrder(price float64, amount float64, side types.Side, orderType types.OrderType) *Order {
	return &Order{
		ID:     uuid.New().String(),
		Price:  price,
		Amount: amount,
		Side:   side,
		Time:   time.Now().UTC(),
		Type:   orderType,
		Status: types.Pending,
	}
}

// Order Checks

func (o *Order) IsValid() bool {
	if o.Price <= 0 || o.Amount <= 0 {
		return false
	}
	if o.Side != types.Buy && o.Side != types.Sell {
		return false
	}
	return true
}
