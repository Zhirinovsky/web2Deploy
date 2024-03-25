package models

type Product struct {
	ID         int
	Name       string
	Price      float64
	Amount     int
	Discount   int
	ImageLink  string `db:"image_link"`
	CategoryID string `db:"category_id"`
	Category   Category
	Sets       []Set
	IsExist    bool `db:"is_exist"`
}
