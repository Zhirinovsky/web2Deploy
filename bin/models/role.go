package models

type Role struct {
	ID      int
	Name    string
	IsExist bool `db:"is_exist"`
}
