package db

import (
	"context"
	"errors"
	"regexp"

	"gorm.io/gorm"
)

type Host struct {
	gorm.Model
	Name   string
	Mac    string
	TaskID int
	Task   Task
}

func CreateHost(mac, hostname string, taskID int, db *gorm.DB) error {
	macRegex := regexp.MustCompile(`^([0-9A-Fa-f]{2}:){5}[0-9A-Fa-f]{2}$`)
	if !macRegex.MatchString(mac) {
		return errors.New("invalid mac address")
	}

	ctx := context.Background()

	err := gorm.G[Host](db).Create(ctx, &Host{Name: hostname, Mac: mac, TaskID: taskID})
	if err != nil {
		return err
	}

	return nil
}

func EditHost(name, mac string, taskID, id int, db *gorm.DB) error {
	ctx := context.Background()

	_, err := gorm.G[Host](db).Where("id = ?", id).Updates(ctx, Host{Name: name, Mac: mac, TaskID: taskID})
	if err != nil {
		return err
	}

	return nil
}

func DeleteHost(id string, db *gorm.DB) error {
	ctx := context.Background()

	_, err := gorm.G[Host](db).Where("id = ?", id).Delete(ctx)
	if err != nil {
		return err
	}

	return nil
}

func GetHostByID(id string, db *gorm.DB) (*Host, error) {
	ctx := context.Background()

	host, err := gorm.G[Host](db).Where("id = ?", id).Preload("Task", nil).First(ctx)
	if err != nil {
		return nil, err
	}

	return &host, nil
}
