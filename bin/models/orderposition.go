package models

type Position struct {
	ID            int
	CheckoutPrice float64 `db:"checkout_price"`
	Amount        int
	OrderID       int `db:"order_id"`
	ProductID     int `db:"product_id"`
	Order         Order
	Product       Product
}
