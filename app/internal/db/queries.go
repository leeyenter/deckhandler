package db

import (
	"errors"

	"github.com/leeyenter/deckhandler/internal/data"
)

// CreateDeck creates a deck in the database,
// given a list of cards and whether they should
// be shuffled. Returns the id.
func (d *Database) CreateDeck(shuffled bool, cards []data.Card) (string, error) {
	d.Conn.Exec()
	return "", errors.New("not implemented")
}

// GetDeck retrieves information of deck `id`.
func (d *Database) GetDeck(id string) (data.Deck, error) {
	return data.Deck{}, errors.New("not implemented")
}

// FetchCards returns the first `count` cards in a deck
// of a given `id`. If `count` is -1, return the whole deck.
func (d *Database) FetchCards(id string, count int) ([]data.Card, error) {
	return nil, errors.New("not implemented")
}

// RemoveCards removes the first `count` cards in a deck
// of given `id`.
func (d *Database) RemoveCards(id string, count int) error {
	return errors.New("not implemented")
}
