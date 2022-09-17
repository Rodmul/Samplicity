package session

import (
	"github.com/gorilla/securecookie"
	"net/http"
)

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

func GetSession(r *http.Request) (token string) {
	if cookie, err := r.Cookie("token"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("token", cookie.Value, &cookieValue); err == nil {
			token = cookieValue["token"]
		}
	}
	return token
}

func SetSession(token string, w http.ResponseWriter) {
	value := map[string]string{
		"token": token,
	}
	if encoded, err := cookieHandler.Encode("token", value); err == nil {
		cookie := &http.Cookie{
			Name:   "token",
			Value:  encoded,
			Path:   "/",
			MaxAge: 3600,
		}
		http.SetCookie(w, cookie)
	}
}

func ClearSession(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}
