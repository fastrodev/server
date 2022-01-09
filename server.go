package server

import (
	"context"
	"fmt"
	"net/http"
)

func New() *Server {
	return &Server{
		server: &http.Server{},
		addr:   PORT,
		routes: map[string]route{},
		ctx:    context.Background(),
	}
}

type Server struct {
	server       *http.Server
	routes       map[string]route
	ctx          context.Context
	addr         string
	staticFolder string
	staticPath   string
}

func (s *Server) ListenAndServe() error {
	fmt.Printf("Listening on http://localhost%v \n", s.addr)
	s.server.Handler = s.newHandler(false)
	s.server.Addr = s.addr
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown() error {
	return s.server.Shutdown(s.ctx)
}

func (s *Server) newHandler(serverless bool) *handler {
	return &handler{
		routes:       s.routes,
		staticPath:   s.staticPath,
		staticFolder: s.staticFolder,
		serverless:   serverless,
	}
}

func (s *Server) SetContext(ctx context.Context) *Server {
	s.ctx = ctx
	return s
}

func (s *Server) SetAddr(addr string) *Server {
	s.addr = addr
	return s
}

func (s *Server) SetStaticFolder(folder string) *Server {
	s.staticFolder = folder
	return s
}

func (s *Server) SetStaticPath(path string) *Server {
	s.staticPath = path
	return s
}

func (s *Server) Get(path string, handler Handler, middleware ...Middleware) *Server {
	key := http.MethodGet + SPLITTER + path
	s.routes[key] = route{path, http.MethodGet, handler, appendMiddleware(middleware)}
	return s
}

func (s *Server) Post(path string, handler Handler, middleware ...Middleware) *Server {
	key := http.MethodPost + SPLITTER + path
	s.routes[key] = route{path, http.MethodPost, handler, appendMiddleware(middleware)}
	return s
}

func (s *Server) Put(path string, handler Handler, middleware ...Middleware) *Server {
	key := http.MethodPut + SPLITTER + path
	s.routes[key] = route{path, http.MethodPut, handler, appendMiddleware(middleware)}
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h := s.newHandler(true)
	h.ServeHTTP(w, r)
}
