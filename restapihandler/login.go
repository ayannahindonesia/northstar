package restapihandler

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"northstar/application"
	"northstar/models"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
)

// ClientLogin func
func ClientLogin(c echo.Context) error {
	defer c.Request().Body.Close()

	data, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Basic "))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "invalid login")
	}

	auth := strings.Split(string(data), ":")
	if len(auth) < 2 {
		return c.JSON(http.StatusUnauthorized, "invalid login")
	}
	type Login struct {
		Key    string `json:"key"`
		Secret string `json:"secret"`
	}

	client := models.Client{}
	err = client.SingleFindFilter(&Login{
		Key:    auth[0],
		Secret: auth[1],
	})
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "invalid login")
	}

	token, err := application.App.JWT.CreateJwtToken(strconv.FormatUint(client.ID, 10))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "invalid login")
	}

	jwtConf := application.App.Config.GetStringMap(fmt.Sprintf("%s.jwt", application.App.ENV))
	expiration := time.Duration(jwtConf["duration"].(int)) * time.Minute

	return c.JSON(http.StatusOK, map[string]interface{}{
		"token":      token,
		"expires_in": expiration.Seconds(),
	})
}
