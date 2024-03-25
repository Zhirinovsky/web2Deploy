package pages

import (
	"github.com/gorilla/csrf"
	log "github.com/sirupsen/logrus"
	"html/template"
	"net/http"
	"web2/bin"
	"web2/bin/models"
	"web2/bin/structures"
)

type LoginPageObject struct {
	BaseObject structures.BaseObject
}

func Login(w http.ResponseWriter, r *http.Request) {
	var pageObject LoginPageObject
	if r.Method == "GET" {
		_, err := r.Cookie("session_token")
		if err == nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			pageObject.BaseObject.CsrfField = csrf.TemplateField(r)
			t, _ := template.ParseFiles("views/login.html", "views/templates/error.html", "views/templates/menu.html")
			err = t.Execute(w, pageObject)
			bin.CheckErr(err)
		}
	} else {
		bin.SaveLog(log.Fields{
			"group":    "authorization",
			"subgroup": "login",
		}, log.TraceLevel, "Authorization attempt...")
		user := map[string]string{
			"Email":    r.FormValue("email"),
			"Password": r.FormValue("password"),
		}
		var result map[string]string
		response := bin.Request("/login", "POST", "", user, &result)
		if response.StatusCode == 200 {
			var authForm models.AuthForm
			authForm.Email = r.FormValue("email")
			authForm.Token = result["token"]
			authForm.Role = result["role"]
			bin.SaveLog(log.Fields{
				"group":    "authorization",
				"subgroup": "login",
				"email":    r.FormValue("email"),
				"role":     result["role"],
			}, log.InfoLevel, "Successful authorization")
			bin.Authorize(w, r, authForm)
		} else {
			bin.SaveLog(log.Fields{
				"group":    "authorization",
				"subgroup": "login",
				"error":    result["message"],
				"email":    r.FormValue("email"),
			}, log.ErrorLevel, "Authorization error")
			pageObject.BaseObject.ErrorStr = result["message"]
			pageObject.BaseObject.CsrfField = csrf.TemplateField(r)
			t, _ := template.ParseFiles("views/login.html", "views/templates/error.html", "views/templates/menu.html")
			err := t.Execute(w, pageObject)
			bin.CheckErr(err)
		}
	}
}
