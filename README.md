# Server
[![][build]](https://github.com/fastrodev/server/actions/workflows/build.yml) [![Coverage Status][cov]](https://coveralls.io/github/fastrodev/server?branch=main) [![][reference]](https://pkg.go.dev/github.com/fastrodev/server?tab=doc)

A simple and idiomatic golang framework for handling HTTP. It uses standard golang requests and responses.
## Get started

Init folder
```
mkdir app && cd app
```
Init app
```
go mod init app
```
Install
```
go get github.com/fastrodev/server
```
Create main.go file:
```go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/fastrodev/server"
)

func main() {
	s := server.New()
	s.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Helo World")
	})
	log.Panic(s.ListenAndServe())
}
```
## Handler
You can use golang standard `http.ResponseWriter` and `http.Request` for the handler signature.
```go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/fastrodev/server"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Helo World")
}

func main() {
	s := server.New()
	s.Get("/", handler)
	log.Panic(s.ListenAndServe())
}
```
## Middleware
You can access all properties and methods of `http.ResponseWriter` and `r *http.Request` before they are handled by the handler.
```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/fastrodev/server"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Helo World")
}

func middleware(w http.ResponseWriter, r *http.Request, next server.Next) {
	dt := time.Now()
	fmt.Println(dt)
	next(w, r)
}

func main() {
	s := server.New()
	s.Get("/", handler, middleware)
	log.Panic(s.ListenAndServe())
}
```

## Routing
You can add multiple handlers for a path with different methods. You can also use common url parameters or with regex.
```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/fastrodev/server"
)

func main() {
	s := server.New()
	s.Get("/api/user", getAllUser)
	s.Get("/api/user/:name", getUserHandler)
	s.Get("/api/account/:id([0-9]+)", getUserByIDHandler)
	s.Post("/api/user", postHandler)
	s.Put("/api/user", putHandler)
	log.Panic(s.ListenAndServe())
}

func getAllUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Helo World")
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, r.URL.Path)
}

func getUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, r.URL.Path)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	type User struct {
		Name string
		Age  int
	}
	var u User

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "User: %+v", u)
}

var putHandler = func(w http.ResponseWriter, r *http.Request) {
	type User struct {
		Name string
		Age  int
	}
	var u User

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "User: %+v", u)
}

```

## Serverless
You can deploy your codes to [google cloud function](https://cloud.google.com/functions). With this approach, you don't call the `ListenAndServe` function again. 
```go
package function

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fastrodev/server"
)

func Main(w http.ResponseWriter, r *http.Request) {
	createServer().ServeHTTP(w, r)
}

func createServer() *server.Server {
	s := server.New()
	s.Get("/api/user", getAllUser)
	s.Get("/api/user/:name", getUserHandler)
	s.Get("/api/account/:id([0-9]+)", getUserByIDHandler)
	s.Post("/api/user", postHandler)
	s.Put("/api/user", putHandler)
	return s
}

func getAllUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Helo World")
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, r.URL.Path)
}

func getUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, r.URL.Path)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	type User struct {
		Name string
		Age  int
	}
	var u User

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "User: %+v", u)
}

var putHandler = func(w http.ResponseWriter, r *http.Request) {
	type User struct {
		Name string
		Age  int
	}
	var u User

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "User: %+v", u)
}

```

## Contributing
We appreciate your help! The main purpose of this repository is to improve performance and readability, making it faster and easier to use.


[build]: https://github.com/fastrodev/server/actions/workflows/build.yml/badge.svg
[reference]: https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white "reference"
[cov]: https://coveralls.io/repos/github/fastrodev/server/badge.svg?branch=main