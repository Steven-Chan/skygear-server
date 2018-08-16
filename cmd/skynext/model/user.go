package model

import (
	"encoding/json"
	"net/http"
)

type User struct {
	Username string `json:"name"`
	Age      int    `json:"age"`
}

func header(i interface{}) http.Header {
	switch i.(type) {
	case *http.Request:
		return (i.(*http.Request)).Header
	case http.ResponseWriter:
		return (i.(http.ResponseWriter)).Header()
	default:
		panic("Invalid type")
	}
}

func GetUserID(i interface{}) string {
	return header(i).Get("X-Skygear-UserID")
}

func GetUsername(i interface{}) string {
	return header(i).Get("X-Skygear-Username")
}

func SetUsername(i interface{}, un string) {
	header(i).Set("X-Skygear-Username", un)
}

func GetUser(i interface{}) *User {
	user := User{}
	if err := json.Unmarshal([]byte(header(i).Get("X-Skygear-User")), &user); err != nil {
		return nil
	}
	return &user
}

func SetUser(i interface{}, u *User) {
	if u == nil {
		header(i).Del("X-Skygear-User")
		return
	}

	user, err := json.Marshal(u)
	if err != nil {
		panic(err)
	}

	header(i).Set("X-Skygear-User", string(user))
}
