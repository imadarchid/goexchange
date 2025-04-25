package order

import (
	"exchange/internal/types"
	"testing"
)

func TestNewOrder(t *testing.T) {
	order := NewOrder(10, 100, types.Buy, types.Limit)
	if len(order.ID) == 0 {
		t.Errorf("The order submitted does not have a valid ID")
	}
}

func TestOrderNegativeParams(t *testing.T) {
	order := NewOrder(-4, 0, types.Buy, types.Limit)
	if order.IsValid() == true {
		t.Errorf("Order has illegal parameters.")
	}
}

func TestOrderNoSides(t *testing.T) {
	order := NewOrder(10, 100, "S", "B")
	if order.IsValid() == true {
		t.Errorf("Order has illegal parameters.")
	}
}
