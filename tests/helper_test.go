package tests

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"northstar/application"
	"northstar/custommodule/seed"
	"northstar/router"
	"os"
	"strings"

	"github.com/gavv/httpexpect"
)

var clientBasicToken = base64.StdEncoding.EncodeToString([]byte("reactkey:reactsecret"))

func init() {
	if application.App.ENV != "development" {
		fmt.Printf("test aren't allowed in %s environment.", application.App.ENV)
		os.Exit(1)
	}
}

// DataBuild builds data for testing purpose
func DataBuild() {
	truncateAllTables()
	seed.Seed()
}

func getBearerToken(auth *httpexpect.Expect) string {
	api := router.NewRouter()

	server := httptest.NewServer(api)

	defer server.Close()

	obj := auth.GET("/login").
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.Keys().Contains("expires_in", "token")

	return obj.Value("token").String().Raw()
}

func truncateAllTables() {
	tables := strings.Join([]string{
		"audittrails",
		"clients",
		"logs",
	}, ", ")
	application.App.DB.Exec(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", tables))
}
