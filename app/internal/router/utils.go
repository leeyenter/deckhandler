package router

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/leeyenter/deckhandler/logger"
)

func routerLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		otherInfo := ""
		err := next(c)
		if err != nil {
			c.Error(err)
			otherInfo = err.Error()
		}

		req := c.Request()
		res := c.Response()
		log := fmt.Sprintf("%3d %-4s %-15s %s %s", res.Status, req.Method, c.RealIP(), req.RequestURI, otherInfo)

		switch {
		case res.Status >= 500:
			logger.Get("ROUTER").Error(log)
		case res.Status >= 400:
			logger.Get("ROUTER").Warn(log)
		default:
			logger.Get("ROUTER").Info(log)
		}

		return nil
	}
}

func generateError(message string) map[string]string {
	return map[string]string{
		"error": message,
	}
}
