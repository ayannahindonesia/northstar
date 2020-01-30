package middlewares

import (
	"fmt"
	"net/http"
	"northstar/application"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// SetClientJWTmiddlewares func
func SetClientJWTmiddlewares(g *echo.Group) {
	jwtConfig := application.App.Config.GetStringMap(fmt.Sprintf("%s.jwt", application.App.ENV))

	g.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: "HS512",
		SigningKey:    []byte(jwtConfig["jwt_secret"].(string)),
	}))

	g.Use(validateJWT)
}

func validateJWT(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user")
		token := user.(*jwt.Token)

		if _, ok := token.Claims.(jwt.MapClaims); ok {
			return next(c)
		}

		return echo.NewHTTPError(http.StatusForbidden, fmt.Sprintf("%s", "invalid token"))
	}
}
