package db

import (
	"time"

	"gorm.io/gorm"
)

type Request struct {
	gorm.Model
	Registered bool
	Time       time.Time
	Mac        string
}
