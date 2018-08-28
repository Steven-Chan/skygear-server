package middleware

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/skygeario/skygear-server/cmd/skyauth/skycontext"
	"github.com/skygeario/skygear-server/pkg/server/asset"
	"github.com/skygeario/skygear-server/pkg/server/authtoken"
	"github.com/skygeario/skygear-server/pkg/server/plugin/hook"
	"github.com/skygeario/skygear-server/pkg/server/recordutil"
	"github.com/skygeario/skygear-server/pkg/server/router"
	"github.com/skygeario/skygear-server/pkg/server/skydb"
	"github.com/skygeario/skygear-server/pkg/server/skyerr"
)

type UserAuthenticator struct {
	ClientKey  string
	MasterKey  string
	AppName    string
	TokenStore authtoken.Store
}

func checkRequestAccessKey(r *http.Request, clientKey string, masterKey string) skyerr.Error {
	var accessKey router.AccessKeyType

	log.Printf("apikey: %+v", r.Header["X-Skygear-Api-Key"])

	apiKey := r.Header["X-Skygear-Api-Key"][0]
	if masterKey != "" && apiKey == masterKey {
		accessKey = router.MasterAccessKey
	} else if clientKey != "" && apiKey == clientKey {
		accessKey = router.ClientAccessKey
	} else if apiKey == "" {
		accessKey = router.NoAccessKey
	} else {
		return skyerr.NewErrorf(skyerr.AccessKeyNotAccepted, "Cannot verify api key: `%v`", apiKey)
	}

	r = r.WithContext(context.WithValue(r.Context(), skycontext.KeyAccessKey, accessKey))
	return nil
}

func (p UserAuthenticator) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := checkRequestAccessKey(r, p.ClientKey, p.MasterKey); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// If payload contains an access token, check whether if the access
		// token is valid. API Key is not required if there is valid access token.
		if tokenString := r.Header["X-Skygear-Access-Token"][0]; tokenString != "" {
			store := p.TokenStore
			token := authtoken.Token{}

			if err := store.Get(tokenString, &token); err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), skycontext.KeyAppName, token.AppName))
			r = r.WithContext(context.WithValue(r.Context(), skycontext.KeyAuthInfoID, token.AuthInfoID))
			r = r.WithContext(context.WithValue(r.Context(), skycontext.KeyAccessToken, token))
			r = r.WithContext(context.WithValue(r.Context(), skycontext.KeyUserIDContext, token.AuthInfoID))

			next.ServeHTTP(w, r)
			return
		}

		accessKey := r.Context().Value(skycontext.KeyAccessKey)
		if accessKey == router.NoAccessKey {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// For master access key, it is possible to impersonate any user of
		// the caller's choosing.
		if accessKey == router.MasterAccessKey {
			if userID := r.Header["X-Skygear-User-ID"][0]; userID != "" {
				r = r.WithContext(context.WithValue(r.Context(), skycontext.KeyAuthInfoID, userID))
			}
		}

		r = r.WithContext(context.WithValue(r.Context(), skycontext.KeyAppName, p.AppName))
		next.ServeHTTP(w, r)
	})
}

type AuthInfoMiddleware struct {
	PwExpiryDays int
}

func isTokenStillValid(token router.AccessToken, authInfo skydb.AuthInfo) bool {
	if authInfo.TokenValidSince == nil {
		return true
	}
	tokenValidSince := *authInfo.TokenValidSince

	// Not all types of access token support this field. The token is
	// still considered if it does not have an issue time.
	if token.IssuedAt().IsZero() {
		return true
	}

	// Due to precision, the issue time of the token can be before
	// AuthInfo.TokenValidSince. We consider the token still valid
	// if the token is issued within 1 second before tokenValidSince.
	return token.IssuedAt().After(tokenValidSince.Add(-1 * time.Second))
}

func (p AuthInfoMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authInfoID := r.Context().Value(skycontext.KeyAuthInfoID).(string)
		accessKey := r.Context().Value(skycontext.KeyAccessKey)

		// If the payload does not already have auth_id, and if the request
		// is authenticated with master key, assume the user is _god.
		if authInfoID == "" && accessKey == router.MasterAccessKey {
			authInfoID = "_god"
			r = r.WithContext(context.WithValue(r.Context(), skycontext.KeyUserIDContext, authInfoID))
		}

		authinfo := skydb.AuthInfo{}

		// Query database to get auth info
		// If an error occurred at this stage, Internal Server Error is returned.
		if authInfoID == "" {
			return
		}

		_, status := p.fetchOrCreateAuth(r, &authinfo)
		if status != 0 {
			return
		}

		accessToken := r.Context().Value(skycontext.KeyAccessToken).(router.AccessToken)

		// If an access token exists checks if the access token has an IssuedAt
		// time that is later than the user's TokenValidSince time. This
		// allows user to invalidate previously issued access token.
		if accessToken != nil && !isTokenStillValid(accessToken, authinfo) {
			http.Error(w, "token does not exist or it has expired", http.StatusUnauthorized)
			return
		}

		// Check if password is expired according to policy
		if authinfo.IsPasswordExpired(p.PwExpiryDays) {
			http.Error(w, "password expired", http.StatusUnauthorized)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), skycontext.KeyAuthInfo, &authinfo))
		next.ServeHTTP(w, r)
		return
	})
}

func (p AuthInfoMiddleware) fetchOrCreateAuth(r *http.Request, authInfo *skydb.AuthInfo) (skyerr.Error, int) {
	dbConn := r.Context().Value(skycontext.KeyDBConn).(skydb.Conn)
	authInfoID := r.Context().Value(skycontext.KeyAuthInfoID).(string)
	accessKey := r.Context().Value(skycontext.KeyAccessKey)

	var err error
	err = dbConn.GetAuth(authInfoID, authInfo)
	if err == skydb.ErrUserNotFound && accessKey == router.MasterAccessKey {
		*authInfo = skydb.AuthInfo{
			ID: authInfoID,
		}
		err = dbConn.CreateAuth(authInfo)
		if err == skydb.ErrUserDuplicated {
			// user already exists, error can be ignored
			err = nil
		}
	}

	if err != nil {
		return skyerr.NewError(skyerr.UnexpectedAuthInfoNotFound, err.Error()), http.StatusInternalServerError
	}

	return nil, 0
}

// UserMiddleware injects a user record to the payload
//
// An AuthInfo must be injected before this, if it is not found, the preprocessor
// would just skip the injection
//
// If AuthInfo is injected but a user record is not found, the preprocessor would
// create a new user record and inject it to the payload
type UserMiddleware struct {
	HookRegistry *hook.Registry
	AssetStore   asset.Store
}

func (p UserMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dbConn := r.Context().Value(skycontext.KeyDBConn).(skydb.Conn)
		db := dbConn.PublicDB()

		user, _ := r.Context().Value(skycontext.KeyUser).(*skydb.Record)
		authInfo, _ := r.Context().Value(skycontext.KeyAuthInfo).(*skydb.AuthInfo)

		if user == nil && authInfo != nil {
			user := skydb.Record{}
			err := db.Get(skydb.NewRecordID("user", authInfo.ID), &user)

			if err == skydb.ErrRecordNotFound {
				user, err = p.createUser(r.Context(), dbConn, authInfo)
			}

			if err != nil {
				http.Error(w, "User not found", http.StatusInternalServerError)
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), skycontext.KeyUser, &user))
		}

		next.ServeHTTP(w, r)
	})
}

func (p UserMiddleware) createUser(ctx context.Context, dbConn skydb.Conn, authInfo *skydb.AuthInfo) (skydb.Record, error) {
	db := dbConn.PublicDB()
	txDB, ok := db.(skydb.Transactional)
	if !ok {
		return skydb.Record{}, skyerr.NewError(skyerr.NotSupported, "database impl does not support transaction")
	}

	var user *skydb.Record
	txErr := skydb.WithTransaction(txDB, func() error {
		userRecord := skydb.Record{
			ID: skydb.NewRecordID(db.UserRecordType(), authInfo.ID),
		}

		recordReq := recordutil.RecordModifyRequest{
			Db:           db,
			Conn:         dbConn,
			AssetStore:   p.AssetStore,
			HookRegistry: p.HookRegistry,
			Atomic:       true,
			Context:      ctx,
			AuthInfo:     authInfo,
			ModifyAt:     time.Now().UTC(),
			RecordsToSave: []*skydb.Record{
				&userRecord,
			},
		}

		recordResp := recordutil.RecordModifyResponse{
			ErrMap: map[skydb.RecordID]skyerr.Error{},
		}

		err := recordutil.RecordSaveHandler(&recordReq, &recordResp)
		if err != nil {
			return err
		}

		user = recordResp.SavedRecords[0]
		return nil
	})

	if txErr != nil {
		return skydb.Record{}, txErr
	}

	return *user, nil
}
