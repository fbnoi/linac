package main

import (
	"net/http"
	"regexp"
)

type mux struct {
	router routes
}

type routes []*route

type route struct {
	regex   *regexp.Regexp
	params  map[int]string
	Handler http.HandlerFunc
}

func main() {
	// m := &mux{}
	// m.Handle("/test/:id(\\d+)", func(w http.ResponseWriter, r *http.Request) {
	// })
	// log.Fatal(http.ListenAndServe(":80", m))
}
