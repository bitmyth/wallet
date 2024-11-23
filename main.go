package main

import (
	_ "embed"
	"github.com/bitmyth/walletserivce/factory"
	"github.com/bitmyth/walletserivce/route"
	"log"
	"net/http"
)

func main() {
	f, err := factory.New()
	if err != nil {
		log.Fatal(err)
	}
	logger := f.Logger()

	db, err := f.DB()
	if err != nil {
		logger.Error(err)
		return
	}
	err = db.Migrate()
	if err != nil {
		logger.Error(err)
		return
	}

	router := route.Router(f)
	f.RegisterRoutes(router)

	logger.Infoln("http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
