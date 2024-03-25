package models

type User struct {
	ID         int
	Email      string
	Password   string
	Phone      string
	LastName   string `db:"last_name"`
	Name       string
	MiddleName string `db:"middle_name"`
	Gender     string
	RoleID     int `db:"role_id"`
	Role       Role
	Card       Card
	IsExist    bool `db:"is_exist"`
}
