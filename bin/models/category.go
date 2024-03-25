package models

type Category struct {
	ID       int
	Name     string
	Relation int
	IsExist  bool `db:"is_exist"`
}
