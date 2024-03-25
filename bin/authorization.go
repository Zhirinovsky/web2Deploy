package bin

import (
	"context"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
	"web2/bin/models"
)

func Authorize(w http.ResponseWriter, r *http.Request, user models.AuthForm) {
	_, err := r.Cookie("session_token")
	if err != nil {
		sessionToken := uuid.NewString()
		err := Client.HMSet(context.Background(), "sessions:"+sessionToken, "login", user.Email, "token", user.Token).Err()
		CheckErr(err)
		Client.Expire(context.Background(), "sessions:"+sessionToken, 2*time.Hour)
		cookie := &http.Cookie{
			Name:     "session_token",
			Value:    sessionToken,
			Path:     "/",
			Expires:  time.Now().Add(120 * time.Minute),
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		}
		http.SetCookie(w, cookie)
	}
	if user.Role == "Клиент" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	}
}

func Refresh(w http.ResponseWriter, r *http.Request) {
	UpdateMetrics()
	c, err := r.Cookie("session_token")
	if err != nil {
		return
	}
	user := Client.HGetAll(context.Background(), "sessions:"+c.Value).Val()
	if len(user) == 0 {
		return
	}
	Client.Del(context.Background(), "sessions:"+c.Value)
	var result map[string]string
	response := Request("/refresh", "GET", user["token"], nil, &result)
	if response.StatusCode == 200 {
		sessionToken := uuid.NewString()
		Client.HMSet(context.Background(), "sessions:"+sessionToken, "login", user["login"], "token", result["token"])
		Client.Expire(context.Background(), "sessions:"+sessionToken, 2*time.Hour)
		cookie := &http.Cookie{
			Name:     "session_token",
			Value:    sessionToken,
			Path:     "/",
			Expires:  time.Now().Add(120 * time.Minute),
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		}
		http.SetCookie(w, cookie)
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		return
	}
	user := Client.HGetAll(context.Background(), "sessions:"+c.Value).Val()
	if len(user) != 0 {
		err = Client.Del(context.Background(), "sessions:"+c.Value).Err()
		CheckErr(err)
		Request("/logout", "GET", user["token"], nil, nil)
	}
	cookie := &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)
	SaveLog(log.Fields{
		"group":    "authorization",
		"subgroup": "logout",
		"email":    user["login"],
	}, log.InfoLevel, "User logout")
	UpdateMetrics()
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func GetCurrentUser(w http.ResponseWriter, r *http.Request) models.User {
	userNew := models.User{}
	c, err := r.Cookie("session_token")
	if err == nil {
		user := Client.HGetAll(context.Background(), "sessions:"+c.Value).Val()
		if len(user) != 0 {
			Request("/getcurrentuser", "GET", user["token"], nil, &userNew)
		} else {
			cookie := &http.Cookie{
				Name:     "session_token",
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteLaxMode,
			}
			http.SetCookie(w, cookie)
		}
	}
	return userNew
}

func GetCurrentToken(w http.ResponseWriter, r *http.Request) string {
	token := ""
	c, err := r.Cookie("session_token")
	if err == nil {
		user := Client.HGetAll(context.Background(), "sessions:"+c.Value).Val()
		if len(user) != 0 {
			token = user["token"]
		}
	}
	return token
}
