package models

import (
	"strings"

	"github.com/jinzhu/gorm"
)

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
