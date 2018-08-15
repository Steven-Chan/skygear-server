package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/skygeario/skygear-server/cmd/skynext/model"
	"github.com/skygeario/skygear-server/cmd/skynext/service"
)

type LoginHandler struct {
	UserAuthenticator service.UserAuthenticator `inject:"UserAuthenticator"`
	UserStore         service.UserStore         `inject:"UserStore"`
}

type LoginHandlerPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	payload := LoginHandlerPayload{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&payload); err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	authenticated := h.UserAuthenticator.Authenticate(payload.Username, payload.Password)
	if !authenticated {
		http.Error(w, "Authentication data not match with credential", http.StatusBadRequest)
		return
	}
	model.SetUsername(w, payload.Username)

	user := h.UserStore.Get(payload.Username)
	if user == nil {
		fmt.Fprintf(w, "{}")
		return
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(user)
}
