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

	g.GET("/log", restapihandler.LogList)
	g.GET("/audittrail", restapihandler.AuditTrailList)
	g.GET("/audittrail/:id", restapihandler.AuditTrailDetail)
	g.GET("/client", restapihandler.ClientList)
}
