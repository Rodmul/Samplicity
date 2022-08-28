package app

import (
	"fmt"
	"github.com/spf13/viper"
	"html/template"
	"log"
	"net/http"
)

type server struct {
	Logger *log.Logger
	router *http.ServeMux
}

type pageData struct {
	Title string
	Body  string
}

const (
	baseURL = "/"
	library = "library/"
	profile = "profile/"
	create  = "create/"
	static  = "static/"
)

func NewServer() *server {
	srv := &server{Logger: log.Default(), router: http.NewServeMux()}
	srv.registerHandlers()
	return srv
}

func (srv *server) Start() error {
	//srv.router.HandleFunc(baseURL+start, srv.handleStart)
	srv.router.HandleFunc(baseURL+create, srv.handleCreate())
	srv.router.HandleFunc(baseURL, srv.handleLibrary())
	srv.router.Handle(baseURL+filePrefix, srv.handleFile())
	fileServer := http.FileServer(http.Dir("./static/"))
	srv.router.Handle(baseURL+static, http.StripPrefix("/static", fileServer))

	//srv.router.HandleFunc(baseURL, srv.Index)
	srv.Logger.Printf("Start to listen to port %s", viper.GetString("port"))
	return fmt.Errorf("failed to listen and serve: %w", http.ListenAndServe(viper.GetString("port"), srv.router))
}

func (srv *server) registerHandlers() {

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

func (srv *server) Index(w http.ResponseWriter, r *http.Request) {
	srv.Logger.Print("index func")
}

func (srv *server) handleStart(w http.ResponseWriter, r *http.Request) {
	srv.Logger.Print("handleStart func")
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
