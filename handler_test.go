package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_handler_ServeHTTP(t *testing.T) {
	not_found := "404 page not found\n"
	type fields struct {
		routes       map[string]route
		staticPath   string
		staticFolder string
		serverless   bool
	}
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "success",
			fields: fields{
				routes: map[string]route{
					"GET:/": {path: "/", method: http.MethodGet, handler: func(w http.ResponseWriter, r *http.Request) {
						fmt.Fprintf(w, "Hello")
					}},
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "/", nil),
			},
			want: "Hello",
		},
		{
			name: "success -- with params",
			fields: fields{
				routes: map[string]route{
					"GET:/:name": {path: "/:name", method: http.MethodGet, handler: func(w http.ResponseWriter, r *http.Request) {
						fmt.Fprintf(w, "Hello")
					}},
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "/agus", nil),
			},
			want: "Hello",
		},
		{
			name: "success -- with regex params",
			fields: fields{
				routes: map[string]route{
					"GET:/user/:id([0-9]+)": {path: "/user/:id([0-9]+)", method: http.MethodGet, handler: func(w http.ResponseWriter, r *http.Request) {
						fmt.Fprintf(w, "Hello")
					}},
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "/user/9", nil),
			},
			want: "Hello",
		},
		{
			name: "fail -- with regex params",
			fields: fields{
				routes: map[string]route{
					"GET:/user/:id([0-9]+)": {path: "/user/:id([0-9]+)", method: http.MethodGet, handler: func(w http.ResponseWriter, r *http.Request) {
						fmt.Fprintf(w, "Hello")
					}},
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "/user/agus", nil),
			},
			want: not_found,
		},
		{
			name: "fail -- with regex params",
			fields: fields{
				routes: map[string]route{
					"GET:/account/:id([0-9]+)": {path: "/account/:id([0-9]+)", method: http.MethodGet, handler: func(w http.ResponseWriter, r *http.Request) {
						fmt.Fprintf(w, "Hello")
					}},
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "/user/9", nil),
			},
			want: not_found,
		},
		{
			name: "not found",
			fields: fields{
				routes: map[string]route{},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "/notfound", nil),
			},
			want: not_found,
		},
		{
			name: "not found - invalid path",
			fields: fields{
				routes: map[string]route{
					"GET:/": {path: "/", method: http.MethodGet, handler: func(w http.ResponseWriter, r *http.Request) {
						fmt.Fprintf(w, "Hello")
					}},
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "/not/found", nil),
			},
			want: not_found,
		},
		{
			name: "not found - serverless",
			fields: fields{
				routes:     map[string]route{},
				serverless: true,
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "/notfound", nil),
			},
			want: not_found,
		},
		{
			name: "fail -- with url params",
			fields: fields{
				routes: map[string]route{
					"GET:/user/:name": {path: "/user/:name", method: http.MethodGet, handler: func(w http.ResponseWriter, r *http.Request) {
						fmt.Fprintf(w, "Hello")
					}},
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "/agus", nil),
			},
			want: not_found,
		},
		{
			name: "middleware with next",
			fields: fields{
				routes: map[string]route{
					"GET:/": {
						path:   "/",
						method: http.MethodGet,
						handler: func(w http.ResponseWriter, r *http.Request) {
							fmt.Fprintf(w, "Hello")
						},
						middlewares: []Middleware{
							func(w http.ResponseWriter, r *http.Request, next Next) {
								next(w, r)
							},
						},
					},
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "/", nil),
			},
			want: "Hello",
		},
		{
			name: "middleware without next",
			fields: fields{
				routes: map[string]route{
					"GET:/": {
						path:   "/",
						method: http.MethodGet,
						handler: func(w http.ResponseWriter, r *http.Request) {
							fmt.Fprintf(w, "Hello")
						},
						middlewares: []Middleware{
							func(w http.ResponseWriter, r *http.Request, next Next) {},
						},
					},
				},
			},

			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "/", nil),
			},
			want: EMPTY,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &handler{
				routes:       tt.fields.routes,
				staticPath:   tt.fields.staticPath,
				staticFolder: tt.fields.staticFolder,
				serverless:   tt.fields.serverless,
			}
			h.ServeHTTP(tt.args.w, tt.args.r)
			resp := tt.args.w.Result()
			body, _ := ioutil.ReadAll(resp.Body)
			got := string(body)
			if got != tt.want {
				t.Errorf("ServeHTTP() = %v, want %v", got, tt.want)
			}
		})
	}
}
