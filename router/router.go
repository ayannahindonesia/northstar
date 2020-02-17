package router

import (
	"northstar/group"
	"northstar/restapihandler"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// NewRouter router func
func NewRouter() *echo.Echo {
	e := echo.New()

	// ignore /api-borrower
	e.Pre(middleware.Rewrite(map[string]string{
		"/northstar/*": "/$1",
	}))

	e.GET("/login", restapihandler.ClientLogin)

	group.NorthstarGroup(e)

	return e
}
