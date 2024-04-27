package structures

import "web2/bin/models"

type BaseObject struct {
	CurrentUser models.User
	ErrorStr    string
	MessageStr  string
	CsrfField   interface{}
}
