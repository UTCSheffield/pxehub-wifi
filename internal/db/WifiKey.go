package db

import (
	"context"
	"errors"
	"strconv"

	"gorm.io/gorm"
)

type WifiKey struct {
	gorm.Model
	Key string `gorm:"uniqueIndex"`
}

func CreateWifiKey(key string, db *gorm.DB) error {
	ctx := context.Background()

	err := gorm.G[WifiKey](db).Create(ctx, &WifiKey{Key: key})
	if err != nil {
		return err
	}

	return nil
}

func EditWifiKey(key string, id uint, db *gorm.DB) error {
	ctx := context.Background()

	_, err := gorm.G[WifiKey](db).Update(ctx, "key", key)

	return err
}

func DeleteWifiKey(id uint, db *gorm.DB) error {
	ctx := context.Background()

	_, err := gorm.G[WifiKey](db).Where("id = ?", id).Delete(ctx)
	if err != nil {
		return err
	}

	return nil
}

func GetWifiKeyByID(id string, db *gorm.DB) (*WifiKey, error) {
	ctx := context.Background()

	key, err := gorm.G[WifiKey](db).Where("id = ?", id).First(ctx)
	if err != nil {
		return nil, err
	}

	return &key, nil
}

var ErrNoWifiKeys = errors.New("no available wifi keys")

func GetUnassignedWifiKey(db *gorm.DB) (*WifiKey, error) {
	var key *WifiKey
	err := db.
		Joins("LEFT JOIN hosts ON hosts.wifi_key_id = wifi_keys.id").
		Where("hosts.id IS NULL").
		First(&key).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNoWifiKeys
	} else if err != nil {
		return nil, err
	}

	return key, nil
}

func AssignWifiKeyToHost(hostID int, db *gorm.DB) (*WifiKey, error) {
	ctx := context.Background()

	key, err := GetUnassignedWifiKey(db)
	if err != nil {
		return nil, err
	} else if key.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	_, err = gorm.G[Host](db).Where("id = ?", hostID).Updates(ctx, Host{WifiKeyID: &key.ID})
	if err != nil {
		return nil, err
	}

	return key, nil
}

func GetOrAssignWifiKeyToHost(hostID uint, db *gorm.DB) (*WifiKey, error) {
	ctx := context.Background()

	host, err := GetHostByID(strconv.Itoa(int(hostID)), db)
	if err != nil {
		return nil, err
	} else if host.WifiKeyID != nil {
		return &host.WifiKey, nil
	}

	key, err := GetUnassignedWifiKey(db)
	if err != nil {
		return nil, err
	} else if key.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	_, err = gorm.G[Host](db).Where("id = ?", hostID).Updates(ctx, Host{WifiKeyID: &key.ID})
	if err != nil {
		return nil, err
	}

	return key, nil
}
