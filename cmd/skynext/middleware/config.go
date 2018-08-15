package middleware

import (
	"net/http"

	rModel "github.com/skygeario/skygear-server/cmd/skynext-router/model"
)

type ConfigMiddleware struct {
}

func (a ConfigMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		config := rModel.Config{}
		config.ReadFromEnv()
		rModel.SetConfig(r, config)
		next.ServeHTTP(w, r)
	})
}
