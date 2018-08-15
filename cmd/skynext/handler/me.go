package handler

import (
	"encoding/json"
	"net/http"

	"github.com/skygeario/skygear-server/cmd/skynext/model"
)

type MeHandler struct{}

func (h MeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user := model.GetUser(r)
	encoder := json.NewEncoder(w)
	encoder.Encode(user)
}
