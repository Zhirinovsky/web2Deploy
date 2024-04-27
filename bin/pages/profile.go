package pages

import (
	"github.com/Masterminds/sprig/v3"
	"github.com/gorilla/csrf"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
	"strconv"
	"time"
	"web2/bin"
	"web2/bin/models"
	"web2/bin/structures"
)

type ProfilePageObject struct {
	BaseObject structures.BaseObject
	Orders     []models.Order
	AllData    bool
}

func Profile(w http.ResponseWriter, r *http.Request) {
	bin.RefreshAPI()
	pageObject := getProfileData(bin.GetCurrentUser(w, r))
	pageObject.BaseObject.CsrfField = csrf.TemplateField(r)
	if pageObject.BaseObject.ErrorStr == "" {
		pageObject.BaseObject.ErrorStr = bin.GetError(&bin.GlobalError)
	}
	bin.Refresh(w, r)
	t, err := template.New("profile.html").Funcs(sprig.FuncMap()).ParseFiles("views/profile.html")
	template.Must(t.ParseGlob("views/templates/*.html"))
	template.Must(t.ParseGlob("views/templates/modals/*.html"))
	bin.NullErrorCheck(&pageObject.BaseObject.ErrorStr, err)
	err = t.Execute(w, pageObject)
	bin.CheckErr(err)
}

func getProfileData(user models.User) ProfilePageObject {
	var pageObject ProfilePageObject
	response := bin.Request("/orders", "GET", bin.ServerToken, nil, &pageObject.Orders)
	bin.ResponseCheck(response, "/orders", "GET")
	pageObject.BaseObject.CurrentUser = user
	if user.Email == "" || user.LastName == "" || user.Name == "" || user.Phone == "" || user.Gender == "" || user.Role.Name == "" {
		pageObject.AllData = false
	} else {
		pageObject.AllData = true
	}
	return pageObject
}

func ChangePersonalData(w http.ResponseWriter, r *http.Request) {
	user := bin.GetCurrentUser(w, r)
	user.LastName = r.FormValue("LastName")
	user.Name = r.FormValue("Name")
	user.MiddleName = r.FormValue("MiddleName")
	user.Gender = r.FormValue("Gender")
	user.Phone = r.FormValue("Phone")
	requestUser := log.Fields{
		"group":       "profile",
		"id":          user.ID,
		"email":       user.Email,
		"phone":       user.Phone,
		"last_name":   user.LastName,
		"name":        user.Name,
		"middle_name": user.MiddleName,
		"gender":      user.Gender,
		"role_id":     user.RoleID,
	}
	bin.SaveLog(requestUser, log.TraceLevel, "Sending profile data change request...")
	var result map[string]string
	response := bin.Request("/users/"+strconv.Itoa(user.ID), "PUT", bin.GetCurrentToken(w, r), user, &result)
	if !(response.StatusCode == 200 || response.StatusCode == 201) {
		bin.GlobalError = result["message"]
		requestUser["error"] = result["message"]
		bin.SaveLog(requestUser, log.ErrorLevel, "Failed profile data change")
	} else {
		bin.SaveLog(requestUser, log.InfoLevel, "Successful profile data change")
	}
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func ChangePassword(w http.ResponseWriter, r *http.Request) {
	user := bin.GetCurrentUser(w, r)
	requestPassword := log.Fields{
		"group":           "profile",
		"id":              user.ID,
		"email":           user.Email,
		"old_password":    r.FormValue("OldPassword"),
		"new_password":    r.FormValue("NewPassword"),
		"repeat_password": r.FormValue("RepeatPassword"),
	}
	bin.SaveLog(requestPassword, log.TraceLevel, "Sending password change request...")
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(r.FormValue("OldPassword")))
	if err == nil {
		if r.FormValue("NewPassword") == r.FormValue("RepeatPassword") {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("NewPassword")), bcrypt.DefaultCost)
			bin.GlobalCheck(err)
			var result map[string]string
			user.Password = string(hashedPassword)
			response := bin.Request("/users/"+strconv.Itoa(user.ID), "PUT", bin.GetCurrentToken(w, r), user, &result)
			if !(response.StatusCode == 200 || response.StatusCode == 201) {
				if bin.GlobalError == "" {
					bin.GlobalError = result["message"]
				}
				requestPassword["error"] = result["message"]
				bin.SaveLog(requestPassword, log.ErrorLevel, "Failed password change")
			} else {
				bin.SaveLog(requestPassword, log.InfoLevel, "Successful password change")
			}
		} else {
			requestPassword["error"] = "Password mismatched."
			bin.SaveLog(requestPassword, log.ErrorLevel, "Failed password change")
			bin.GlobalError = "Введёные новые пароли не совпадают"
		}
	} else {
		requestPassword["error"] = "Incorrect password."
		bin.SaveLog(requestPassword, log.ErrorLevel, "Failed password change")
		bin.GlobalError = "Введён неверный пароль"
	}
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func CreateCard(w http.ResponseWriter, r *http.Request) {
	user := bin.GetCurrentUser(w, r)
	card := models.Card{
		ID:       user.ID,
		Date:     time.Now().Format(time.DateTime),
		Discount: 0,
	}
	requestCard := log.Fields{
		"group":    "profile",
		"email":    user.Email,
		"id":       card.ID,
		"date":     card.Date,
		"discount": card.Discount,
	}
	bin.SaveLog(requestCard, log.TraceLevel, "Sending discount card creation request...")
	var result map[string]string
	response := bin.Request("/cards", "POST", bin.GetCurrentToken(w, r), card, &result)
	if !(response.StatusCode == 200 || response.StatusCode == 201) {
		bin.GlobalError = result["message"]
		requestCard["error"] = result["message"]
		bin.SaveLog(requestCard, log.ErrorLevel, "Failed discount card creation")
	} else {
		bin.SaveLog(requestCard, log.InfoLevel, "Successful discount card creation")
	}
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func DeleteCard(w http.ResponseWriter, r *http.Request) {
	user := bin.GetCurrentUser(w, r)
	requestCard := log.Fields{
		"group": "profile",
		"email": user.Email,
	}
	bin.SaveLog(requestCard, log.TraceLevel, "Sending discount card deletion request...")
	var result map[string]string
	response := bin.Request("/cards/"+strconv.Itoa(user.ID), "DELETE", bin.GetCurrentToken(w, r), nil, &result)
	if !(response.StatusCode == 200 || response.StatusCode == 201) {
		bin.GlobalError = result["message"]
		requestCard["error"] = result["message"]
		bin.SaveLog(requestCard, log.ErrorLevel, "Failed discount card deletion")
	} else {
		bin.SaveLog(requestCard, log.InfoLevel, "Successful discount card deletion")
	}
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}
