package main

import (
	"net/http"

	"github.com/facebookgo/inject"
	"github.com/gorilla/handlers"
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
	serveMux := http.NewServeMux()

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

	baseMiddleware := middleware.ChainMiddleware(
		handlers.RecoveryHandler(),
		userRecordMiddleware.Handle,
	)
	serveMux.Handle("/me", baseMiddleware(meHandler))

	authMiddleware := middleware.ChainMiddleware(
		handlers.RecoveryHandler(),
	)
	serveMux.Handle("/auth/login", authMiddleware(loginHandler))

	http.ListenAndServe(":3000", serveMux)
}
