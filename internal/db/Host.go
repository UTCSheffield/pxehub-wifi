package db

import "gorm.io/gorm"

type Host struct {
	gorm.Model
	Name   string
	Mac    string
	TaskID int
	Task   Task
}
