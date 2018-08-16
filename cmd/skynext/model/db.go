package model

func GetDatabaseID(i interface{}) string {
	return header(i).Get("X-Skygear-DatabaseID")
}

func SetDatabaseID(i interface{}, id string) {
	header(i).Set("X-Skygear-DatabaseID", id)
}
