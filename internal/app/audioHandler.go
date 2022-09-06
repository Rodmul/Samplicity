package app

import (
	"DriveApi/store"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type FileInfo struct {
	Name string
}

const (
	filePrefix = "samples/"
	root       = "./samples"
)

func (srv *server) handleFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Path[len(baseURL+filePrefix):]
		if name == "" {
			err := serverDirDB(w, r, srv.store)
			if err != nil {
				srv.Logger.Fatal(err)
			}
			return
		}
		sample, err := srv.store.Sample().GetByName(name)
		path := fmt.Sprintf(sample.Path + sample.Name + "." + sample.Type)
		srv.Logger.Println(path)
		if err != nil {
			srv.Logger.Println("failed to get sample by name; %v", err)
		}
		http.ServeFile(w, r, path)
	}
}

func (srv *server) handleFiles() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join(root, r.URL.Path[len(baseURL+filePrefix)-1:])
		stat, err := os.Stat(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if stat.IsDir() {
			serveDir(w, r, path)
			return
		}
		http.ServeFile(w, r, path)
	}
}

func serveDir(w http.ResponseWriter, r *http.Request, path string) {
	defer func() {
		if err, ok := recover().(error); ok {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}()
	file, err := os.Open(path)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalf("Failed to close file; %v", err)
		}
	}(file)
	if err != nil {
		panic(err)
	}

	files, err := file.Readdir(-1)
	if err != nil {
		panic(err)
	}

	fileInfos := make([]FileInfo, len(files), len(files))

	for i := range files {
		fileInfos[i].Name = files[i].Name()
	}

	j := json.NewEncoder(w)

	if err := j.Encode(&fileInfos); err != nil {
		panic(err)
	}
}

func serverDirDB(w http.ResponseWriter, r *http.Request, s *store.Store) error {
	samples, err := s.Sample().GetAll()
	if err != nil {
		return fmt.Errorf("failed to get samples; %w", err)
	}

	j := json.NewEncoder(w)
	if err := j.Encode(&samples); err != nil {
		return fmt.Errorf("failed to encode samples info to json format; %w", err)
	}

	return nil
}
