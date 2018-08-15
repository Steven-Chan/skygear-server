package handler

import (
	"encoding/json"
	"net/http"

	rModel "github.com/skygeario/skygear-server/cmd/skynext-router/model"
)

type ConfigHandler struct{}

func (h ConfigHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	accessKeyType := rModel.GetAccessKeyType(r)
	if accessKeyType != rModel.MasterAccessKey {
		http.Error(w, "Master key required", http.StatusUnauthorized)
		return
	}

	config := rModel.GetAuthConfig(r)
	encoder := json.NewEncoder(w)
	encoder.Encode(config)
}
