package db

import (
	"gorm.io/gorm"
)

type Script struct {
	gorm.Model
	ID      int
	Name    string
	Content string `gorm:"type:longtext"`
}

type Host struct {
	gorm.Model
	Name     string
	Mac      string
	ScriptID int
	Script   Script
}
