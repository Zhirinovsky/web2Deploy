package models

type Set struct {
	ID               int
	Value            float64
	ProductID        int `db:"product_id"`
	CharacteristicID int `db:"characteristic_id"`
	Product          Product
	Characteristic   Characteristic
}
