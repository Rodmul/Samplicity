package app

import (
	"DriveApi/internal/model"
	"DriveApi/internal/session"
	"fmt"
	"net/http"
)

type ProfileInfo struct {
	Username string
	Email    string
	Liked    int
	Created  int
}

func (srv *server) getProfileInfo(w http.ResponseWriter, r *http.Request) (*ProfileInfo, error) {
	info := &ProfileInfo{}
	token := session.GetSession(r)
	userID, err := srv.store.User().ParseToken(token)
	if err != nil {
		srv.Logger.Println(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return nil, fmt.Errorf("failed to parse token; %w", err)
	}

	user := &model.User{}
	createdSamples, err := srv.store.Sample().GetUserCreatedAmount(userID)
	if err != nil {
		srv.Logger.Printf("failed to get user created samples amount from database")
		return nil, err
	}
	user, err = srv.store.User().GetUserByID(userID)
	if err != nil {
		srv.Logger.Printf("failed to get user by id from database")
		return nil, err
	}

	info.Username = user.Username
	info.Email = user.Email
	info.Created = createdSamples
	info.Liked = 0

	return info, nil
}
