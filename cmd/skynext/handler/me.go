package handler

import (
	"encoding/json"
	"net/http"

	"github.com/skygeario/skygear-server/cmd/skynext/model"
)

type MeHandler struct{}

func (h MeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if username := model.GetUsername(r); username == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user := model.GetUser(r)
	encoder := json.NewEncoder(w)
	encoder.Encode(user)
}
