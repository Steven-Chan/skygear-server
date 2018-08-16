package handler

import (
	"fmt"
	"net/http"

	"github.com/skygeario/skygear-server/pkg/server/skyerr"

	rModel "github.com/skygeario/skygear-server/cmd/skynext-router/model"
	"github.com/skygeario/skygear-server/cmd/skynext/service"
	skyHandler "github.com/skygeario/skygear-server/pkg/server/handler"
	"github.com/skygeario/skygear-server/pkg/server/router"
)

type SchemaFetchHandler struct {
	handler          skyHandler.SchemaFetchHandler
	DatabaseProvider service.DatabaseProvider `inject:"DatabaseProvider"`
}

func (h SchemaFetchHandler) Preprocess(r *http.Request, rpayload *router.Payload, response *router.Response) int {
	if rModel.GetAccessKeyType(r) != rModel.MasterAccessKey {
		response.Err = skyerr.NewError(skyerr.AccessKeyNotAccepted, "Master key expected")
		return http.StatusBadRequest
	}

	h.handler = skyHandler.SchemaFetchHandler{}

	var err error
	if rpayload.Database, err = h.DatabaseProvider.GetDatabase(r); err != nil || rpayload.Database == nil {
		response.Err = skyerr.MakeError(err)
		return http.StatusServiceUnavailable
	}
	rpayload.DBConn = rpayload.Database.Conn()
	return http.StatusOK
}

func (h SchemaFetchHandler) Handle(rpayload *router.Payload, response *router.Response) {
	fmt.Printf("database: %+v\n", rpayload.Database)
	h.handler.Handle(rpayload, response)
}
