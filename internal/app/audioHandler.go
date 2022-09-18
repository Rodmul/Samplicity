package app

import (
	"DriveApi/internal/session"
	"DriveApi/store"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	samples = "samples/"
	created = "created/"
	liked   = "liked/"
)

func (srv *server) handleFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Path[len(baseURL+samples):]
		if name == "" {
			err := serveDirDB(w, srv.store)
			if err != nil {
				srv.Logger.Println(err)
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

func (srv *server) handleCreatedFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Path[len(baseURL+created):]
		if name == "" {
			err := serveCreatedDir(w, r, srv.store)
			if err != nil {
				srv.Logger.Println(err)
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

func (srv *server) handleLikedFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Path[len(baseURL+liked):]
		if name == "" {
			err := serveLikedDir(w, r, srv.store)
			if err != nil {
				srv.Logger.Println(err)
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

func serveCreatedDir(w http.ResponseWriter, r *http.Request, s *store.Store) error {

	token := session.GetSession(r)
	userID, err := s.User().ParseToken(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return fmt.Errorf("failed to parse token; %w", err)
	}

	samples, err := s.Sample().GetUserCreated(userID)
	if err != nil {
		return fmt.Errorf("failed to get created user samples; %w", err)
	}

	j := json.NewEncoder(w)
	if err := j.Encode(&samples); err != nil {
		return fmt.Errorf("failed to encode samples info to json format; %w", err)
	}

	return nil
}

func serveLikedDir(w http.ResponseWriter, r *http.Request, s *store.Store) error {
	token := session.GetSession(r)
	userID, err := s.User().ParseToken(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return fmt.Errorf("failed to parse token; %w", err)
	}

	samples, err := s.LikedSample().GetUserAll(userID)
	if err != nil {
		return fmt.Errorf("failed to get created user samples; %w", err)
	}

	j := json.NewEncoder(w)
	if err := j.Encode(&samples); err != nil {
		return fmt.Errorf("failed to encode samples info to json format; %w", err)
	}

	return nil
}

func serveDirDB(w http.ResponseWriter, s *store.Store) error {
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
