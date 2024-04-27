package models

type Cart struct {
	ID        int
	Amount    int
	ProductID int `db:"product_id"`
	UserID    int `db:"user_id"`
	Active    bool
	Product   Product
	User      User
}
