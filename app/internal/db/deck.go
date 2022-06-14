package db

import (
	"context"
	"math/rand"

	"github.com/jackc/pgx/v4"
	"github.com/leeyenter/deckhandler/internal/data"
	"github.com/leeyenter/deckhandler/logger"
)

// CreateDeck creates a deck in the database,
// given a list of cards and whether they should
// be shuffled. Returns the id.
func (d *Database) CreateDeck(shuffled bool, cards []data.Card) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if shuffled {
		rand.Shuffle(len(cards), func(i, j int) {
			cards[i], cards[j] = cards[j], cards[i]
		})
	}

	tx, err := d.Conn.Begin(ctx)
	if err != nil {
		logger.Get("DB-DECK").Error("CreateDeck: Transaction failed - " + err.Error())
		return "", err
	}

	defer tx.Rollback(ctx)

	row := tx.QueryRow(ctx, `INSERT INTO decks (shuffled) VALUES ($1) RETURNING id`, shuffled)
	var id string
	if err := row.Scan(&id); err != nil {
		logger.Get("DB-DECK").Error("CreateDeck: Insert deck failed - " + err.Error())
		return "", err
	}

	for _, card := range cards {
		if _, err := tx.Exec(ctx, `INSERT INTO deck_cards (deck_id, card_code) VALUES ($1, $2)`, id, card.ID); err != nil {
			logger.Get("DB-DECK").Error("CreateDeck: Insert card failed - " + err.Error())
			return "", err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		logger.Get("DB-DECK").Error("CreateDeck: Committing - " + err.Error())
		return "", err
	}

	return id, nil
}

// GetDeck retrieves information of deck `id`.
func (d *Database) GetDeck(id string) (data.Deck, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	deck := data.Deck{ID: id}

	row := d.Conn.QueryRow(ctx, `SELECT shuffled FROM decks WHERE id = $1`, id)
	if err := row.Scan(&deck.Shuffled); err != nil {
		logger.Get("DB-DECK").Error("GetDeck: Select failed - " + err.Error())
		return deck, err
	}

	var err error
	if deck.Cards, err = d.FetchCardsFromDeck(id, -1); err != nil {
		return deck, err
	}

	return deck, nil
}

// FetchCardsFromDeck returns the first `count` cards in a deck
// of a given `id`. If `count` is -1, return the whole deck.
func (d *Database) FetchCardsFromDeck(id string, count int) ([]data.Card, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	const query = `SELECT code, value FROM deck_cards INNER JOIN cards ON cards.code = deck_cards.card_code WHERE deck_id = $1`
	var rows pgx.Rows
	var err error

	if count < 0 {
		rows, err = d.Conn.Query(ctx, query, id)
	} else {
		rows, err = d.Conn.Query(ctx, query+` LIMIT $2`, id, count)
	}

	if err != nil {
		logger.Get("DB-DECK").Error("FetchCardsFromDeck: Select failed - " + err.Error())
		return nil, err
	}

	cards := make([]data.Card, 0)

	for rows.Next() {
		var card data.Card
		if err = rows.Scan(&card.ID, &card.Values); err != nil {
			logger.Get("DB-DECK").Error("FetchCardsFromDeck: Scan failed - " + err.Error())
			return nil, err
		} else {
			cards = append(cards, card)
		}
	}

	return cards, nil
}

// RemoveCardsFromDeck removes the first `count` cards in a deck
// of given `id`.
func (d *Database) RemoveCardsFromDeck(id string, cardCodes []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for _, cardCode := range cardCodes {
		_, err := d.Conn.Exec(ctx, `DELETE FROM deck_cards WHERE deck_id = $1 AND card_code = $2`, id, cardCode)
		if err != nil {
			logger.Get("DB-DECK").Error("RemoveCardsFromDeck - " + err.Error())
			return err
		}
	}

	return nil
}

// ClearDecks is a helper function that removes
// all decks in the database. Used mainly for testing.
func (d *Database) ClearDecks() error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if _, err := d.Conn.Exec(ctx, `DELETE FROM decks`); err != nil {
		logger.Get("DB-DECK").Error("ClearDecks - " + err.Error())
	}

	return nil
}
