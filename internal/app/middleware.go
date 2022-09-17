package app

import (
	"DriveApi/internal/session"
	"net/http"
	"strconv"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userID"
)

func (srv *server) userIdentity(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		token := session.GetSession(r)
		r.Header.Set(authorizationHeader, token)

		header := r.Header.Get(authorizationHeader)
		if header == "" {
			http.Redirect(w, r, "http://localhost:8000/auth/", http.StatusTemporaryRedirect)
			srv.Logger.Println("empty auth header")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 1 {
			srv.Logger.Printf("invalid auth header")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if len(headerParts[0]) == 0 {
			srv.Logger.Println("token is empty")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userID, err := srv.store.User().ParseToken(headerParts[0])
		if err != nil {
			srv.Logger.Println(err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.Header().Set(userCtx, strconv.Itoa(userID))
		next.ServeHTTP(w, r)
	}
}
