package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/leeyenter/deckhandler/internal/db"
)

type Router struct {
	DB     *db.Database
	Server string
}

func New() (*Router, error) {
	dbObj, err := db.GetDB()
	if err != nil {
		return nil, err
	}

	r := Router{DB: dbObj}
	return &r, nil
}

func (r *Router) BuildServer() *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.Use(routerLogger)
	e.Logger.SetLevel(log.OFF)
	e.POST("/", r.CreateDeck)
	e.GET("/:id", r.OpenDeck)
	e.POST("/:id/draw", r.DrawCards)
	return e
}
