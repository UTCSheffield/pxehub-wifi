package db

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func OpenDB(dbPath string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic(fmt.Sprintf("failed to open db: %s", err))
	}

	db.AutoMigrate(&Task{})
	db.AutoMigrate(&Host{})
	db.AutoMigrate(&Request{})
	db.AutoMigrate(&WifiKey{})

	return db
}
