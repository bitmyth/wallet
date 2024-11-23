package db

import (
	"github.com/bitmyth/walletserivce/config"
	"testing"
)

func TestMain(m *testing.M) {
	config.SetConfigPath("../")
	m.Run()
}

func TestOpenPG(t *testing.T) {
	c, _ := config.NewConfig()

	d, err := Open(c)
	if err != nil {
		t.Error(err)
	}

	alive := d.IsAlive()
	if !alive {
		t.Error("pg is not alive")
	}

	c.Host = ""
	_, err = Open(c)
	if err == nil {
		t.Error("expect return error")
	}
}

func TestOpenRedis(t *testing.T) {
	c, _ := config.NewConfig()

	r, err := OpenRedis(c)
	if err != nil {
		t.Error(err)
	}
	alive := r.IsAlive()
	if !alive {
		t.Error("redis is not alive")
	}

	c.Redis.Addr = "0"
	_, err = OpenRedis(c)
	if err == nil {
		t.Error("expect return error")
	}
}
