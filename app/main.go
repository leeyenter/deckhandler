package main

import (
	"log"
	"os"

	"github.com/leeyenter/deckhandler/internal/db"
)

var dbObj *db.Database

func main() {
	var err error

	dbObj, err = db.GetDB()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
