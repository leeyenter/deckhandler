package main

import (
	"log"

	"github.com/leeyenter/deckhandler/internal/db"
	"github.com/leeyenter/deckhandler/internal/router"
)

// var dbObj *db.Database

func main() {

	db.Init("")

	r, err := router.New()
	if err != nil {
		log.Fatal(err)
	}

	server := r.BuildServer()
	server.Start(":3000")
}
