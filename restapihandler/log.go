package restapihandler

import (
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
		rows      int
		page      int
		startDate string
		endDate   string
		err       error
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
		if endDate = c.QueryParam("end_date"); len(endDate) < 1 {
			endDate = startDate
		}
	}

	result, err = logs.PagedFindFilter(page, rows, orderby, sort, &models.LogQueryFilter{
		Client:    c.QueryParam("client"),
		Tag:       c.QueryParam("tag"),
		Note:      c.QueryParam("note"),
		UID:       c.QueryParam("uid"),
		Username:  c.QueryParam("username"),
		Level:     c.QueryParam("level"),
		StartDate: startDate,
		EndDate:   endDate,
		Messages:  strings.Split(c.QueryParam("messages"), ","),
	})
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}

	return c.JSON(http.StatusOK, result)
}

// LogDetail func
func LogDetail(c echo.Context) error {
	defer c.Request().Body.Close()

	var logDetail models.Log

	type LogSearchID struct {
		ID string `json:"id"`
	}

	err := logDetail.SingleFindFilter(&LogSearchID{ID: c.Param("id")})
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}

	return c.JSON(http.StatusOK, logDetail)
}
