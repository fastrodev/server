package server

import "net/http"

const (
	PORT     = ":8080"
	SPLITTER = ":"
	EMPTY    = ""
	NOTFOUND = "!"
	SLASH    = "/"
	TMP      = "tmp"

	SERVERLESS_FOLDER = "serverless_function_source_code/"
)

type route struct {
	path        string
	method      string
	handler     Handler
	middlewares []Middleware
}

type Next func(w http.ResponseWriter, r *http.Request)

type Middleware func(w http.ResponseWriter, r *http.Request, next Next)
