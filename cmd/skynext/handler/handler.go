package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/skygeario/skygear-server/pkg/server/router"
	"github.com/skygeario/skygear-server/pkg/server/skyerr"
	"github.com/skygeario/skygear-server/pkg/server/skyversion"
)

type SkygearHandler interface {
	Preprocess(r *http.Request, rpayload *router.Payload, response *router.Response) int
	Handle(rpayload *router.Payload, response *router.Response)
}

func FromSkygearHandler(h SkygearHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// start: copy from skygear router
		resp := router.NewResponse(w)
		version := strings.TrimPrefix(skyversion.Version(), "v")
		w.Header().Set("Server", fmt.Sprintf("Skygear Server/%s", version))

		payload, err := router.NewPayload(r)
		if err != nil {
			http.Error(w, "Unable to create request payload", http.StatusInternalServerError)
			return
		}
		// end: copy from skygear router

		// TODO: Set request tag

		httpStatus := h.Preprocess(r, payload, resp)
		if resp.Err != nil {
			if httpStatus == http.StatusOK {
				httpStatus = defaultStatusCode(resp.Err)
			}
			http.Error(w, resp.Err.Error(), httpStatus)
			return
		}

		h.Handle(payload, resp)

		// start: copy from skygear router
		writer := resp.Writer()
		if writer == nil {
			// The response is already written.
			return
		}

		writer.Header().Set("Content-Type", "application/json")

		if resp.Err != nil && httpStatus >= 200 && httpStatus <= 299 {
			httpStatus = defaultStatusCode(resp.Err)
		}

		writer.WriteHeader(httpStatus)
		if err := writeEntity(writer, resp); err != nil {
			panic(err)
		}
		// end: copy from skygear router
	})
}

// start: copy from skygear router
func writeEntity(w http.ResponseWriter, i interface{}) error {
	if w == nil {
		return errors.New("writer is nil")
	}
	return json.NewEncoder(w).Encode(i)
}

func defaultStatusCode(err skyerr.Error) int {
	httpStatus, ok := map[skyerr.ErrorCode]int{
		skyerr.NotAuthenticated:        http.StatusUnauthorized,
		skyerr.PermissionDenied:        http.StatusForbidden,
		skyerr.AccessKeyNotAccepted:    http.StatusUnauthorized,
		skyerr.AccessTokenNotAccepted:  http.StatusUnauthorized,
		skyerr.InvalidCredentials:      http.StatusUnauthorized,
		skyerr.InvalidSignature:        http.StatusUnauthorized,
		skyerr.BadRequest:              http.StatusBadRequest,
		skyerr.InvalidArgument:         http.StatusBadRequest,
		skyerr.IncompatibleSchema:      http.StatusConflict,
		skyerr.AtomicOperationFailure:  http.StatusConflict,
		skyerr.PartialOperationFailure: http.StatusOK,
		skyerr.Duplicated:              http.StatusConflict,
		skyerr.ConstraintViolated:      http.StatusConflict,
		skyerr.ResourceNotFound:        http.StatusNotFound,
		skyerr.UndefinedOperation:      http.StatusNotFound,
		skyerr.NotSupported:            http.StatusNotImplemented,
		skyerr.NotImplemented:          http.StatusNotImplemented,
		skyerr.PluginUnavailable:       http.StatusServiceUnavailable,
		skyerr.PluginTimeout:           http.StatusGatewayTimeout,
		skyerr.RecordQueryInvalid:      http.StatusBadRequest,
		skyerr.ResponseTimeout:         http.StatusServiceUnavailable,
		skyerr.DeniedArgument:          http.StatusForbidden,
		skyerr.RecordQueryDenied:       http.StatusForbidden,
		skyerr.NotConfigured:           http.StatusServiceUnavailable,
		skyerr.UserDisabled:            http.StatusForbidden,
		skyerr.VerificationRequired:    http.StatusForbidden,
	}[err.Code()]
	if !ok {
		if err.Code() < 10000 {
			logrus.Warnf("Error code %d (%v) does not have a default status code set. Assumed 500.", err.Code(), err.Code())
		}
		httpStatus = http.StatusInternalServerError
	}
	return httpStatus
}

// end: copy from skygear router
