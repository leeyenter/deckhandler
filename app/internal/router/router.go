package router

import (
	"errors"

	"github.com/labstack/echo/v4"
)

type Router struct {
}

func New() *Router {
	r := Router{}
	return &r
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
	// return c.JSON(http.StatusCreated, map[string]interface{}{
	// 	"deck_id":   "hi",
	// 	"shuffled":  true,
	// 	"remaining": 30,
	// })
	return errors.New("not implemented")
}

// GET /:id
func (r *Router) OpenDeck(c echo.Context) error {
	return errors.New("not implemented")
}

// POST/:id/draw
// Payload: count
func (r *Router) DrawCards(c echo.Context) error {
	return errors.New("not implemented")
}
