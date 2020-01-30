package router

import (
	"os"

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

	// files url
	gopath, _ := os.Getwd()
	e.Static("/", gopath+"/assets")

	return e
}
