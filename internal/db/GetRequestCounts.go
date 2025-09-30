package db

import (
	"context"
	"time"

	"gorm.io/gorm"
)

func GetTotalRequestCounts(db *gorm.DB) (totalRequests, totalUnregisteredRequests, totalRegisteredRequests int64, err error) {
	ctx := context.Background()

	totalRequests, err = gorm.G[Request](db).Count(ctx, "ID")
	if err != nil {
		return 0, 0, 0, err
	}

	totalUnregisteredRequests, err = gorm.G[Request](db).Where("registered = ?", false).Count(ctx, "ID")
	if err != nil {
		return 0, 0, 0, err
	}

	totalRegisteredRequests, err = gorm.G[Request](db).Where("registered = ?", true).Count(ctx, "ID")
	if err != nil {
		return 0, 0, 0, err
	}

	return
}

func GetRequestGraphData(date time.Time, db *gorm.DB) (allPerDay, registeredPerDay, unregisteredPerDay []int, graphDates []string, err error) {
	ctx := context.Background()

	year, month, _ := date.Date()
	location := date.Location()
	firstOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, location)
	firstOfNextMonth := firstOfMonth.AddDate(0, 1, 0)
	daysInMonth := firstOfNextMonth.Add(-time.Nanosecond).Day()

	allPerDay = make([]int, daysInMonth)
	registeredPerDay = make([]int, daysInMonth)
	unregisteredPerDay = make([]int, daysInMonth)
	graphDates = make([]string, daysInMonth)

	for i := 0; i < daysInMonth; i++ {
		currentDay := firstOfMonth.AddDate(0, 0, i)
		graphDates[i] = currentDay.Format("2006-01-02")
	}

	var requests []Request
	requests, err = gorm.G[Request](db).Where("time >= ? AND time < ?", firstOfMonth, firstOfNextMonth).Find(ctx)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	for _, r := range requests {
		day := r.Time.Day() - 1
		allPerDay[day]++
		if r.Registered {
			registeredPerDay[day]++
		} else {
			unregisteredPerDay[day]++
		}
	}

	return
}
