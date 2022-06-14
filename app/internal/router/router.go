package router

import "github.com/labstack/echo/v4"

type Router struct {
}

func New() *Router {
	r := Router{}
	return &r
}

func (r *Router) CreateDeck(c echo.Context) error {
	return nil
}
