package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/facebookgo/inject"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/julienschmidt/httprouter"
	"github.com/skygeario/skygear-server/cmd/skynext/handler"
	"github.com/skygeario/skygear-server/cmd/skynext/middleware"
	"github.com/skygeario/skygear-server/cmd/skynext/model"
	"github.com/skygeario/skygear-server/cmd/skynext/service"
)

var passwordStore map[string]string
var userMap map[string]model.User

func init() {
	userMap = map[string]model.User{
		"god": model.User{Username: "god", Age: 100},
		"ten": model.User{Username: "ten", Age: 18},
	}

	passwordStore = map[string]string{
		"ten": "i am your father",
	}
}

func main() {
	injector := inject.Graph{}
	injector.Provide(
		&inject.Object{
			Value:    service.MemoryUserAuthenticator{Store: passwordStore},
			Complete: true,
			Name:     "UserAuthenticator",
		},
		&inject.Object{
			Value:    service.MemoryUserStore{Store: userMap},
			Complete: true,
			Name:     "UserStore",
		},
	)

	// prepare middleware
	userRecordMiddleware := middleware.UserRecordMiddleware{}
	injector.Provide(
		&inject.Object{Value: &userRecordMiddleware},
	)

	// prepare handlers
	meHandler := handler.MeHandler{}
	loginHandler := handler.LoginHandler{}

	injector.Provide(
		&inject.Object{Value: &meHandler},
		&inject.Object{Value: &loginHandler},
	)
	if err := injector.Populate(); err != nil {
		panic(err)
	}

	switch os.Getenv("IMPL") {
	case "":
		fallthrough
	case "net/http":
		serveMux := http.NewServeMux()

		baseMiddleware := middleware.ChainMiddleware(
			handlers.RecoveryHandler(),
			userRecordMiddleware.Handle,
		)
		serveMux.Handle("/me", baseMiddleware(meHandler))

		authMiddleware := middleware.ChainMiddleware(
			handlers.RecoveryHandler(),
		)
		serveMux.Handle("/auth/login", authMiddleware(loginHandler))

		fmt.Println("Start net/http server")
		if err := http.ListenAndServe(":3000", serveMux); err != nil {
			log.Println(err)
		}
	case "gorilla":
		r := mux.NewRouter()

		r.Use(handlers.RecoveryHandler())
		r.Handle("/me", userRecordMiddleware.Handle(meHandler))

		authR := r.PathPrefix("/auth").Subrouter()
		authR.Handle("/login", loginHandler)

		/*
			userR := r.PathPrefix("/user").Subrouter()
			userR.Use(userRecordMiddleware.Handle)
			userR.Handle("/me", meHandler)

			authR := r.PathPrefix("/auth").Subrouter()
			authR.Handle("/login", loginHandler)
		*/

		srv := &http.Server{
			Addr: "0.0.0.0:3000",
			// Good practice to set timeouts to avoid Slowloris attacks.
			WriteTimeout: time.Second * 15,
			ReadTimeout:  time.Second * 15,
			IdleTimeout:  time.Second * 60,
			Handler:      r, // Pass our instance of gorilla/mux in.
		}

		fmt.Println("Start gorilla mux server")
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	case "httprouter":
		r := httprouter.New()

		baseMiddleware := middleware.ChainMiddleware(
			LoggingMiddleware{}.Handle,
			handlers.RecoveryHandler(),
			userRecordMiddleware.Handle,
		)
		r.POST("/me", wrapHandler(baseMiddleware(meHandler)))

		authMiddleware := middleware.ChainMiddleware(
			LoggingMiddleware{}.Handle,
			handlers.RecoveryHandler(),
		)
		r.POST("/auth/login", wrapHandler(authMiddleware(loginHandler)))

		r.POST("/hello/:name", wrapMiddleware(baseMiddleware)(Hello))

		fmt.Println("Start httprouter server")
		if err := http.ListenAndServe(":3000", r); err != nil {
			log.Println(err)
		}
	}
}

func wrapHandler(h http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		h.ServeHTTP(w, r)
	}
}

type httprouterMiddlewareFunc func(httprouter.Handle) httprouter.Handle

func wrapMiddleware(h middleware.MiddlewareFunc) httprouterMiddlewareFunc {
	return func(next httprouter.Handle) httprouter.Handle {
		return func(rw http.ResponseWriter, rr *http.Request, ps httprouter.Params) {
			h(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next(w, r, ps)
			})).ServeHTTP(rw, rr)
		}
	}
}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	user := model.GetUser(r)
	if user == nil {
		http.Error(w, "Current user not found", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "Hello %s from %s!\n", ps.ByName("name"), user.Username)
}

type LoggingMiddleware struct{}

func (m LoggingMiddleware) Handle(next http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, next)
}
