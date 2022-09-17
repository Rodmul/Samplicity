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
	baseURL       = "/"
	profile       = "profile/"
	create        = "create/"
	static        = "static/"
	upload        = "upload/"
	auth          = "auth/"
	login         = "login/"
	registration  = "registration/"
	authorization = "authorization/"
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
	srv.router.HandleFunc(baseURL+samples, srv.handleFile())
	srv.router.HandleFunc(baseURL+profile, srv.userIdentity(srv.handleProfile()))
	srv.router.HandleFunc(baseURL+upload, srv.userIdentity(srv.uploadFile()))
	srv.router.HandleFunc(baseURL+auth, srv.handleAuth())
	srv.router.HandleFunc(baseURL+login, srv.handleLogin())
	srv.router.HandleFunc(baseURL+login+registration, srv.handleRegistration())
	srv.router.HandleFunc(baseURL+auth+authorization, srv.handleAuthorization())
	fileServer := http.FileServer(http.Dir("./static/"))
	srv.router.Handle(baseURL+static, http.StripPrefix("/static", fileServer))
}

func (srv *server) handleProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ts, err := template.ParseFiles("./static/templates/profilePage.html")
		if err != nil {
			srv.Logger.Printf(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		info, err := srv.getProfileInfo(w, r)
		if err != nil {
			srv.Logger.Println("failed to get profile info")
			return
		}

		err = ts.Execute(w, info)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (srv *server) handleAuth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ts, err := template.ParseFiles("./static/templates/authorizationPage.html")
		if err != nil {
			srv.Logger.Printf(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = ts.Execute(w, nil)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (srv *server) handleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ts, err := template.ParseFiles("./static/templates/registrationPage.html")
		if err != nil {
			srv.Logger.Printf(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = ts.Execute(w, nil)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (srv *server) handleCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ts, err := template.ParseFiles("./static/templates/createPage.html")
		if err != nil {
			srv.Logger.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = ts.Execute(w, nil)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (srv *server) handleLibrary() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ts, err := template.ParseFiles("./static/templates/libraryPage.html")
		if err != nil {
			srv.Logger.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = ts.Execute(w, nil)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
