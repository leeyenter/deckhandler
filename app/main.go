package main

import (
	"log"
	"os"

	"github.com/leeyenter/deckhandler/internal/data"
	"github.com/leeyenter/deckhandler/internal/db"
)

var cards []data.Card
var dbObj *db.Database

func main() {
	var err error
	cards, err = data.LoadCSVFile("assets/cards.csv")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	dbObj, err = db.GetDB()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
