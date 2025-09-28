package db

import (
	"context"

	"gorm.io/gorm"
)

func GetScriptByMAC(mac string, db *gorm.DB) (string, error) {
	ctx := context.Background()

	host, err := gorm.G[Host](db).Where("mac = ?", mac).First(ctx)
	if err != nil {
		return "", err
	}

	script, err := gorm.G[Script](db).Where("ID = ?", host.ScriptID).First(ctx)
	if err != nil {
		return "", err
	}

	return script.Content, nil
}
