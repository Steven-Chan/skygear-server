package service

import (
	"net/http"
	"strings"

	rModel "github.com/skygeario/skygear-server/cmd/skynext-router/model"
	"github.com/skygeario/skygear-server/cmd/skynext/model"
	"github.com/skygeario/skygear-server/pkg/server/skydb"
	"github.com/skygeario/skygear-server/pkg/server/skyerr"
)

type DatabaseProvider interface {
	GetDatabase(r *http.Request) (skydb.Database, error)
}

type SkygearDatabaseProvider struct {
	AppName       string
	AccessControl string
	DBOpener      skydb.DBOpener
	DBImpl        string
	Option        string
	DBConfig      skydb.DBConfig
}

func (p SkygearDatabaseProvider) GetDBConn(r *http.Request) (skydb.Conn, error) {
	dbConfig := p.DBConfig
	accessKeyType := rModel.GetAccessKeyType(r)
	if accessKeyType == rModel.MasterAccessKey {
		dbConfig.CanMigrate = true
	}

	conn, err := p.DBOpener(r.Context(), p.DBImpl, p.AppName, p.AccessControl, p.Option, dbConfig)
	return conn, err
}

func (p SkygearDatabaseProvider) GetDatabase(r *http.Request) (skydb.Database, error) {
	conn, err := p.GetDBConn(r)
	if err != nil {
		return nil, err
	}

	databaseID := model.GetDatabaseID(r)
	if databaseID == "" {
		databaseID = "_public"
	}

	accessKeyType := rModel.GetAccessKeyType(r)
	authInfoID := model.GetAuthInfoID(r)
	var database skydb.Database

	switch databaseID {
	case "_private":
		if authInfoID != "" {
			database = conn.PrivateDB(authInfoID)
		} else {
			return nil, skyerr.NewError(skyerr.NotAuthenticated, "Authentication is needed for private DB access")
		}
	case "_public":
		database = conn.PublicDB()
	case "_union":
		if accessKeyType != rModel.MasterAccessKey {
			return nil, skyerr.NewError(skyerr.NotAuthenticated, "Master key is needed for union DB access")
		}
		database = conn.UnionDB()
	default:
		if strings.HasPrefix(databaseID, "_") {
			return nil, skyerr.NewInvalidArgument("invalid database ID", []string{"database_id"})
		} else if accessKeyType == rModel.MasterAccessKey {
			database = conn.PrivateDB(databaseID)
		} else if authInfoID != "" && databaseID == authInfoID {
			database = conn.PrivateDB(databaseID)
		} else {
			return nil, skyerr.NewError(skyerr.PermissionDenied, "The selected DB cannot be accessed because permission is denied")
		}
	}

	return database, nil
}
