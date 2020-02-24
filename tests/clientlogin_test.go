package tests

import (
	"net/http"
	"net/http/httptest"
	"northstar/router"
	"testing"

	"github.com/gavv/httpexpect"
)

func TestClientLogin(t *testing.T) {
	DataBuild()

	api := router.NewRouter()

	server := httptest.NewServer(api)

	defer server.Close()

	e := httpexpect.New(t, server.URL)

	auth := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Basic "+clientBasicToken)
	})

	obj := auth.GET("/login").
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.Keys().Contains("expires_in", "token")
}
