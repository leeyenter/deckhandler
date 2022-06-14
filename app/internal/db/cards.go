package db

import (
	"errors"

	"github.com/leeyenter/deckhandler/internal/data"
)

// CreateCard adds a new card into the database, that
// can subsequently be added into decks.
func (d *Database) CreateCard(card data.Card) error {
	return errors.New("not implemented")
}

// FetchCards retrieves all cards from the database.
func (d *Database) FetchCards() ([]data.Card, error) {
	return nil, errors.New("not implemented")
}

// ClearCards is a helper function that removes
// all cards in the database. Will result in decks
// having no cards, so `ClearDecks` needs to be called
// as well.
func (d *Database) ClearCards() error {
	return errors.New("not implemented")
}
