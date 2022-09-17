package app

import (
	"DriveApi/internal/session"
	"net/http"
)

func (srv *server) handleAuthorization() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			srv.Logger.Printf("failed to parse form; %v", err)
		}

		username := r.FormValue("username")
		password := r.FormValue("password")

		token, err := srv.store.User().GenerateToken(username, password)
		if err != nil {
			srv.Logger.Printf("failed to generate token; %v", err)

			return
		}

		session.SetSession(token, w)

		srv.Logger.Printf(token)
		w.WriteHeader(http.StatusOK)
	}
}
