package util

import (
	"encoding/json"
	"testing"
)

func TestConfig(t *testing.T) {
	sc := ServerConfig{
		DBConfig:  &DBConfig{},
		LogConfig: &LogConfig{},
	}
	bts, _ := json.Marshal(sc)
	println(string(bts))
}

func TestGetBundlerPrivateKeysFromEnv(t *testing.T) {
	type DBKeys struct {
		BundlerPrivateKeys []string `json:"bundler_private_keys"`
	}

	var dbKeys DBKeys = DBKeys{
		BundlerPrivateKeys: []string{"7013b62758059b6fbd08bd38a987c54e6a50cc4d306788db744c1818f18a08a5"},
	}
	bts, _ := json.Marshal(dbKeys)
	println(string(bts))
}
