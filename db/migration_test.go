package db

import (
	"github.com/bitmyth/walletserivce/config"
	"testing"
)

func TestDB_Migrate(t *testing.T) {
	c, _ := config.NewConfig()
	db, _ := Open(c)

	if err := db.Migrate(); err != nil {
		t.Error(err)
		return
	}
}
