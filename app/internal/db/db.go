package db

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/jackc/pgx/v4"
)

const timeout = time.Second * 2

// Database - used to connect to the database
type Database struct {
	Conn *pgx.Conn
}

var singleton *Database
var once sync.Once

// GetDB returns the singleton database instance
func GetDB() (*Database, error) {
	var err error
	once.Do(func() {
		rand.Seed(time.Now().UTC().UnixNano())
		singleton = &Database{}
		err = singleton.init()
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

func (d *Database) init() error {
	config, err := pgx.ParseConfig(getDbConnStr())
	if err != nil {
		log.Fatalln("Unable to parse database config:", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	d.Conn, err = pgx.ConnectConfig(ctx, config)

	return err
}
