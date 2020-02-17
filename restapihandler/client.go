package restapihandler

import (
	"net/http"
	"northstar/application"
	"northstar/models"

	"github.com/labstack/echo"
)

// ClientList shows client list
func ClientList(c echo.Context) error {
	defer c.Request().Body.Close()

	var (
		client models.Client
		err    error
	)

	result, err := client.ClientNameList(application.App.DB)
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}

	return c.JSON(http.StatusOK, result)
}
