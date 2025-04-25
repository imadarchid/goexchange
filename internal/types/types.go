package types

type Side string

const (
	Buy  Side = "BUY"
	Sell Side = "SELL"
)

type OrderType string

const (
	Market OrderType = "MARKET"
	Limit  OrderType = "LIMIT"
)

type Status string

const (
	Pending         Status = "PENDING"
	Filled          Status = "FILLED"
	PartiallyFilled Status = "PARTIAL_FILL"
	Cancelled       Status = "CANCELLED"
)
