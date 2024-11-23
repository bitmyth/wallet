package fixtures

import (
	"github.com/bitmyth/walletserivce/db"
	"log"
)

type dependency interface {
	DB() (*db.DB, error)
	Redis() (*db.Redis, error)
}

func PreloadTestingData(f dependency) {
	d, _ := f.DB()

	_, _ = d.Exec("truncate table users CASCADE")
	_, _ = d.Exec("truncate table transactions")
	// Prepare the SQL statement
	insertSQL := `INSERT INTO users (username, balance) VALUES ($1, $2), ($3, $4)`

	username1 := "user1"
	balance1 := 100.0
	username2 := "user2"
	balance2 := 100.0

	// Execute the SQL statement
	_, err := d.Exec(insertSQL, username1, balance1, username2, balance2)
	if err != nil {
		log.Fatal("error inserting rows:", err)
	}
}
