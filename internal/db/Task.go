package db

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	ID     int
	Name   string
	Script string `gorm:"type:longtext"`
}
