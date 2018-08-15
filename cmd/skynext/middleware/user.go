package middleware

import (
	"net/http"

	"github.com/skygeario/skygear-server/cmd/skynext/model"
	"github.com/skygeario/skygear-server/cmd/skynext/service"
)

type UserRecordMiddleware struct {
	UserStore service.UserStore `inject:"UserStore"`
}

func (m UserRecordMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username := model.GetUsername(r)

		if username == "Kuen kuen" {
			panic("bomb")
		}

		if user := m.UserStore.Get(username); user != nil {
			model.SetUser(r, user)
		}

		next.ServeHTTP(w, r)
	})
}
