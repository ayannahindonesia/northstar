package models

import (
	"math"
	"strings"

	"github.com/ayannahindonesia/basemodel"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
)

// Log main type
type Log struct {
	basemodel.BaseModel
	Level    string         `json:"level" gorm:"column:level;type:varchar(255)"`
	Messages postgres.Jsonb `json:"messages" gorm:"column:messages"`
}

// Create func
func (model *Log) Create() error {
	return basemodel.Create(&model)
}

// Save func
func (model *Log) Save() error {
	return basemodel.Save(&model)
}

// FirstOrCreate func
func (model *Log) FirstOrCreate() error {
	return basemodel.FirstOrCreate(&model)
}

// Delete log
func (model *Log) Delete() error {
	return basemodel.Delete(&model)
}

// SingleFindFilter search using filter and return last
func (model *Log) SingleFindFilter(filter interface{}) error {
	return basemodel.SingleFindFilter(&model, filter)
}

// PagedFindFilter search using filter and return with pagination format
func (model *Log) PagedFindFilter(page int, rows int, order []string, sort []string, filter map[string]interface{}) (basemodel.PagedFindResult, error) {
	if page <= 0 {
		page = 1
	}

	query := basemodel.DB

	query = conditionQuery(query, filter)
	query = orderSortQuery(query, order, sort)

	temp := query
	var totalRows int

	temp.Find(&model).Count(&totalRows)

	var (
		offset   int
		lastPage int
		err      error
	)

	if rows > 0 {
		offset = (page * rows) - rows
		lastPage = int(math.Ceil(float64(totalRows) / float64(rows)))

		query = query.Limit(rows).Offset(offset)
	}

	err = query.Find(&model).Error

	result := basemodel.PagedFindResult{
		TotalData:   totalRows,
		Rows:        rows,
		CurrentPage: page,
		LastPage:    lastPage,
		From:        offset + 1,
		To:          offset + rows,
		Data:        model,
	}

	return result, err
}

func conditionQuery(query *gorm.DB, filter map[string]interface{}) *gorm.DB {
	query = query.Joins("JOIN LATERAL jsonb_array_elements(log.messages) j ON true")
	for k, v := range filter {
		query = query.Where("value->>'?' = ?", k, v)
	}

	return query
}

func orderSortQuery(query *gorm.DB, order []string, sort []string) *gorm.DB {
	for k, v := range order {
		q := v
		if len(sort) > k {
			value := sort[k]
			if strings.ToUpper(value) == "ASC" || strings.ToUpper(value) == "DESC" {
				q = v + " " + strings.ToUpper(value)
			}
		}
		query = query.Order(q)
	}

	return query
}
