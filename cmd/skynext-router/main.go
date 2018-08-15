package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/handlers"

	"github.com/gorilla/mux"
)

type Key struct {
	APIKey    string
	MasterKey string
}

var keyMap map[string]Key
var routerMap map[string]*url.URL

func init() {
	keyMap = map[string]Key{
		"skygear": Key{
			APIKey:    "apikey",
			MasterKey: "masterkey",
		},
		"skygear-next": Key{
			APIKey:    "apikey-next",
			MasterKey: "masterkey-next",
		},
	}
	auth, _ := url.Parse("http://localhost:3000")
	routerMap = map[string]*url.URL{
		"auth": auth,
	}
}

func main() {
	r := mux.NewRouter()

	r.Use(LoggingMiddleware{}.Handle)
	r.Use(APIKeyMiddleware{}.Handle)

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

func keyMapLookUp(apiKey string) (KeyType, string) {
	for appName, key := range keyMap {
		if key.APIKey == key.APIKey {
			return APIAccessKey, appName
		}

		if key.MasterKey == key.MasterKey {
			return MasterAccessKey, appName
		}
	}

	return NoAccessKey, ""
}

type LoggingMiddleware struct{}

func (m LoggingMiddleware) Handle(next http.Handler) http.Handler {
	return handlers.LoggingHandler(os.Stdout, next)
}

type APIKeyMiddleware struct {
}

func (a APIKeyMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := GetAPIKey(r)
		keyType, appName := keyMapLookUp(apiKey)
		if keyType == NoAccessKey {
			http.Error(w, "API key not set", http.StatusBadRequest)
			return
		}

		SetAccessKeyType(r, keyType)
		SetAppName(r, appName)
		next.ServeHTTP(w, r)
	})
}

type KeyType int

const (
	NoAccessKey KeyType = iota
	APIAccessKey
	MasterAccessKey
)

func header(i interface{}) http.Header {
	switch i.(type) {
	case *http.Request:
		return (i.(*http.Request)).Header
	case http.ResponseWriter:
		return (i.(http.ResponseWriter)).Header()
	default:
		panic("Invalid type")
	}
}

func GetAccessKeyType(i interface{}) KeyType {
	i, err := strconv.Atoi(header(i).Get("X-Skygear-AccessKeyType"))
	if err != nil {
		return NoAccessKey
	}

	kt, ok := i.(KeyType)
	if !ok {
		return NoAccessKey
	}

	return kt
}

func SetAccessKeyType(i interface{}, kt KeyType) {
	header(i).Set("X-Skygear-AccessKeyType", strconv.Itoa(int(kt)))
}

func GetAPIKey(i interface{}) string {
	return header(i).Get("X-Skygear-APIKey")
}

func SetAppName(i interface{}, appName string) {
	header(i).Set("X-Skygear-AppName", appName)
}
