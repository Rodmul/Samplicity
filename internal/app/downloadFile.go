package app

import (
	"fmt"
	"net/http"
)

func (srv *server) downloadFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			srv.Logger.Printf("Error while uploading file to server; %v", err)
		}
		sampleName := r.FormValue("sampleName")
		srv.Logger.Println(sampleName)

		sample, err := srv.store.Sample().GetByName(sampleName)
		path := fmt.Sprintf(sample.Path + sample.Name + "." + sample.Type)
		srv.Logger.Println(path)
		if err != nil {
			srv.Logger.Println("failed to get sample by name; %v", err)
		}

		http.ServeFile(w, r, path)
	}
}
