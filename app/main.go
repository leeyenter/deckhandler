package main

import (
	"log"
	"os"

	"github.com/leeyenter/deckhandler/internal/db"
	"github.com/leeyenter/deckhandler/internal/router"
	"github.com/leeyenter/deckhandler/internal/utils"
	"github.com/leeyenter/deckhandler/logger"
)

func main() {
	seedPath := utils.Getenv("CARDS_PATH", "assets/cards.csv")
	if os.Getenv("SKIP_SEED") == "true" {
		logger.Get("MAIN").Info("Skipping seeding")
		seedPath = ""
	} else {
		logger.Get("MAIN").Info("Seeding database using data from " + seedPath + ". All preexisting data will be removed.")
	}

	db.Init(seedPath)

	r, err := router.New()
	if err != nil {
		log.Fatal(err)
	}

	server := r.BuildServer()
	server.Start(":" + utils.Getenv("PORT", "3000"))
}
