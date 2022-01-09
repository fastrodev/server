package server

import (
	"net/http"
)

type Handler func(w http.ResponseWriter, r *http.Request)

type handler struct {
	routes       map[string]route
	staticPath   string
	staticFolder string
	serverless   bool
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route, ok := h.routes[h.getKey(r.Method, r.URL.Path)]
	if !ok {
		h.handleNotFoundRouteKey(w, r)
		return
	}

	if len(route.middlewares) > 0 {
		h.handleMiddleware(route, w, r)
		return
	}

	route.handler(w, r)
}

func (h *handler) getKey(incomingMethod string, incomingPath string) string {
	for _, r := range h.routes {
		if incomingMethod == r.method && validate(r.path, incomingPath) {
			return r.method + SPLITTER + r.path
		}
	}
	return NOTFOUND
}

func (h *handler) handleNotFoundRouteKey(w http.ResponseWriter, r *http.Request) {
	folder := h.staticFolder
	path := h.staticPath
	if h.serverless {
		folder = SERVERLESS_FOLDER + h.staticFolder
	}
	if path == EMPTY {
		path = SLASH
	}
	if folder == EMPTY {
		folder = TMP
	}
	fileHandler := http.FileServer(http.Dir(folder))
	http.StripPrefix(path, fileHandler).ServeHTTP(w, r)
}

func (h *handler) handleMiddleware(
	route route,
	w http.ResponseWriter,
	r *http.Request,
) {
	var next bool
	lengthOfRouteMiddleware := len(route.middlewares)
	if lengthOfRouteMiddleware > 0 {
		next, w, r = h.rangeMiddleware(
			route, route.middlewares, w, r, lengthOfRouteMiddleware,
		)
		if !next {
			return
		}
	}
	if next {
		route.handler(w, r)
	}
}

func (h *handler) rangeMiddleware(
	route route,
	middlewares []Middleware,
	w http.ResponseWriter,
	r *http.Request,
	length int) (bool, http.ResponseWriter, *http.Request) {
	var (
		next bool
		res  http.ResponseWriter
		req  *http.Request
	)
	for i := range middlewares {
		middlewares[length-1-i](
			w, r,
			func(w http.ResponseWriter, r *http.Request) {
				next = true
				res = w
				req = r
			})
	}
	return next, res, req
}
