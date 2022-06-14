package router

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/leeyenter/deckhandler/internal/data"
	"github.com/leeyenter/deckhandler/internal/db"
	"github.com/leeyenter/deckhandler/logger"
)

type Router struct {
	DB *db.Database
}

func New() (*Router, error) {
	dbObj, err := db.GetDB()
	if err != nil {
		return nil, err
	}

	r := Router{DB: dbObj}
	return &r, nil
}

type DeckResponse struct {
	DeckID    string                   `json:"deck_id"`
	Shuffled  bool                     `json:"shuffled"`
	Remaining int                      `json:"remaining"`
	Cards     []map[string]interface{} `json:"cards"`
}

// POST /
// URL params may include `cards`, `shuffled`
func (r *Router) CreateDeck(c echo.Context) error {
	cards := c.QueryParam("cards")
	shuffled := c.QueryParam("shuffled")

	if shuffled != "" && shuffled != "true" && shuffled != "false" {
		logger.Get("ROUTER").Info("Unrecognised shuffled value: " + shuffled)
		return c.NoContent(http.StatusBadRequest)
	}

	shuffledBool := shuffled == "true" // if "true", then shuffling; else, leave it

	// Get cards
	allCards, err := r.DB.FetchCards()
	if err != nil {
		logger.Get("ROUTER").Info("Error retrieving cards")
		return c.NoContent(http.StatusInternalServerError)
	}

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
				logger.Get("ROUTER").Info("Unrecognised card code: " + code)
				return c.NoContent(http.StatusBadRequest)
			} else {
				deckCards = append(deckCards, card)
			}
		}
	}

	var resp DeckResponse
	if resp.DeckID, err = r.DB.CreateDeck(shuffledBool, deckCards); err != nil {
		logger.Get("ROUTER").Info("Error creating deck")
		return c.NoContent(http.StatusInternalServerError)
	}

	resp.Shuffled = shuffledBool
	resp.Remaining = len(deckCards)

	return c.JSON(http.StatusOK, resp)
}

// GET /:id
func (r *Router) OpenDeck(c echo.Context) error {
	deckId := c.Param("id")
	if deck, err := r.DB.GetDeck(deckId); err != nil {
		return c.NoContent(http.StatusBadRequest)
	} else {
		deckCards := make([]map[string]interface{}, len(deck.Cards))
		for i, card := range deck.Cards {
			deckCards[i] = card.ToMap()
		}

		return c.JSON(http.StatusOK, DeckResponse{
			DeckID:    deck.ID,
			Shuffled:  deck.Shuffled,
			Remaining: len(deck.Cards),
			Cards:     deckCards,
		})
	}
}

// POST/:id/draw?count=
// Payload: count
func (r *Router) DrawCards(c echo.Context) error {
	deckId := c.Param("id")
	countStr := c.QueryParam("count")

	if deckId == "" || countStr == "" {
		return c.NoContent(http.StatusBadRequest)
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	cards, err := r.DB.FetchCardsFromDeck(deckId, count)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	} else if len(cards) != count {
		return c.NoContent(http.StatusBadRequest)
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
	return c.JSON(http.StatusOK, DeckResponse{
		Cards: mappedCards,
	})
}
