package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/leeyenter/deckhandler/internal/data"
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
				return
			}
		}

	})
	return singleton, err
}

// GetDB returns the singleton database instance
func GetDB() (*Database, error) {
	var err error
	once.Do(func() {
		// If haven't yet initialised the database,
		// init with values for running tests
		singleton = &Database{}
		err = singleton.Init("../../assets/cards.csv")
	})
	return singleton, err
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func getDbConnStr() string {
	dbUser := getenv("DB_USER", "deckhandler")
	dbPass := getenv("DB_PASS", "wmWLWyoqsKJtXwisAqwaPkA9yT8MvrzRj")
	dbHost := getenv("DB_HOST", "127.0.0.1")
	dbPort := getenv("DB_PORT", "5432")
	return fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%s", dbUser, dbUser, dbPass, dbHost, dbPort)
}

func (d *Database) Init(file string) error {
	config, err := pgx.ParseConfig(getDbConnStr())
	if err != nil {
		log.Fatalln("Unable to parse database config:", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if d.Conn, err = pgx.ConnectConfig(ctx, config); err != nil {
		return err
	}

	return d.seedData(file) // Included here to make development easier/cleaner
}

// Clears all existing data in the database,
// and adds in the cards data.
func (d *Database) seedData(file string) error {
	if err := d.ClearCards(); err != nil {
		return err
	}

	if err := d.ClearDecks(); err != nil {
		return err
	}

	cards, err := data.LoadCSVFile(file)
	if err != nil {
		return err
	}

	if err = d.CreateCards(cards); err != nil {
		return err
	}

	return nil
}
