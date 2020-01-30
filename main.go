package main

import (
	"flag"
	"log"
	"northstar/application"
	"northstar/router"
	"os"

	"github.com/labstack/echo/middleware"
)

var (
	flags = flag.NewFlagSet("goose", flag.ExitOnError)
)

func main() {
	defer application.App.Close()

	flags.Usage = usage
	flags.Parse(os.Args[1:])
	args := flags.Args()

	switch args[0] {
	default:
		flags.Usage()
		break
	case "run":
		e := router.NewRouter()
		// CORS react handle
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowHeaders: []string{"*"},
		}))

		e.Logger.Fatal(e.Start(":" + application.App.Port))
		os.Exit(0)
		break
	}
}

func usage() {
	usagestring := ``

	log.Print(usagestring)
}
