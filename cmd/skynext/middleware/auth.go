package middleware

import (
	"fmt"
	"net/http"

	rModel "github.com/skygeario/skygear-server/cmd/skynext-router/model"
	"github.com/skygeario/skygear-server/cmd/skynext/model"
	"github.com/skygeario/skygear-server/pkg/server/authtoken"
)

type AuthInfoMiddleware struct {
	TokenStore authtoken.Store `inject:"TokenStore"`
}

func (m AuthInfoMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessKeyType := rModel.GetAccessKeyType(r)
		tokenString := model.GetAccessToken(r)

		fmt.Println("tokenString: ", tokenString)

		// If payload contains an access token, check whether if the access
		// token is valid. API Key is not required if there is valid access token.
		if tokenString != "" {
			store := m.TokenStore
			token := authtoken.Token{}

			if err := store.Get(tokenString, &token); err == nil {
				model.SetAuthInfoID(r, token.AuthInfoID)
			}
		} else if accessKeyType == rModel.MasterAccessKey {
			// For master access key, it is possible to impersonate any user of
			// the caller's choosing.
			if userID := model.GetUserID(r); userID != "" {
				model.SetAuthInfoID(r, userID)
			}
		}

		// TODO: Inject auth info here?

		next.ServeHTTP(w, r)
	})
}
