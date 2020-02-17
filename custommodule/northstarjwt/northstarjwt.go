package northstarjwt

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// JWT main type
type (
	JWT struct {
		Secret   string
		Duration int64
	}
	// JWTclaims type
	JWTclaims struct {
		ClientID string `json:"client_id"`
		jwt.StandardClaims
	}
)

// New jwt instance
func New(secret string, duration int64) JWT {
	return JWT{Secret: secret, Duration: duration}
}

// SetClientJWTmiddlewares func
func (j *JWT) SetClientJWTmiddlewares(g *echo.Group) {
	g.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: "HS512",
		SigningKey:    []byte(j.Secret),
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

// CreateJwtToken func
func (j *JWT) CreateJwtToken(id string) (string, error) {
	claim := JWTclaims{
		id,
		jwt.StandardClaims{
			Id:        id,
			ExpiresAt: time.Now().Add(time.Duration(j.Duration) * time.Minute).Unix(),
		},
	}

	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claim)
	token, err := rawToken.SignedString([]byte(j.Secret))
	if err != nil {
		return "", err
	}

	return token, nil
}
