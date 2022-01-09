package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *Server
	}{
		{
			name: "success",
			want: &Server{
				server: &http.Server{},
				addr:   PORT,
				routes: map[string]route{},
				ctx:    context.Background(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_router_ListenAndServe(t *testing.T) {
	tests := []struct {
		name    string
		addr    string
		path    string
		folder  string
		wantErr bool
	}{
		{
			name:    "success shutdown",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New()
			s.SetAddr(tt.addr)
			s.SetStaticPath(tt.path)
			s.SetStaticFolder(tt.folder)
			s.SetContext(context.Background())
			go func() {
				time.Sleep(1 * time.Second)
				_ = s.Shutdown()
			}()
			if err := s.ListenAndServe(); (err != nil) != tt.wantErr {
				t.Errorf("router.ListenAndServe() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_server_Get(t *testing.T) {
	t.Run("GET", func(t *testing.T) {
		r := &Server{routes: map[string]route{}}
		want := &Server{
			routes: map[string]route{
				"GET:/": {path: "/", method: "GET"},
			},
		}
		if got := r.Get("/", nil); !reflect.DeepEqual(got, want) {
			t.Errorf("httpRouter.Get() = %v, want %v", r.routes, want.routes)
		}
	})
}

func Test_server_Post(t *testing.T) {
	t.Run("POST", func(t *testing.T) {
		r := &Server{routes: map[string]route{}}
		want := &Server{
			routes: map[string]route{
				"POST:/": {path: "/", method: "POST"},
			},
		}
		if got := r.Post("/", nil); !reflect.DeepEqual(got, want) {
			t.Errorf("httpRouter.Post() = %v, want %v", r.routes, want.routes)
		}
	})
}

func Test_server_Put(t *testing.T) {
	t.Run("PUT", func(t *testing.T) {
		r := &Server{routes: map[string]route{}}
		want := &Server{
			routes: map[string]route{
				"PUT:/": {path: "/", method: "PUT"},
			},
		}
		if got := r.Put("/", nil); !reflect.DeepEqual(got, want) {
			t.Errorf("httpRouter.PUT() = %v, want %v", r.routes, want.routes)
		}
	})
}

func Test_server_ServeHTTP(t *testing.T) {
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name       string
		handler    Handler
		middleware Middleware
		args       args
		want       string
	}{
		{
			name:       "success",
			handler:    func(w http.ResponseWriter, r *http.Request) {},
			middleware: func(w http.ResponseWriter, r *http.Request, next Next) {},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "/", nil),
			},
			want: "200 OK",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New()
			s.Get("/", tt.handler, tt.middleware)
			s.ServeHTTP(tt.args.w, tt.args.r)
			got := tt.args.w.Result().Status
			if got != tt.want {
				t.Errorf("httpRouter.GET() = %v, want %v", got, tt.want)
			}
		})
	}
}
