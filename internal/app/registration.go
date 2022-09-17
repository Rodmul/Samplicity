package app

import (
	"DriveApi/internal/model"
	"net/http"
)

func (srv *server) handleRegistration() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			srv.Logger.Printf("failed to parse form; %v", err)
		}

		username := r.FormValue("username")
		password := r.FormValue("password")
		email := r.FormValue("email")

		user := &model.User{Username: username, Password: password, Email: email}

		err = srv.store.User().Create(user)
		if err != nil {
			srv.Logger.Printf("failed to create user; %v", err)
		}
	}
}
