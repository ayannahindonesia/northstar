package tests

import (
	"net/http"
	"net/http/httptest"
	"northstar/router"
	"testing"

	"github.com/gavv/httpexpect"
)

func TestLogList(t *testing.T) {
	DataBuild()

	api := router.NewRouter()

	server := httptest.NewServer(api)

	defer server.Close()

	e := httpexpect.New(t, server.URL)

	auth := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Basic "+clientBasicToken)
	})

	token := getBearerToken(auth)

	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer "+token)
	})

	obj := auth.GET("/ns/log").
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("total_data").ValueNotEqual("total_data", 0)
}

func TestLogDetail(t *testing.T) {
	DataBuild()

	api := router.NewRouter()

	server := httptest.NewServer(api)

	defer server.Close()

	e := httpexpect.New(t, server.URL)

	auth := e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Basic "+clientBasicToken)
	})

	token := getBearerToken(auth)

	auth = e.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", "Bearer "+token)
	})

	obj := auth.GET("/ns/log/1").
		Expect().
		Status(http.StatusOK).JSON().Object()
	obj.ContainsKey("id").ValueEqual("id", 1)
}
