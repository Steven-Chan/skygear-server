package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/skygeario/skygear-server/pkg/core/asset"
	"github.com/skygeario/skygear-server/pkg/core/auth"
	"github.com/skygeario/skygear-server/pkg/core/auth/authz"
	"github.com/skygeario/skygear-server/pkg/core/auth/authz/policy"
	"github.com/skygeario/skygear-server/pkg/core/db"
	"github.com/skygeario/skygear-server/pkg/core/handler"
	"github.com/skygeario/skygear-server/pkg/core/inject"
	"github.com/skygeario/skygear-server/pkg/core/model"
	"github.com/skygeario/skygear-server/pkg/core/server"
	recordGear "github.com/skygeario/skygear-server/pkg/record"
	"github.com/skygeario/skygear-server/pkg/record/dependency/record"
	"github.com/skygeario/skygear-server/pkg/record/dependency/recordconv"
	"github.com/skygeario/skygear-server/pkg/server/skyerr"
)

func AttachSaveHandler(
	server *server.Server,
	recordDependency recordGear.DependencyMap,
) *server.Server {
	server.Handle("/save", &SaveHandlerFactory{
		recordDependency,
	}).Methods("POST")
	return server
}

type SaveHandlerFactory struct {
	Dependency recordGear.DependencyMap
}

func (f SaveHandlerFactory) NewHandler(request *http.Request) http.Handler {
	h := &SaveHandler{}
	inject.DefaultInject(h, f.Dependency, request)
	return handler.APIHandlerToHandler(h, h.TxContext)
}

func (f SaveHandlerFactory) ProvideAuthzPolicy() authz.Policy {
	return policy.AllOf(
		authz.PolicyFunc(policy.DenyNoAccessKey),
		authz.PolicyFunc(policy.RequireAuthenticated),
		authz.PolicyFunc(policy.DenyDisabledUser),
	)
}

type SaveRequestPayload struct {
	Atomic bool `json:"atomic"`

	// RawMaps stores the original incoming `records`.
	RawMaps []map[string]interface{} `json:"records"`

	// IncomigItems contains de-serialized recordID or de-serialization error,
	// the item is one-one corresponding to RawMaps.
	IncomingItems []interface{}

	// Records contains the successfully de-serialized record
	Records []*record.Record

	// Errs is the array of de-serialization errors
	Errs []skyerr.Error
}

func (s SaveRequestPayload) Validate() error {
	if len(s.RawMaps) == 0 {
		return skyerr.NewInvalidArgument("expected list of record", []string{"records"})
	}

	return nil
}

func (s SaveRequestPayload) isClean() bool {
	return len(s.Errs) == 0
}

/*
SaveHandler is dummy implementation on save/modify Records
curl -X POST -H "Content-Type: application/json" \
  -d @- http://localhost:3000/save <<EOF
{
    "records": [{
        "_id": "note/EA6A3E68-90F3-49B5-B470-5FFDB7A0D4E8",
        "content": "ewdsa",
        "_access": [{
            "role": "admin",
            "level": "write"
        }]
    }]
}
EOF

Save with reference
curl -X POST -H "Content-Type: application/json" \
  -d @- http://localhost:3000/save <<EOF
{
  "records": [
    {
      "collection": {
        "$type": "ref",
        "$id": "collection/10"
      },
      "noteOrder": 1,
      "content": "hi",
      "_id": "note/71BAE736-E9C5-43CB-ADD1-D8633B80CAFA",
      "_type": "record",
      "_access": [{
          "role": "admin",
          "level": "write"
      }]
    }
  ]
}
EOF
*/
type SaveHandler struct {
	AuthContext auth.ContextGetter `dependency:"AuthContextGetter"`
	TxContext   db.TxContext       `dependency:"TxContext"`
	RecordStore record.Store       `dependency:"RecordStore"`
	Logger      *logrus.Entry      `dependency:"HandlerLogger"`
	AssetStore  asset.Store        `dependency:"AssetStore"`
}

func (h SaveHandler) WithTx() bool {
	return false
}

func (h SaveHandler) DecodeRequest(request *http.Request) (handler.RequestPayload, error) {
	payload := SaveRequestPayload{}
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		return nil, err
	}

	for _, recordMap := range payload.RawMaps {
		var r record.Record
		if err := (*recordconv.JSONRecord)(&r).FromMap(recordMap); err != nil {
			skyErr := skyerr.NewError(skyerr.InvalidArgument, err.Error())
			payload.Errs = append(payload.Errs, skyErr)
			payload.IncomingItems = append(payload.IncomingItems, skyErr)
		} else {
			r.SanitizeForInput()
			payload.IncomingItems = append(payload.IncomingItems, r.ID)
			payload.Records = append(payload.Records, &r)
		}
	}

	return payload, nil
}

func (h SaveHandler) Handle(req interface{}) (resp interface{}, err error) {
	payload := req.(SaveRequestPayload)

	resultFilter, err := NewRecordResultFilter(
		h.RecordStore,
		h.TxContext,
		h.AssetStore,
		h.AuthContext.AuthInfo(),
		h.AuthContext.AccessKeyType() == model.MasterAccessKey,
	)
	if err != nil {
		err = skyerr.MakeError(err)
		return
	}

	modifyReq := RecordModifyRequest{
		RecordStore:   h.RecordStore,
		TxContext:     h.TxContext,
		AssetStore:    h.AssetStore,
		Logger:        h.Logger,
		AuthInfo:      h.AuthContext.AuthInfo(),
		RecordsToSave: payload.Records,
		Atomic:        payload.Atomic,
		WithMasterKey: h.AuthContext.AccessKeyType() == model.MasterAccessKey,
		ModifyAt:      timeNow(),
	}
	modifyResp := RecordModifyResponse{
		ErrMap: map[record.ID]skyerr.Error{},
	}

	// Open transaction for whole operation if atomic save
	if payload.Atomic {
		if err = h.TxContext.BeginTx(); err != nil {
			return
		}
	}

	var opErr error
	defer func() {
		if payload.Atomic {
			if txErr := db.EndTx(h.TxContext, opErr); txErr != nil {
				err = txErr
			}
		} else {
			err = opErr
		}
	}()

	if opErr = h.ExtendRecordSchemaWithTx(payload); opErr != nil {
		return
	}

	if opErr = RecordSaveHandler(&modifyReq, &modifyResp); opErr != nil {
		// Override error in atomic save
		if payload.Atomic && len(modifyResp.ErrMap) > 0 {
			info := map[string]interface{}{}
			for recordID, err := range modifyResp.ErrMap {
				info[recordID.String()] = err
			}

			opErr = skyerr.NewErrorWithInfo(skyerr.AtomicOperationFailure,
				"Atomic Operation rolled back due to one or more errors",
				info)
			return
		}

		return
	}

	results := make([]interface{}, 0, len(payload.RawMaps))
	h.makeResultsFromIncomingItem(payload.IncomingItems, modifyResp, resultFilter, &results)

	resp = results

	return
}

// ExtendRecordSchemaWithTx ensure the operation is within a transaction.
// When the request is inatomic, the operation would be wrapped in a new transaction.
func (h SaveHandler) ExtendRecordSchemaWithTx(payload SaveRequestPayload) error {
	return executeFuncInTx(h.TxContext, payload.Atomic, func() (err error) {
		// TODO: emit schema updated event
		_, err = ExtendRecordSchema(h.RecordStore, h.Logger, payload.Records)
		if err != nil {
			h.Logger.WithError(err).Errorln("failed to migrate record schema")
			if _, ok := err.(skyerr.Error); !ok {
				err = skyerr.NewError(skyerr.IncompatibleSchema, "failed to migrate record schema")
			}

			return
		}

		return
	})
}

func (h SaveHandler) makeResultsFromIncomingItem(incomingItems []interface{}, resp RecordModifyResponse, resultFilter RecordResultFilter, results *[]interface{}) {
	currRecordIdx := 0
	for _, itemi := range incomingItems {
		var result interface{}

		switch item := itemi.(type) {
		case skyerr.Error:
			result = newSerializedError("", item)
		case record.ID:
			if err, ok := resp.ErrMap[item]; ok {
				h.Logger.WithFields(logrus.Fields{
					"recordID": item,
					"err":      err,
				}).Debugln("failed to save record")

				result = newSerializedError(item.String(), err)
			} else {
				record := resp.SavedRecords[currRecordIdx]
				currRecordIdx++
				result = resultFilter.JSONResult(record)
			}
		default:
			panic(fmt.Sprintf("unknown type of incoming item: %T", itemi))
		}

		*results = append(*results, result)
	}
}
