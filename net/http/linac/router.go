package linac

import (
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// Router model
type Router struct {
	routes []*Route
	engine *Engine
}

// Handle http
func (router *Router) Handle(pattern string, handler http.HandlerFunc) {
	parts := strings.Split(pattern, "/")
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
	pattern = strings.Join(parts, "/")
	log.Print(pattern)
	regex, regexErr := regexp.Compile(pattern)
	if regexErr != nil {
		panic(regexErr)
	}
	router.routes = append(router.routes, &Route{
		Regex:   regex,
		Params:  params,
		Handler: handler,
	})
}

// ServeHTTP ServeHTTP
func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range router.routes {
		if !route.Regex.MatchString(r.RequestURI) {
			continue
		}
		matches := route.Regex.FindStringSubmatch(r.RequestURI)

		//double check that the Route matches the URL pattern.
		if len(matches[0]) != len(r.RequestURI) {
			continue
		}

		params := make(map[string]string)

		if len(route.Params) > 0 {
			//add url parameters to the query param map
			values := r.URL.Query()
			for i, match := range matches[1:] {
				values.Add(route.Params[i], match)
				params[route.Params[i]] = match
				log.Print(match)
				log.Print(route.Params[i])
			}
			log.Print(params)
			//reassemble query params and add to RawQuery
			r.URL.RawQuery = url.Values(values).Encode() + "&" + r.URL.RawQuery
			//r.URL.RawQuery = url.Values(values).Encode()
		}
		route.Handler(w, r)
	}
}
