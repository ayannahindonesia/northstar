package group

import (
	"northstar/application"
	"northstar/restapihandler"

	"github.com/labstack/echo"
)

// NorthstarGroup group
func NorthstarGroup(e *echo.Echo) {
	g := e.Group("/ns")
	application.App.JWT.SetClientJWTmiddlewares(g)

	g.POST("/log", restapihandler.LogList)
}
