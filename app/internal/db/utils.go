package db

import (
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgconn"
	"github.com/leeyenter/deckhandler/internal/data"
)

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

// GetPGErr is a helper function to quickly extract the PG error code
// from a given error.
func GetPGErr(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code
	}
	return ""
}
