package config

import (
	"testing"
)

func TestNewConfig(t *testing.T) {
	SetConfigPath("../")
	config, err := NewConfig()
	if err != nil {
		t.Error(err)
		return
	}
	if config.Host == "" {
		t.Error("postgre host is empty")
	}
}
