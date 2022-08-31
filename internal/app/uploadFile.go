package app

import (
	"github.com/200sc/klangsynthese"
	"github.com/200sc/klangsynthese/audio/filter"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

func (srv *server) uploadFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			srv.Logger.Printf("Error while uploading file to server; %v", err)
		}

		file, _, err := r.FormFile("myFile")
		if err != nil {
			srv.Logger.Printf("Error retrieving the file; %v", err)

			return
		}

		name := r.FormValue("fileName")
		pitch := r.FormValue("pitch")
		start := r.FormValue("start")
		end := r.FormValue("end")
		srv.Logger.Println(name, pitch, start, end)

		defer func(file multipart.File) {
			err := file.Close()
			if err != nil {
				log.Fatalf("Failed to close file; %v", err)
			}
		}(file)
		/*srv.Logger.Printf("Uploaded file: %v\n", handler.Filename)
		srv.Logger.Println("File Size: %v\n", handler.Size)
		srv.Logger.Println("MIME Header: %v\n", handler.Header)*/

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

		_, err = klangsynthese.LoadFile("./samples/" + name)
		if err != nil {
			srv.Logger.Printf("Failed to modify the audio; %v", err)
		}

		hQShifter := filter.HighQualityShifter
		hQShifter.PitchShift(2.0)

		w.WriteHeader(http.StatusOK)
	}
}
