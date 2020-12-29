package main

import (
	"log"
	"net/http"
)

type mux struct {
	router routes
}

type routes []*route

type route struct {
	Path    string
	Handler http.HandlerFunc
}

func (m *mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range m.router {
		if route.Path == r.RequestURI {
			route.Handler(w, r)
			return
		}
	}
	log.Print("404 not found")
}

func (m *mux) Handle(path string, handler http.HandlerFunc) {
	m.router = append(m.router, &route{
		Path:    path,
		Handler: handler,
	})
}

func main() {
	m := &mux{}
	m.Handle("/", func(w http.ResponseWriter, r *http.Request) {
		log.Print("hello world")
	})
	log.Fatal(http.ListenAndServe(":80", m))
}
