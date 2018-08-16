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
	"github.com/skygeario/skygear-server/pkg/server/authtoken"
	"github.com/skygeario/skygear-server/pkg/server/skyconfig"
	"github.com/skygeario/skygear-server/pkg/server/skydb"
	"github.com/skygeario/skygear-server/pkg/server/skydb/pq"
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
	env := os.Getenv("GO_ENV")

	pq.WarmUp()

	// TODO: Read config from request header
	config := skyconfig.NewConfiguration()
	config.ReadFromEnv()
	if err := config.Validate(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	dbConfig := baseDBConfig(config)

	tokenStore := authtoken.InitTokenStore(authtoken.Configuration{
		Implementation: config.TokenStore.ImplName,
		Path:           config.TokenStore.Path,
		Prefix:         config.TokenStore.Prefix,
		Expiry:         config.TokenStore.Expiry,
		Secret:         config.TokenStore.Secret,
	})

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
		&inject.Object{
			Value: service.SkygearDatabaseProvider{
				AppName:       config.App.Name,
				AccessControl: config.App.AccessControl,
				DBOpener:      skydb.Open,
				DBImpl:        config.DB.ImplName,
				Option:        config.DB.Option,
				DBConfig:      dbConfig,
			},
			Complete: true,
			Name:     "DatabaseProvider",
		},
		&inject.Object{
			Value:    tokenStore,
			Complete: true,
			Name:     "TokenStore",
		},
	)

	// prepare middleware
	authMiddleware := middleware.AuthInfoMiddleware{}
	userRecordMiddleware := middleware.UserRecordMiddleware{}
	var configMiddleware middleware.Middleware
	if env == "dev" {
		configMiddleware = middleware.ConfigMiddleware{}
	} else {
		configMiddleware = middleware.NoOpMiddleware{}
	}

	injector.Provide(
		&inject.Object{Value: &authMiddleware},
		&inject.Object{Value: &userRecordMiddleware},
		&inject.Object{Value: &configMiddleware},
	)

	// prepare handlers
	meHandler := handler.MeHandler{}
	loginHandler := handler.LoginHandler{}
	configHandler := handler.ConfigHandler{}
	skySchemaFetchHandler := handler.SchemaFetchHandler{}

	injector.Provide(
		&inject.Object{Value: &meHandler},
		&inject.Object{Value: &loginHandler},
		&inject.Object{Value: &configHandler},
		&inject.Object{Value: &skySchemaFetchHandler},
	)
	if err := injector.Populate(); err != nil {
		panic(err)
	}

	// transform skygear handler to skygear next handler
	schemaFetchHandler := handler.FromSkygearHandler(skySchemaFetchHandler)

	switch os.Getenv("IMPL") {
	case "":
		fallthrough
	case "net/http":
		serveMux := http.NewServeMux()

		baseMiddleware := middleware.ChainMiddleware(
			LoggingMiddleware{}.Handle,
			handlers.RecoveryHandler(handlers.PrintRecoveryStack(true)),
			configMiddleware.Handle,
			authMiddleware.Handle,
			userRecordMiddleware.Handle,
		)
		serveMux.Handle("/me", baseMiddleware(meHandler))
		serveMux.Handle("/config", baseMiddleware(configHandler))
		serveMux.Handle("/schema/fetch", baseMiddleware(schemaFetchHandler))

		authMiddleware := middleware.ChainMiddleware(
			LoggingMiddleware{}.Handle,
			handlers.RecoveryHandler(handlers.PrintRecoveryStack(true)),
			configMiddleware.Handle,
		)
		serveMux.Handle("/auth/login", authMiddleware(loginHandler))

		fmt.Println("Start net/http server")
		if err := http.ListenAndServe(":3000", serveMux); err != nil {
			log.Println(err)
		}
	case "gorilla":
		r := mux.NewRouter()

		r.Use(LoggingMiddleware{}.Handle)
		r.Use(handlers.RecoveryHandler(handlers.PrintRecoveryStack(true)))
		r.Use(authMiddleware.Handle)
		r.Use(configMiddleware.Handle)
		r.Handle("/me", userRecordMiddleware.Handle(meHandler))
		r.Handle("/config", configHandler)
		r.Handle("/schema/fetch", schemaFetchHandler)

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
			handlers.RecoveryHandler(handlers.PrintRecoveryStack(true)),
			configMiddleware.Handle,
			authMiddleware.Handle,
			userRecordMiddleware.Handle,
		)
		r.POST("/me", wrapHandler(baseMiddleware(meHandler)))
		r.POST("/config", wrapHandler(baseMiddleware(configHandler)))
		r.POST("/schema/fetch", wrapHandler(baseMiddleware(schemaFetchHandler)))

		authMiddleware := middleware.ChainMiddleware(
			LoggingMiddleware{}.Handle,
			handlers.RecoveryHandler(handlers.PrintRecoveryStack(true)),
			configMiddleware.Handle,
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

// TODO: Remove all configuration related code from this file
func baseDBConfig(config skyconfig.Configuration) skydb.DBConfig {
	passwordHistoryEnabled := config.UserAudit.PwHistorySize > 0 ||
		config.UserAudit.PwHistoryDays > 0

	return skydb.DBConfig{
		CanMigrate:             config.App.DevMode,
		PasswordHistoryEnabled: passwordHistoryEnabled,
	}
}
