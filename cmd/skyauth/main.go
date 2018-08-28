package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/skygeario/skygear-server/cmd/skyauth/middleware"
	"github.com/skygeario/skygear-server/pkg/server/authtoken"
	"github.com/skygeario/skygear-server/pkg/server/handler"
	"github.com/skygeario/skygear-server/pkg/server/plugin/hook"
	"github.com/skygeario/skygear-server/pkg/server/router"
	"github.com/skygeario/skygear-server/pkg/server/skyconfig"
	"github.com/skygeario/skygear-server/pkg/server/skydb"
	"github.com/skygeario/skygear-server/pkg/server/skydb/pq"
)

func main() {
	pq.Init()

	r := mux.NewRouter()

	config := skyconfig.NewConfiguration()
	config.ReadFromEnv()
	if err := config.Validate(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	connOpener := ensureDB(config) // Fatal on DB failed
	initUserAuthRecordKeys(connOpener, config.App.AuthRecordKeys)

	tokenStore := authtoken.InitTokenStore(authtoken.Configuration{
		Implementation: config.TokenStore.ImplName,
		Path:           config.TokenStore.Path,
		Prefix:         config.TokenStore.Prefix,
		Expiry:         config.TokenStore.Expiry,
		Secret:         config.TokenStore.Secret,
	})

	tokenStoreMiddleware := middleware.TokenStoreMiddleware{tokenStore}
	userAuthenticator := UserAuthenticator{
		ClientKey:  config.App.APIKey,
		MasterKey:  config.App.MasterKey,
		AppName:    config.App.Name,
		TokenStore: tokenStore,
	}

	authInfoMiddleware := AuthInfoMiddleware{
		PwExpiryDays: config.UserAudit.PwExpiryDays,
	}

	dbConfig := baseDBConfig(config)
	dbConnMiddleware := DBConnMiddleware{
		AppName:       config.App.Name,
		AccessControl: config.App.AccessControl,
		DBOpener:      skydb.Open,
		DBImpl:        config.DB.ImplName,
		Option:        config.DB.Option,
		DBConfig:      dbConfig,
	}

	userMiddleware := UserMiddleware{
		hook.NewRegistry(),
		nil,
	}

	r.Use(tokenStoreMiddleware.Handle)
	r.Use(userAuthenticator.Handle)
	r.Use(dbConnMiddleware.Handle)
	r.Use(authInfoMiddleware.Handle)
	r.Use(userMiddleware.Handle)
	r.Handle("/me", MeHandler{})

	srv := &http.Server{
		Addr: "0.0.0.0:3000",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}

type MeHandler struct {
}

func (m MeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	info := r.Context().Value(KeyAuthInfo).(*skydb.AuthInfo)
	if info == nil {
		http.Error(w, "Authentication is needed to get current user", http.StatusUnauthorized)
		return
	}

	store := r.Context().Value(KeyTokenStore).(authtoken.Store)
	appName := r.Context().Value(KeyAppName).(string)
	token, err := store.NewToken(appName, info.ID)
	if err != nil {
		panic(err)
	}

	dbConn := r.Context().Value(KeyDBConn).(skydb.Conn)
	db := dbConn.PublicDB()
	user := r.Context().Value(KeyUser).(*skydb.Record)
	accessKey := r.Context().Value(KeyAccessKey)

	// We will return the last seen in DB, not current time stamp
	authResponse, err := handler.AuthResponseFactory{
		AssetStore: nil,
		Conn:       dbConn,
		Database:   db,
	}.NewAuthResponse(*info, *user, token.AccessToken, accessKey == router.MasterAccessKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Populate the activity time to user
	now := time.Now().UTC()
	info.LastSeenAt = &now
	if err := dbConn.UpdateAuth(info); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Printf("%+v", authResponse)
	// fmt.Fprintf(w, "%+v", authResponse)

	w.Header().Set("Content-Type", "applicaiton/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(authResponse)
}
