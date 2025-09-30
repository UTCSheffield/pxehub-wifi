package db

import (
	"context"
	"errors"
	"regexp"

	"gorm.io/gorm"
)

func CreateHost(mac, hostname string, db *gorm.DB) error {
	macRegex := regexp.MustCompile(`^([0-9A-Fa-f]{2}:){5}[0-9A-Fa-f]{2}$`)
	if !macRegex.MatchString(mac) {
		return errors.New("invalid mac address")
	}

	ctx := context.Background()

	err := gorm.G[Host](db).Create(ctx, &Host{Name: hostname, Mac: mac})
	if err != nil {
		return err
	}

	return nil
}
