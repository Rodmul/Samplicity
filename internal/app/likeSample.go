package app

import (
	"DriveApi/internal/model"
	"DriveApi/internal/session"
	"net/http"
)

func (srv *server) likeSample() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			srv.Logger.Printf("Error while uploading file to server; %v", err)
		}
		sampleName := r.FormValue("sampleName")
		srv.Logger.Println(sampleName)

		token := session.GetSession(r)
		userID, err := srv.store.User().ParseToken(token)
		if err != nil {
			srv.Logger.Println(err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		sample, err := srv.store.Sample().GetByName(sampleName)
		if err != nil {
			srv.Logger.Println("failed to get sample by name")
			return
		}

		likedSample := &model.LikedSamples{SampleID: sample.ID, UserID: userID}
		err = srv.store.LikedSample().Create(likedSample)
		if err != nil {
			srv.Logger.Println("failed to create likedSample")
		}
	}
}
