package order

import (
	"context"
	"exchange/internal/db"
	"sync"

	"github.com/google/uuid"
)

type OrderService struct {
	orderbook OrderBookInterface
	db        *db.Queries
	mu        sync.RWMutex
}

// OrderBookInterface defines the interface for orderbook operations
type OrderBookInterface interface {
	Submit(*Order) bool
	Withdraw(*Order) bool
}

func NewOrderService(orderbook OrderBookInterface, db *db.Queries) *OrderService {
	return &OrderService{
		orderbook: orderbook,
		db:        db,
	}
}

// SubmitOrder handles the submission of a new order, ensuring consistency between memory and database
func (s *OrderService) SubmitOrder(ctx context.Context, o *Order) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// First, try to submit to the orderbook
	success := s.orderbook.Submit(o)
	if !success {
		return nil // Order was rejected by the orderbook
	}

	// If successful in memory, persist to database
	_, err := s.db.CreateOrder(ctx, db.CreateOrderParams{
		Price:     o.Price,
		Amount:    o.Amount,
		Side:      o.Side,
		OrderType: o.Type,
		Asset:     o.Ticker,
		CreatedBy: o.CreatedBy,
	})

	if err != nil {
		// If database operation fails, we need to rollback the memory operation
		s.orderbook.Withdraw(o)
		return err
	}

	return nil
}

// UpdateOrderStatus updates both memory and database order status
func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status db.OrderStatusType) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Update in database first
	_, err := s.db.UpdateOrderStatus(ctx, db.UpdateOrderStatusParams{
		OrderStatus: status,
		ID:          orderID,
	})
	if err != nil {
		return err
	}

	// If database update successful, update in memory
	// Note: This is a simplified version. In a real implementation, you'd need to
	// find the order in the appropriate heap and update its status
	return nil
}

// GetOrder retrieves an order from the database
func (s *OrderService) GetOrder(ctx context.Context, orderID uuid.UUID) (*Order, error) {
	dbOrder, err := s.db.GetOrderById(ctx, orderID)
	if err != nil {
		return nil, err
	}

	return &Order{
		ID:        dbOrder.ID,
		Price:     dbOrder.Price,
		Amount:    dbOrder.Amount,
		Side:      dbOrder.Side,
		Time:      dbOrder.CreatedAt,
		Type:      dbOrder.OrderType,
		Status:    dbOrder.OrderStatus,
		Ticker:    dbOrder.Asset,
		CreatedBy: dbOrder.CreatedBy,
	}, nil
}
