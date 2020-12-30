package linac

import (
	"net/http"
	"regexp"
)

// RouteGroup RouteGroup
type RouteGroup struct {
	router   Router
	Prefixed string
	Routes   []*Route
}

// Route model
type Route struct {
	Regex   *regexp.Regexp
	Params  map[int]string
	Handler http.HandlerFunc
}
