package linac

import (
	"log"
	"net/http"
	"regexp"
	"strings"
)

// RouteGroup RouteGroup
type RouteGroup struct {
	router   Router
	Prefixed string
	Routes   []*Route
}

// Group Group
func (group *RouteGroup) Group(path string, fu func(*RouteGroup) *RouteGroup) {
	xGroup := &RouteGroup{
		Prefixed: group.Prefixed,
		router:   group.router,
	}
	xGroup = fu(xGroup)
	for _, route := range xGroup.Routes {
		group.Routes = append(group.Routes, route)
	}
}

func (group *RouteGroup) addRoute(path string, handler http.HandlerFunc) {
	parts := strings.Split(path, "/")
	j := 0
	params := make(map[int]string)
	for i, part := range parts {
		if strings.HasPrefix(part, ":") {
			expr := "([^/]+)"

			// a user may choose to override the defult expression
			// similar to expressjs: ‘/user/:id([0-9]+)’
			if index := strings.Index(part, "("); index != -1 {
				expr = part[index:]
				part = part[:index]
			}
			params[j] = part[1:]
			parts[i] = expr
			j++
		}
	}
	path = strings.Join(parts, "/")
	log.Print(path)
	regex, regexErr := regexp.Compile(path)
	if regexErr != nil {
		panic(regexErr)
	}
	group.Routes = append(group.Routes, &Route{
		Regex:   regex,
		Params:  params,
		Handler: handler,
	})
}

// Route model
type Route struct {
	Regex   *regexp.Regexp
	Params  map[int]string
	Handler http.HandlerFunc
}
