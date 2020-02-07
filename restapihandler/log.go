package restapihandler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"northstar/models"
	"strconv"
	"strings"

	"github.com/ayannahindonesia/basemodel"

	"github.com/labstack/echo"
)

// LogList shows log list
func LogList(c echo.Context) error {
	defer c.Request().Body.Close()

	var (
		logs      models.Log
		result    basemodel.PagedFindResult
		unmarsh   []string
		rows      int
		page      int
		startDate string
		endDate   string
	)

	// pagination parameters
	rows, _ = strconv.Atoi(c.QueryParam("rows"))
	if rows > 0 {
		page, _ = strconv.Atoi(c.QueryParam("page"))
		if page <= 0 {
			page = 1
		}
	}
	orderby := strings.Split(c.QueryParam("orderby"), ",")
	sort := strings.Split(c.QueryParam("sort"), ",")

	if startDate = c.QueryParam("start_date"); len(startDate) > 0 {
		if endDate := c.QueryParam("end_date"); len(endDate) < 1 {
			endDate = startDate
		}
	}

	b, err := ioutil.ReadAll(c.Request().Body)
	log.Printf("body : %v", string(b))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	err = json.Unmarshal(b, &unmarsh)
	if err != nil {
		unmarsh = []string{}
	}

	result, err = logs.PagedFindFilter(page, rows, orderby, sort, &models.LogQueryFilter{
		Client:    c.QueryParam("client"),
		Level:     c.QueryParam("level"),
		StartDate: startDate,
		EndDate:   endDate,
		Messages:  unmarsh,
	})
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}

	return c.JSON(http.StatusOK, result)
}
