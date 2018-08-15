package model

import (
	"strconv"
)

type Config struct {
	Auth       AuthConfig
	Record     RecordConfig
	AssetStore AssetStoreConfig
}

type AuthConfig struct {
	PasswordLength int
}

type RecordConfig struct {
	AutoMigration bool
}

type AssetStoreConfig struct {
	Impl   string
	Secret string
}

func SetConfig(i interface{}, config Config) {
	SetAuthConfig(i, config.Auth)
	SetRecordConfig(i, config.Record)
	SetAssetStoreConfig(i, config.AssetStore)
}

func GetAuthConfig(i interface{}) AuthConfig {
	var pl int
	plv, err := strconv.Atoi(header(i).Get("X-Skygear-Config-Auth-PasswordLength"))
	if err == nil {
		pl = plv
	} else {
		pl = 0
	}

	return AuthConfig{
		PasswordLength: pl,
	}
}

func SetAuthConfig(i interface{}, config AuthConfig) {
	if config.PasswordLength != 0 {
		header(i).Set("X-Skygear-Config-Auth-PasswordLength", strconv.Itoa(config.PasswordLength))
	} else {
		header(i).Del("X-Skygear-Config-Auth-PasswordLength")
	}
}

func GetRecordConfig(i interface{}) RecordConfig {
	var am bool
	amv, err := strconv.ParseBool(header(i).Get("X-Skygear-Config-Record-AutoMigration"))
	if err == nil {
		am = amv
	} else {
		am = false
	}

	return RecordConfig{
		AutoMigration: am,
	}
}

func SetRecordConfig(i interface{}, config RecordConfig) {
	if config.AutoMigration != false {
		header(i).Set("X-Skygear-Config-Record-AutoMigration", strconv.FormatBool(config.AutoMigration))
	} else {
		header(i).Del("X-Skygear-Config-Record-AutoMigration")
	}
}

func GetAssetStoreConfig(i interface{}) AssetStoreConfig {
	impl := header(i).Get("X-Skygear-Config-AssetStore-Impl")
	secret := header(i).Get("X-Skygear-Config-AssetStore-Secret")

	return AssetStoreConfig{
		Impl:   impl,
		Secret: secret,
	}
}

func SetAssetStoreConfig(i interface{}, config AssetStoreConfig) {
	header(i).Set("X-Skygear-Config-AssetStore-Impl", config.Impl)
	header(i).Set("X-Skygear-Config-AssetStore-Secret", config.Secret)
}
