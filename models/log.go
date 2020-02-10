package models

import (
	"math"
	"strings"

	"github.com/ayannahindonesia/basemodel"
	"github.com/jinzhu/gorm"
)

type (
	// Log main type
	Log struct {
		basemodel.BaseModel
		Client   string `json:"client" gorm:"column:client;type:varchar(255)"`
		Tag      string `json:"tag" gorm:"column:tag;type:varchar(255)"`
		Level    string `json:"level" gorm:"column:level;type:varchar(255);default:'info'"`
		Messages string `json:"messages" gorm:"column:messages;type:text"`
	}
	// LogQueryFilter filter struct
	LogQueryFilter struct {
		Client    string
		Level     string
		StartDate string
		EndDate   string
		Messages  []string
	}
)

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
func (model *Log) PagedFindFilter(page int, rows int, order []string, sort []string, filter *LogQueryFilter) (basemodel.PagedFindResult, error) {
	if page <= 0 {
		page = 1
	}

	query := basemodel.DB
	models := []Log{}

	query = conditionQuery(query, filter)
	query = orderSortQuery(query, order, sort)

	temp := query
	var totalRows int

	temp.Find(&models).Count(&totalRows)

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

	err = query.Find(&models).Error

	result := basemodel.PagedFindResult{
		TotalData:   totalRows,
		Rows:        rows,
		CurrentPage: page,
		LastPage:    lastPage,
		From:        offset + 1,
		To:          offset + rows,
		Data:        models,
	}

	return result, err
}

func conditionQuery(query *gorm.DB, filter *LogQueryFilter) *gorm.DB {
	for _, v := range filter.Messages {
		query = query.Where("messages LIKE ?", "%"+v+"%")
	}

	if len(filter.Client) > 0 {
		query = query.Where("client = ?", filter.Client)
	}

	if len(filter.StartDate) > 0 {
		if len(filter.EndDate) < 1 {
			filter.EndDate = filter.StartDate
		}
		query = query.Where("created_at BETWEEN ? AND ?", filter.StartDate, filter.EndDate)
	}

	if len(filter.Level) > 0 {
		query = query.Where("level = ?", filter.Level)
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
