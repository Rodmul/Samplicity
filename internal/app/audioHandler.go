package app

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type FileInfo struct {
	Name  string      `json:"Name"`
	IsDir bool        `json:"IsDir"`
	Mode  os.FileMode `json:"Mode"`
}

const (
	filePrefix = "samples/"
	root       = "./samples"
)

func (srv *server) handleFile() http.HandlerFunc {
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
		fileInfos[i].IsDir = files[i].IsDir()
		fileInfos[i].Mode = files[i].Mode()
	}

	j := json.NewEncoder(w)

	if err := j.Encode(&fileInfos); err != nil {
		panic(err)
	}
}
