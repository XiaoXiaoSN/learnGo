package main

import (
	"time"

	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Name     string
	Age      int64
	Birthday time.Time
	Email    string `gorm:"type:varchar(100);unique_index"`
	Role     string `gorm:"size:255"`   // set field size to 255
	Address  string `gorm:"index:addr"` // create index with name `addr` for address
	IgnoreMe int    `gorm:"-"`          // ignore this field
}
