package db

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type WifiKey struct {
	gorm.Model
	Key string `gorm:"uniqueIndex"`
}

var ErrNoWifiKeys = errors.New("no available wifi keys")

func GetUnassignedWifiKey(db *gorm.DB) (*WifiKey, error) {
	var key *WifiKey
	db.
		Joins("LEFT JOIN hosts ON hosts.wifi_key_id = wifi_keys.id").
		Where("hosts.id IS NULL").
		First(&key)

	return key, nil
}

func setWifiKeyStatus(id uint, used bool, db *gorm.DB) error {
	ctx := context.Background()

	_, err := gorm.G[WifiKey](db).Where("id = ?", id).Update(ctx, "used", used)
	if err != nil {
		return err
	}

	return nil
}

func SetWifiKeyAsUnused(id uint, db *gorm.DB) error {
	return setWifiKeyStatus(id, false, db)
}

func AssignWifiKeyToHost(hostID int, db *gorm.DB) (*WifiKey, error) {
	ctx := context.Background()

	key, err := GetUnassignedWifiKey(db)
	if err != nil {
		return nil, err
	}

	_, err = gorm.G[Host](db).Where("id = ?", hostID).Updates(ctx, Host{WifiKey: *key})
	if err != nil {
		return nil, err
	}

	return key, nil
}
