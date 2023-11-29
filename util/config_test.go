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
