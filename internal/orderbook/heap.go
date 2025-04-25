package orderbook

import (
	"exchange/internal/order"
)

type OrderHeap struct {
	orders []*order.Order
	isMax  bool
}

func NewOrderHeap(isMax bool) *OrderHeap {
	return &OrderHeap{
		orders: []*order.Order{},
		isMax:  isMax,
	}
}

func (h *OrderHeap) Insert(o *order.Order) {
	h.orders = append(h.orders, o)
	h.heapifyUp(len(h.orders) - 1)
}

func (h *OrderHeap) Delete() *order.Order {
	if len(h.orders) == 0 {
		return nil
	}
	top := h.orders[0]
	last := h.orders[len(h.orders)-1]
	h.orders = h.orders[:len(h.orders)-1]

	if len(h.orders) > 0 {
		h.orders[0] = last
		h.heapifyDown(0)
	}

	return top
}

func (h *OrderHeap) Len() int {
	return len(h.orders)
}

func (h *OrderHeap) Peek() *order.Order {
	if len(h.orders) == 0 {
		return nil
	}
	return h.orders[0]
}

func (h *OrderHeap) heapifyDown(index int) {
	for {
		swapIndex := index
		if left(index) < len(h.orders) && h.compareOrders(left(index), index) > 0 {
			swapIndex = left(index)
		}
		if right(index) < len(h.orders) && h.compareOrders(right(index), index) > 0 {
			swapIndex = right(index)
		}

		if swapIndex == index {
			break
		}
		h.swap(index, swapIndex)
		index = swapIndex
	}
}

func (h *OrderHeap) heapifyUp(index int) {
	for index > 0 && h.compareOrders(parent(index), index) < 0 {
		h.swap(parent(index), index)
		index = parent(index)
	}
}

func (h *OrderHeap) compareOrders(i, j int) int {
	if h.isMax {
		if h.orders[i].Price > h.orders[j].Price {
			return 1
		} else if h.orders[i].Price < h.orders[j].Price {
			return -1
		}
		return 0
	}
	if h.orders[i].Price < h.orders[j].Price {
		return 1
	} else if h.orders[i].Price > h.orders[j].Price {
		return -1
	}
	return 0
}

func (ob *OrderBook) removeFromHeap(h *OrderHeap, id string) bool {
	for i, val := range h.orders {
		if val.ID == id {
			lastIndex := len(h.orders) - 1
			if i != lastIndex {
				h.orders[i], h.orders[lastIndex] = h.orders[lastIndex], h.orders[i]
			}
			h.orders = h.orders[:lastIndex]

			h.heapifyUp(i)
			h.heapifyDown(i)

			return true
		}
	}
	return false
}

func parent(i int) int {
	return (i - 1) / 2
}

func left(i int) int {
	return 2*i + 1
}

func right(i int) int {
	return 2*i + 2
}

func (h *OrderHeap) swap(i1, i2 int) {
	h.orders[i1], h.orders[i2] = h.orders[i2], h.orders[i1]
}
