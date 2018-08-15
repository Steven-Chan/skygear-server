package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/skygeario/skygear-server/cmd/skynext-router/middleware"

	"github.com/gorilla/mux"
)

var routerMap map[string]*url.URL

func init() {
	auth, _ := url.Parse("http://localhost:3000")
	routerMap = map[string]*url.URL{
		"auth": auth,
	}
}

func main() {
	r := mux.NewRouter()

	r.Use(LoggingMiddleware{}.Handle)
	r.Use(middleware.APIKeyMiddleware{}.Handle)
	r.Use(middleware.ConfigMiddleware{}.Handle)

	proxy := NewReverseProxy()
	r.HandleFunc("/{module}/{rest:.*}", rewriteHandler(proxy))

	srv := &http.Server{
		Addr: "0.0.0.0:3001",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	fmt.Println("Start gateway server")
	if err := srv.ListenAndServe(); err != nil {
		log.Println(err)
	}
}

func NewReverseProxy() *httputil.ReverseProxy {
	director := func(req *http.Request) {
		path := req.URL.Path
		req.URL = routerMap[req.Header.Get("X-Skygear-Module")]
		req.URL.Path = path
	}
	return &httputil.ReverseProxy{Director: director}
}

func rewriteHandler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("X-Skygear-Module", mux.Vars(r)["module"])
		r.URL.Path = "/" + mux.Vars(r)["rest"]
		p.ServeHTTP(w, r)
	}
}

type LoggingMiddleware struct{}

func (m LoggingMiddleware) Handle(next http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, next)
}
