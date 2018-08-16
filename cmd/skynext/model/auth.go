package model

import (
	"encoding/json"

	"github.com/skygeario/skygear-server/pkg/server/skydb"
)

func GetAuthInfo(i interface{}) *skydb.AuthInfo {
	str := header(i).Get("X-Skygear-AuthInfo")
	info := skydb.AuthInfo{}
	err := json.Unmarshal([]byte(str), &info)
	if err != nil {
		return nil
	}

	return &info
}

func SetAuthInfo(i interface{}, authInfo *skydb.AuthInfo) {
	data, err := json.Marshal(*authInfo)
	if err != nil {
		header(i).Del("X-Skygear-AuthInfo")
		return
	}

	header(i).Set("X-Skygear-AuthInfo", string(data))
}

func GetAccessToken(i interface{}) string {
	return header(i).Get("X-Skygear-AccessToken")
}

func SetAccessToken(i interface{}, at string) {
	header(i).Set("X-Skygear-AccessToken", at)
}

func GetAuthInfoID(i interface{}) string {
	return header(i).Get("X-Skygear-AuthInfoID")
}

func SetAuthInfoID(i interface{}, at string) {
	header(i).Set("X-Skygear-AuthInfoID", at)
}
