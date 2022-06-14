package router

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
	"github.com/leeyenter/deckhandler/internal/data"
	"github.com/leeyenter/deckhandler/internal/db"
	"github.com/leeyenter/deckhandler/logger"
)

type DeckResponseWithCards struct {
	DeckID    string                   `json:"deck_id"`
	Shuffled  bool                     `json:"shuffled"`
	Remaining int                      `json:"remaining"`
	Cards     []map[string]interface{} `json:"cards"`
}

type DeckResponseWithoutCards struct {
	DeckID    string `json:"deck_id"`
	Shuffled  bool   `json:"shuffled"`
	Remaining int    `json:"remaining"`
}

// CreateDeck creates a deck and returns it. It takes two optional parameters:
// `shuffle` (boolean): if the deck should be shuffled or not. Defaults to false.
// `cards` ([]string): list of card codes to be included. Defaults to include the whole deck.
func (r *Router) CreateDeck(c echo.Context) error {
	cards := c.QueryParam("cards")
	shuffled := c.QueryParam("shuffle")

	// Check if shuffled is either none, "true" or "false"
	if shuffled != "" && shuffled != "true" && shuffled != "false" {
		logger.Get("ROUTER").Warn("Unrecognised shuffle value: " + shuffled)
		return c.JSON(http.StatusBadRequest, generateError("Unrecognised shuffle value"))
	}

	shuffledBool := shuffled == "true" // if "true", then shuffling; else, leave it

	// Get cards
	allCards, err := r.DB.FetchCards()
	if err != nil {
		logger.Get("ROUTER").Error("Error retrieving cards")
		return c.NoContent(http.StatusInternalServerError)
	}

	// First prepare the list of cards that are to be added to the deck
	deckCards := make([]data.Card, 0)

	if cards == "" {
		// No list provided; store all
		deckCards = allCards
	} else {
		// Prepare look-up for list of cards to make subsetting easier
		deckCardsMap := make(map[string]data.Card)
		for _, card := range allCards {
			deckCardsMap[card.ID] = card
		}

		cardCodes := strings.Split(cards, ",")
		for _, code := range cardCodes {
			if card, ok := deckCardsMap[code]; !ok {
				logger.Get("ROUTER").Warn("Unrecognised card code: " + code)
				return c.JSON(http.StatusBadRequest, generateError("Unrecognised card code '"+code+"'"))
			} else {
				deckCards = append(deckCards, card)
			}
		}
	}

	// Create deck, with given config and cards
	var resp DeckResponseWithoutCards
	if resp.DeckID, err = r.DB.CreateDeck(shuffledBool, deckCards); err != nil {
		logger.Get("ROUTER").Error("Error creating deck")
		return c.NoContent(http.StatusInternalServerError)
	}

	resp.Shuffled = shuffledBool
	resp.Remaining = len(deckCards)

	return c.JSON(http.StatusOK, resp)
}

// OpenDeck takes a URL param `id`, and retrieves data about the deck
// with that given `id`.
func (r *Router) OpenDeck(c echo.Context) error {
	deckId := c.Param("id")

	if deckId == "" {
		return c.JSON(http.StatusBadRequest, generateError("Deck ID must be included."))
	}

	// Retrieve deck info
	deck, err := r.DB.GetDeck(deckId)
	if err == pgx.ErrNoRows || db.GetPGErr(err) == "22P02" {
		return c.JSON(http.StatusBadRequest, generateError("Unrecognised deck ID '"+deckId+"'"))
	} else if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	// Format the deck's cards into expected JSON format
	deckCards := make([]map[string]interface{}, len(deck.Cards))
	for i, card := range deck.Cards {
		deckCards[i] = card.ToMap()
	}

	return c.JSON(http.StatusOK, DeckResponseWithCards{
		DeckID:    deck.ID,
		Shuffled:  deck.Shuffled,
		Remaining: len(deck.Cards),
		Cards:     deckCards,
	})
}

// DrawCards lists the first `count` cards from deck `id`,
// by the order that it was created, and removes them from the deck.
func (r *Router) DrawCards(c echo.Context) error {
	deckId := c.Param("id")
	countStr := c.QueryParam("count")

	// Input validation
	if deckId == "" {
		return c.JSON(http.StatusBadRequest, generateError("Deck ID must be included."))
	}

	if countStr == "" {
		return c.JSON(http.StatusBadRequest, generateError("Count must be included."))
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, generateError("Invalid count value"))
	}

	// First check to see if deck ID is valid
	_, err = r.DB.GetDeck(deckId)
	if err == pgx.ErrNoRows || db.GetPGErr(err) == "22P02" {
		return c.JSON(http.StatusBadRequest, generateError("Unrecognised deck ID '"+deckId+"'"))
	}

	// Retrieve the cards from the deck
	cards, err := r.DB.FetchCardsFromDeck(deckId, count)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	} else if len(cards) != count {
		return c.JSON(http.StatusBadRequest, generateError("Drawing too many cards"))
	}

	cardCodes := make([]string, len(cards))
	mappedCards := make([]map[string]interface{}, len(cards))
	for i, card := range cards {
		cardCodes[i] = card.ID
		mappedCards[i] = card.ToMap()
	}

	if err := r.DB.RemoveCardsFromDeck(deckId, cardCodes); err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	// All ok, return the drawn cards
	return c.JSON(http.StatusOK, map[string][]map[string]interface{}{
		"cards": mappedCards,
	})
}
