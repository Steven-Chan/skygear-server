package ssohandler

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/skygeario/skygear-server/pkg/auth/dependency/provider/oauth"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/sso"
	"github.com/skygeario/skygear-server/pkg/core/auth"
	"github.com/skygeario/skygear-server/pkg/core/auth/authinfo"
	"github.com/skygeario/skygear-server/pkg/core/auth/authtoken"
	"github.com/skygeario/skygear-server/pkg/core/auth/role"
	"github.com/skygeario/skygear-server/pkg/core/db"
	"github.com/skygeario/skygear-server/pkg/server/skyerr"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAuthPayload(t *testing.T) {
	Convey("Test AuthRequestPayload", t, func() {
		// callback URL and ux_mode is required
		Convey("validate valid payload", func() {
			payload := AuthRequestPayload{
				Code:         "code",
				EncodedState: "state",
			}
			So(payload.Validate(), ShouldBeNil)
		})

		Convey("validate payload without code", func() {
			payload := AuthRequestPayload{
				EncodedState: "state",
			}
			err := payload.Validate()
			errResponse := err.(skyerr.Error)
			So(errResponse.Code(), ShouldEqual, skyerr.InvalidArgument)
		})

		Convey("validate payload without state", func() {
			payload := AuthRequestPayload{
				Code: "code",
			}
			err := payload.Validate()
			errResponse := err.(skyerr.Error)
			So(errResponse.Code(), ShouldEqual, skyerr.InvalidArgument)
		})
	})
}

func TestAuthHandler(t *testing.T) {
	Convey("Test TestAuthURLHandler", t, func() {
		Convey("when user is not existed", func() {
			code := "code"
			action := "login"
			UXMode := "web_redirect"
			stateJWTSecret := "secret"

			sh := &AuthHandler{}
			sh.TxContext = db.NewMockTxContext()
			sh.AuthContext = auth.NewMockContextGetterWithDefaultUser()
			setting := sso.Setting{
				URLPrefix:      "http://localhost:3000",
				StateJWTSecret: stateJWTSecret,
			}
			config := sso.Config{
				Name:         "mock",
				ClientID:     "mock_client_id",
				ClientSecret: "mock_client_secret",
			}
			mockProvider := sso.MockSSOProverImpl{
				BaseURL: "http://mock/auth",
				Setting: setting,
				Config:  config,
			}
			sh.Provider = &mockProvider
			mockOAuthProvider := oauth.NewMockProvider(
				map[string]oauth.Principal{},
			)
			sh.OAuthAuthProvider = mockOAuthProvider
			authInfoStore := authinfo.NewMockStoreWithAuthInfoMap(
				map[string]authinfo.AuthInfo{},
			)
			sh.AuthInfoStore = authInfoStore
			mockTokenStore := authtoken.NewJWTStore("myApp", "secret", 0)
			sh.TokenStore = mockTokenStore
			sh.RoleStore = role.NewMockStore()
			h := sh.APIHandler()

			state := sso.State{
				CallbackURL: "callbackURL",
				UXMode:      UXMode,
				Action:      action,
			}
			encodedState, _ := sso.EncodeState(stateJWTSecret, state)

			Convey("should return login auth response", func() {
				v := url.Values{}
				v.Set("code", code)
				v.Add("state", encodedState)
				u := url.URL{
					RawQuery: v.Encode(),
				}

				req, _ := http.NewRequest("GET", u.RequestURI(), nil)
				resp := httptest.NewRecorder()

				h.ServeHTTP(resp, req)
				// for web_redirect, it should redirect to original callback url
				So(resp.Code, ShouldEqual, 302)
			})
		})
	})
}
