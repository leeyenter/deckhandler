package db

import (
	"context"

	"github.com/leeyenter/deckhandler/internal/data"
)

// CreateCards adds a list of new cards into the database, that
// can subsequently be added into decks.
func (d *Database) CreateCards(cards []data.Card) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	tx, err := d.Conn.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	for _, card := range cards {
		if _, err := tx.Exec(
			ctx,
			`INSERT INTO cards (code, value) VALUES ($1, $2)`,
			card.ID, card.Values,
		); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

// FetchCards retrieves all cards from the database.
func (d *Database) FetchCards() ([]data.Card, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	rows, err := d.Conn.Query(ctx, `SELECT code, value FROM cards`)
	if err != nil {
		return nil, err
	}

	cards := make([]data.Card, 0)

	for rows.Next() {
		var card data.Card
		if err = rows.Scan(&card.ID, &card.Values); err != nil {
			return nil, err
		} else {
			cards = append(cards, card)
		}
	}

	return cards, nil
}

// ClearCards is a helper function that removes
// all cards in the database. Will result in decks
// having no cards, so `ClearDecks` needs to be called
// as well.
func (d *Database) ClearCards() error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_, err := d.Conn.Exec(ctx, `DELETE FROM cards`)

	return err
}
