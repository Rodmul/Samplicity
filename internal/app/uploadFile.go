package app

import (
	"DriveApi/internal/model"
	"DriveApi/internal/session"
	"fmt"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
)

func (srv *server) uploadFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		token := session.GetSession(r)
		userID, err := srv.store.User().ParseToken(token)
		if err != nil {
			srv.Logger.Println(err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		err = r.ParseMultipartForm(10 << 20)
		if err != nil {
			srv.Logger.Printf("Error while uploading file to server; %v", err)
		}

		file, _, err := r.FormFile("myFile")
		if err != nil {
			srv.Logger.Printf("Error retrieving the file; %v", err)

			return
		}

		name := r.FormValue("fileName")
		fileType := r.FormValue("type")
		start := r.FormValue("start")
		end := r.FormValue("end")

		author, err := srv.store.User().GetUserByID(userID)
		if err != nil {
			srv.Logger.Println("failed to get user by id")

			return
		}

		sample := &model.Sample{Name: name, Author: author.Username, AuthorID: userID, Path: "./samples/", Type: fileType}

		startInt, err := strconv.Atoi(start)
		if err != nil {
			srv.Logger.Println("Failed to convert start values")
		}
		endInt, err := strconv.Atoi(end)
		if err != nil {
			srv.Logger.Println("Failed to convert start values")
		}

		srv.Logger.Println(name, fileType, start, end)

		defer func(file multipart.File) {
			err := file.Close()
			if err != nil {
				log.Fatalf("Failed to close file; %v", err)
			}
		}(file)

		dst, err := os.Create("./samples/" + name)
		defer func(dst *os.File) {
			err := dst.Close()
			if err != nil {
				log.Fatalf("Failed to close file; %v", err)
			}
		}(dst)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		path := fmt.Sprintf(name + "." + fileType)
		err = ffmpeg.Input("./samples/"+name, ffmpeg.KwArgs{"ss": startInt, "to": endInt}).
			Output("./samples/" + path).OverWriteOutput().Run()
		if err != nil {
			srv.Logger.Printf("failed to cut audio %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		_, err = srv.store.Sample().Create(sample)
		if err != nil {
			srv.Logger.Printf("failed to create sample; %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		/*createdSample := &model.CreatedSample{
			SampleID: sampleID,
			UserID:   userID,
		}

		err = srv.store.CreatedSample().Create(createdSample)
		if err != nil {
			srv.Logger.Printf("failed to create createdSample; %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}*/
		w.WriteHeader(http.StatusOK)
	}
}
