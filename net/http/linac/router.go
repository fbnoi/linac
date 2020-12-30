package linac

import (
	"log"
	"net/http"
	"net/url"
)

// Router model
type Router struct {
	routes []*Route
	engine *Engine
}

// AddRoute 向路由器中添加路由
func (router *Router) AddRoute(path, method string, handler Handler) *Router {
	router.engine.addRoute(path, method, handler)
	return router
}

// handleFunc 添加路由处理方法
// pattern 路由模式，必须以 '/' 开头，
func (router *Router) handleFunc(pattern, method string, handler http.HandlerFunc) {
	if pattern[0] != '/' {
		panic("pattern must start with '/'")
	}
	router.routes = append(router.routes, newRoute(pattern, method, handler))
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
