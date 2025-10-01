package db

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Request struct {
	gorm.Model
	Registered bool
	Time       time.Time
	Mac        string
}

func LogRequest(registered bool, datetime time.Time, mac string, db *gorm.DB) error {
	ctx := context.Background()

	err := gorm.G[Request](db).Create(ctx, &Request{Registered: registered, Time: datetime, Mac: mac})
	if err != nil {
		return err
	}

	return nil
}
