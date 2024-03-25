package models

type Order struct {
	ID        int
	Date      string
	Address   string
	StatusID  int `db:"status_id"`
	UserID    int `db:"user_id"`
	Status    Status
	User      User
	Positions []Position
}
