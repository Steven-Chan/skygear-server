package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	server := NewServer(&RealDBProvider{})

	server.Handle("/", &LoginHandlerFactory{}).Methods("POST")

	srv := &http.Server{
		Addr: "0.0.0.0:3000",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      server.Router, // Pass our instance of gorilla/mux in.
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Println(err)
	}
}

type TenantConfiguration struct {
	DBConnectionStr string
}

type HandlerContext struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	Configuration  TenantConfiguration
}

type Handler interface {
	Handle(HandlerContext)
}

type HandlerFactory interface {
	NewHandler() Handler
}

/*
 * Dependency Provider
 */
type DBProvider interface {
	GetDB(TenantConfiguration) IDB
}

func NewGet(dependencyName string, tConfig TenantConfiguration, dbProvider DBProvider) GetDB {
	switch dependencyName {
	case "DB":
		return func() IDB {
			return dbProvider.GetDB(tConfig)
		}
	default:
		return nil
	}
}

type RealDBProvider struct{}

func (p RealDBProvider) GetDB(tConfig TenantConfiguration) IDB {
	return &DB{tConfig.DBConnectionStr}
}

/*
 * Dependency
 */
type IDB interface {
	GetRecord(string) string
}

type DB struct {
	ConnectionStr string
}

func (db DB) GetRecord(recordID string) string {
	return db.ConnectionStr + ":" + recordID
}

/*
 * Handler factory (multi-tenant)
 */
type LoginHandlerFactory struct{}

func (f LoginHandlerFactory) NewHandler() Handler {
	return &LoginHandler{}
}

type Server struct {
	*mux.Router

	DBProvider DBProvider
}

func NewServer(dbProvider DBProvider) Server {
	return Server{
		Router:     mux.NewRouter(),
		DBProvider: dbProvider,
	}
}

func (s *Server) Handle(path string, hf HandlerFactory) *mux.Route {
	return s.NewRoute().Path(path).Handler(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		configuration := TenantConfiguration{"public"}
		handler := hf.NewHandler()

		t := reflect.TypeOf(handler).Elem()
		v := reflect.ValueOf(handler).Elem()

		numField := t.NumField()
		for i := 0; i < numField; i++ {
			dependency := t.Field(i).Tag.Get("dependency")
			field := v.Field(i)
			field.Set(reflect.ValueOf(NewGet(dependency, configuration, s.DBProvider)))
		}

		handler.Handle(HandlerContext{rw, r, TenantConfiguration{}})
	}))
}

type GetDB func() IDB

type LoginHandler struct {
	GetDB `dependency:"DB"`
}

func (h LoginHandler) Handle(ctx HandlerContext) {
	input, _ := ioutil.ReadAll(ctx.Request.Body)
	fmt.Fprintln(ctx.ResponseWriter, `{"user": "`+h.GetDB().GetRecord("user:"+string(input))+`"}`)
}
