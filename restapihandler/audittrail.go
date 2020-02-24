package restapihandler

import (
	"log"
	"net/http"
	"northstar/models"
	"strconv"
	"strings"

	"github.com/ayannahindonesia/basemodel"
	"github.com/labstack/echo"
)

// AuditTrailList list
func AuditTrailList(c echo.Context) error {
	defer c.Request().Body.Close()

	var (
		audittrail models.Audittrail
		result     basemodel.PagedFindResult
		rows       int
		page       int
		startDate  string
		endDate    string
		err        error
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

	log.Println(c.QueryParam("username"))
	result, err = audittrail.PagedFindFilter(page, rows, orderby, sort, &models.AudittrailQueryFilter{
		Client:    c.QueryParam("client"),
		UserID:    c.QueryParam("user"),
		Username:  c.QueryParam("username"),
		Entity:    c.QueryParam("entity"),
		EntityID:  c.QueryParam("entity_id"),
		Action:    c.QueryParam("action"),
		Original:  strings.Split(c.QueryParam("original"), ","),
		New:       strings.Split(c.QueryParam("new"), ","),
		StartDate: startDate,
		EndDate:   endDate,
	})
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}

	return c.JSON(http.StatusOK, result)
}

// AuditTrailDetail func
func AuditTrailDetail(c echo.Context) error {
	defer c.Request().Body.Close()

	var audittrail models.Audittrail

	type AudittrailSearchID struct {
		ID string `json:"id"`
	}

	err := audittrail.SingleFindFilter(&AudittrailSearchID{ID: c.Param("id")})
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}

	return c.JSON(http.StatusOK, audittrail)
}
