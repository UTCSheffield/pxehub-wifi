package db

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func OpenDB(dbPath string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(dbPath))
	if err != nil {
		panic(fmt.Sprintf("failed to open db: %s", err))
	}

	db.AutoMigrate(&Script{})
	db.AutoMigrate(&Host{})

	return db
}
