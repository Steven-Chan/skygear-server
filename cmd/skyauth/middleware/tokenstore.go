package middleware

import (
	"context"
	"net/http"

	"github.com/skygeario/skygear-server/cmd/skyauth/skycontext"
	"github.com/skygeario/skygear-server/pkg/server/authtoken"
)

type TokenStoreMiddleware struct {
	tokenStore authtoken.Store
}

func (m TokenStoreMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), skycontext.KeyTokenStore, m.tokenStore)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
