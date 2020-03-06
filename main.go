package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

type server struct {
	Things map[int]thing
	Router *http.ServeMux
}

func newServer() *server {
	s := &server{
		Router: http.NewServeMux(),
		Things: map[int]thing{},
	}
	s.routes()
	s.populateThings()
	return s
}

func (s *server) routes() {
	ms := s.recoverMiddleware
	s.Router.Handle("/things", ms(s.handleThings()))
}

func (s *server) errorJSON(w http.ResponseWriter, error interface{}, code int) {
	wrapper := struct{ Data interface{} }{Data: error}
	resData, err := json.Marshal(wrapper)
	if err != nil {
		http.Error(w, "Really big bad could not even json the error", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	w.Write(resData)
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, data interface{}, status int) {
	resData, err := json.Marshal(data)
	if err != nil {
		s.errorJSON(w, "no json for u", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(resData)
}

func (s *server) populateThings() {
	s.Things = map[int]thing{
		1: thing{
			ThingId: 1,
			Name:    "First Thing.",
		},
		2: thing{
			ThingId: 2,
			Name:    "Grants Thing",
		},
		3: thing{
			ThingId: 3,
			Name:    "Garrets Thing",
		},
		4: thing{
			ThingId: 4,
			Name:    "Correys Thing",
		},
		5: thing{
			ThingId: 5,
			Name:    "Some Thing",
		},
		6: thing{
			ThingId: 6,
			Name:    "Universe Thing",
		},
		7: thing{
			ThingId: 7,
			Name:    "Underverse Thing",
		},
	}
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	s := newServer()
	err := http.ListenAndServe(":3000", s)
	if err != nil {
		return err
	}
	return nil
}

func (s *server) recoverMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer func() {
			r := recover()
			if r != nil {
				switch t := r.(type) {
				case string:
					err = errors.New(t)
				case error:
					err = t
				default:
					err = errors.New("Critical error.")
				}
				s.errorJSON(w, err.Error(), http.StatusInternalServerError)
			}
		}()
		h.ServeHTTP(w, r)
	}
}

type thing struct {
	ThingId int    `json:"thingId"`
	Name    string `json:"name"`
}

func (s *server) handleThings() http.HandlerFunc {
	type response struct {
		Data map[int]thing `json:"data"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			s.handleGETThings()(w, r)
		default:
			s.errorJSON(w, "Not supported", http.StatusMethodNotAllowed)
		}
	}
}

func (s *server) handleGETThings() http.HandlerFunc {
	type response struct {
		Data map[int]thing `json:"data"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		response := response{Data: s.Things}
		s.respond(w, r, response, 200)
	}
}

// REQUIREMENTS ------------------

// panics must not reach the user

// support GET at /things
