package order

import (
	"exchange/internal/db"
	"testing"

	"github.com/google/uuid"
)

func TestNewOrder(t *testing.T) {
	if len(NewOrder(10, 100, db.OrderSideTypeBUY, db.OrderTypeLIMIT, "LINK", uuid.New()).ID) == 0 {
		t.Errorf("The order submitted does not have a valid ID")
	}
}

func TestOrderNegativeParams(t *testing.T) {
	order := NewOrder(-4, 0, db.OrderSideTypeBUY, db.OrderTypeLIMIT, "LINK", uuid.New())
	if order.IsValid() == true {
		t.Errorf("Order has illegal parameters.")
	}
}

func TestOrderNoSides(t *testing.T) {
	order := NewOrder(10, 100, "S", "B", "LINK", uuid.New())
	if order.IsValid() == true {
		t.Errorf("Order has illegal parameters.")
	}
}

func TestValidOrder(t *testing.T) {
	order := NewOrder(10, 100, db.OrderSideTypeBUY, db.OrderTypeLIMIT, "LINK", uuid.New())
	if order.IsValid() != true {
		t.Errorf("Expected valid order.")
	}
}
