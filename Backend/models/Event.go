package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Event struct {
	gorm.Model
	DateTime    time.Time
	Description string
	Location    string
	fee         float64
}

const (
	ColumnDateTime = "date_time"
	ColumnDescription= "description"
	ColumnLocation = "location"
	ColumnDateFee = "fee"
)
