package app

import (
	"DriveApi/internal/repository"
	"DriveApi/store"
	"fmt"
	"github.com/spf13/viper"
	"html/template"
	"log"
	"net/http"
)

type server struct {
	Logger *log.Logger
	router *http.ServeMux
	store  *store.Store
}

const (
	baseURL = "/"
	profile = "profile/"
	create  = "create/"
	static  = "static/"
	upload  = "upload/"
)

func NewServer() *server {
	srv := &server{Logger: log.Default(), router: http.NewServeMux()}
	srv.registerHandlers()
	return srv
}

func (srv *server) Start() error {
	database, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: "1234",
	})
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}

	srv.store = store.New(database, srv.Logger)
	srv.Logger.Printf("Start to listen to port %s", viper.GetString("port"))
	return fmt.Errorf("failed to listen and serve: %w", http.ListenAndServe(viper.GetString("port"), srv.router))
}

func (srv *server) registerHandlers() {
	srv.router.HandleFunc(baseURL+create, srv.handleCreate())
	srv.router.HandleFunc(baseURL, srv.handleLibrary())
	srv.router.Handle(baseURL+filePrefix, srv.handleFile())
	srv.router.Handle(baseURL+upload, srv.uploadFile())
	fileServer := http.FileServer(http.Dir("./static/"))
	srv.router.Handle(baseURL+static, http.StripPrefix("/static", fileServer))
}

func (srv *server) handleCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ts, err := template.ParseFiles("./static/templates/createPage.html")
		if err != nil {
			srv.Logger.Println(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		err = ts.Execute(w, nil)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}

func (srv *server) handleLibrary() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ts, err := template.ParseFiles("./static/templates/libraryPage.html")
		if err != nil {
			srv.Logger.Println(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		err = ts.Execute(w, nil)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}
