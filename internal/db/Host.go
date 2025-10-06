package db

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"gorm.io/gorm"
)

type Host struct {
	gorm.Model
	Name string `gorm:"unique"`
	Mac  string `gorm:"unique"`

	TaskID        *int
	Task          Task
	PermanentTask bool

	WifiKeyID *uint `gorm:"unique"`
	WifiKey   WifiKey
}

func CreateHost(mac, hostname string, taskID int, taskPerm bool, db *gorm.DB) error {
	macRegex := regexp.MustCompile(`^([0-9A-Fa-f]{2}:){5}[0-9A-Fa-f]{2}$`)
	if !macRegex.MatchString(mac) {
		return errors.New("invalid mac address")
	}
	mac = strings.ToLower(mac)

	ctx := context.Background()

	err := gorm.G[Host](db).Create(ctx, &Host{Name: hostname, Mac: mac, TaskID: &taskID, PermanentTask: taskPerm})
	if err != nil {
		return err
	}

	return nil
}

func EditHost(name, mac string, taskID *int, taskPerm bool, id uint, db *gorm.DB) error {
	var host Host

	if err := db.First(&host, id).Error; err != nil {
		return err
	}

	host.Name = name
	host.Mac = mac
	if *taskID == 0 {
		host.TaskID = nil
	} else {
		host.TaskID = taskID
	}
	host.PermanentTask = taskPerm

	db.Save(&host)

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

	host, err := gorm.G[Host](db).Where("id = ?", id).Preload("Task", nil).Preload("WifiKey", nil).First(ctx)
	if err != nil {
		return nil, err
	}

	return &host, nil
}

func GetHostByMAC(mac string, db *gorm.DB) (*Host, error) {
	ctx := context.Background()

	host, err := gorm.G[Host](db).Where("LOWER(mac) = LOWER(?)", mac).Preload("Task", nil).Preload("WifiKey", nil).First(ctx)
	if err != nil {
		return nil, err
	}

	return &host, nil
}
