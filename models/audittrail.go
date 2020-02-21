package models

import (
	"math"

	"github.com/ayannahindonesia/basemodel"
	"github.com/jinzhu/gorm"
)

type (
	// Audittrail main type
	Audittrail struct {
		basemodel.BaseModel
		Client   string `json:"client" gorm:"column:client;type:varchar(255)"`
		UserID   string `json:"user_id" gorm:"column:user_id;type:varchar(255)"`
		Username string `json:"username" gorm:"column:username;type:varchar(255)"`
		Roles    string `json:"roles" gorm:"column:roles;type:varchar(255)"`
		Entity   string `json:"entity" gorm:"column:entity;type:varchar(255)"`
		EntityID string `json:"entity_id" gorm:"column:entity_id;type:varchar(255)"`
		Action   string `json:"action" gorm:"column:action;type:varchar(255)"`
		Original string `json:"original" gorm:"column:original;type:text"`
		New      string `json:"new" gorm:"column:new;type:text"`
	}

	// AudittrailQueryFilter filter struct
	AudittrailQueryFilter struct {
		Client    string
		User      string
		Username  string
		Entity    string
		EntityID  string
		Action    string
		Original  []string
		New       []string
		StartDate string
		EndDate   string
	}
)

// Create func
func (model *Audittrail) Create() error {
	return basemodel.Create(&model)
}

// Save func
func (model *Audittrail) Save() error {
	return basemodel.Save(&model)
}

// FirstOrCreate func
func (model *Audittrail) FirstOrCreate() error {
	return basemodel.FirstOrCreate(&model)
}

// Delete trail
func (model *Audittrail) Delete() error {
	return basemodel.Delete(&model)
}

// SingleFindFilter search using filter and return last
func (model *Audittrail) SingleFindFilter(filter interface{}) error {
	return basemodel.SingleFindFilter(&model, filter)
}

// PagedFindFilter search using filter and return with pagination format
func (model *Audittrail) PagedFindFilter(page int, rows int, order []string, sort []string, filter *AudittrailQueryFilter) (basemodel.PagedFindResult, error) {
	if page <= 0 {
		page = 1
	}

	query := basemodel.DB
	models := []Audittrail{}

	query = audittrailConditionQuery(query, filter)
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

func audittrailConditionQuery(query *gorm.DB, filter *AudittrailQueryFilter) *gorm.DB {
	for _, v := range filter.Original {
		if len(v) > 0 {
			query = query.Where("original LIKE ?", "%"+v+"%")
		}
	}
	for _, v := range filter.New {
		if len(v) > 0 {
			query = query.Where("new LIKE ?", "%"+v+"%")
		}
	}

	if len(filter.Client) > 0 {
		query = query.Where("client = ?", filter.Client)
	}

	if len(filter.Entity) > 0 {
		query = query.Where("entity = ?", filter.Entity)
	}

	if len(filter.EntityID) > 0 {
		query = query.Where("entity_id = ?", filter.EntityID)
	}

	if len(filter.Action) > 0 {
		query = query.Where("action = ?", filter.Action)
	}

	if len(filter.StartDate) > 0 {
		if len(filter.EndDate) < 1 {
			filter.EndDate = filter.StartDate
		}
		query = query.Where("created_at BETWEEN ? AND ?", filter.StartDate, filter.EndDate)
	}

	return query
}
