package middleware

import (
	"context"
	"net/http"

	"github.com/skygeario/skygear-server/cmd/skyauth/skycontext"
	"github.com/skygeario/skygear-server/pkg/server/router"
	"github.com/skygeario/skygear-server/pkg/server/skydb"
)

type DBConnMiddleware struct {
	AppName       string
	AccessControl string
	DBOpener      skydb.DBOpener
	DBImpl        string
	Option        string
	DBConfig      skydb.DBConfig
}

func (m DBConnMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessKey := r.Context().Value(skycontext.KeyAccessKey)

		dbConfig := m.DBConfig
		if accessKey == router.MasterAccessKey {
			dbConfig.CanMigrate = true
		}
		conn, err := m.DBOpener(r.Context(), m.DBImpl, m.AppName, m.AccessControl, m.Option, dbConfig)
		if err != nil {
			http.Error(w, "Unable to open database", http.StatusServiceUnavailable)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), skycontext.KeyDBConn, conn))
		next.ServeHTTP(w, r)
	})
}
