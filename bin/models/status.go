package models

type Status struct {
	ID      int
	Status  string
	IsExist bool `db:"is_exist"`
}
