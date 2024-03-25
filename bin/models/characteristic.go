package models

type Characteristic struct {
	ID       int
	Name     string
	Type     string
	Relation int
	Products []Product
	IsExist  bool `db:"is_exist"`
}
