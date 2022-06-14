package db

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/leeyenter/deckhandler/logger"
)

const timeout = time.Second * 2

// Database - used to connect to the database
type Database struct {
	Conn *pgx.Conn
}

var singleton *Database
var once sync.Once

// Initialises the database. Includes creating the singleton object,
// and seeding the database if needed.
func Init(file string) (*Database, error) {
	var err error
	once.Do(func() {
		config, err := pgx.ParseConfig(getDbConnStr())
		if err != nil {
			log.Fatalln("Unable to parse database config:", err)
		}

		config.Logger = zapadapter.NewLogger(logger.Get("DB"))
		config.LogLevel = pgx.LogLevelWarn

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		singleton = &Database{}
		singleton.Conn, err = pgx.ConnectConfig(ctx, config)
		if err != nil {
			return
		}

		// Optionally pass in a file to seed the database.
		if file != "" {
			err = singleton.seedData(file)
			if err != nil {
				logger.Get("DB").Error("Error seeding database")
				return
			} else {
				logger.Get("DB").Info("Database successfully seeded")
			}
		}

	})
	return singleton, err
}

// GetDB returns the singleton database instance
func GetDB() (*Database, error) {
	var err error
	Init("../../assets/cards.csv") // call init with testing-suitable path for seeding the database
	return singleton, err
}
