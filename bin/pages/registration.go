package pages

import (
	"github.com/gorilla/csrf"
	log "github.com/sirupsen/logrus"
	"html/template"
	"net/http"
	"web2/bin"
	"web2/bin/structures"
)

type RegistrationPageObject struct {
	BaseObject structures.BaseObject
}

func Registration(w http.ResponseWriter, r *http.Request) {
	var pageObject RegistrationPageObject
	if r.Method == "GET" {
		t, _ := template.ParseFiles("views/registration.html", "views/templates/error.html", "views/templates/menu.html")
		pageObject.BaseObject.CsrfField = csrf.TemplateField(r)
		err := t.Execute(w, pageObject)
		bin.CheckErr(err)
	} else {
		bin.SaveLog(log.Fields{
			"group":    "authorization",
			"subgroup": "registration",
		}, log.TraceLevel, "Registration attempt...")
		if r.FormValue("password") == r.FormValue("passwordAgain") {
			user := map[string]string{
				"Email":    r.FormValue("email"),
				"Password": r.FormValue("password"),
			}
			var result map[string]string
			response := bin.Request("/registration", "POST", "", user, &result)
			if response.StatusCode == 201 {
				bin.SaveLog(log.Fields{
					"group":    "authorization",
					"subgroup": "registration",
					"email":    r.FormValue("email"),
				}, log.InfoLevel, "Successful registration")
				http.Redirect(w, r, "/login", http.StatusSeeOther)
			} else {
				bin.SaveLog(log.Fields{
					"group":    "authorization",
					"subgroup": "registration",
					"error":    result["message"],
					"email":    r.FormValue("email"),
				}, log.ErrorLevel, "Registration error")
				pageObject.BaseObject.ErrorStr = result["message"]
				pageObject.BaseObject.CsrfField = csrf.TemplateField(r)
				t, _ := template.ParseFiles("views/registration.html", "views/templates/error.html", "views/templates/menu.html")
				err := t.Execute(w, pageObject)
				bin.CheckErr(err)
			}
		} else {
			bin.SaveLog(log.Fields{
				"group":    "authorization",
				"subgroup": "registration",
				"error":    "Passwords doesn't match!",
				"email":    r.FormValue("email"),
			}, log.ErrorLevel, "Registration error")
			pageObject.BaseObject.ErrorStr = "Введённые пароли не совпадают"
			pageObject.BaseObject.CsrfField = csrf.TemplateField(r)
			t, _ := template.ParseFiles("views/registration.html", "views/templates/error.html", "views/templates/menu.html")
			err := t.Execute(w, pageObject)
			bin.CheckErr(err)
		}
	}
}
