package db

import (
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

var unregisteredScript = `#!ipxe

menu Boot Menu - Unregistered ({mac})
item --gap -- -------------------------------
item local      Boot from local disk
item register   Register Device
item netbootxyz Boot to netboot.xyz menu
choose target && goto ${target}

:local
exit

:register
echo -n Hostname:
read hostname
chain --autofree http://${next-server}/api/new/host/${net0/mac}/${hostname}

:netbootxyz
chain --autofree http://boot.netboot.xyz/
`

var registeredScript = `#!ipxe

menu Boot Menu - Registered as {hostname}
item --gap -- -------------------------------
item local      Boot from local disk
item netbootxyz Boot to netboot.xyz menu
choose --default local --timeout 3000 target || goto local
goto ${target}

:local
exit

:register
read hostname
chain --autofree http://${next-server}/api/new-host/${net0/mac}/${hostname}

:netbootxyz
chain --autofree http://boot.netboot.xyz/
`

func GetScriptByMAC(mac string, db *gorm.DB, log bool) (string, error) {
	ctx := context.Background()

	host, err := gorm.G[Host](db).Where("mac = ?", mac).First(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if log {
			LogRequest(false, time.Now(), mac, db)
		}
		script := strings.ReplaceAll(unregisteredScript, "{mac}", mac)
		return script, nil
	} else if err != nil {
		return "", err
	}
	if log {
		LogRequest(true, time.Now(), mac, db)
	}

	task, err := gorm.G[Task](db).Where("id = ?", host.TaskID).First(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		script := strings.ReplaceAll(registeredScript, "{hostname}", host.Name)
		return script, nil
	} else if err != nil {
		return "", err
	}

	script := strings.ReplaceAll(task.Script, "{hostname}", host.Name)

	return script, nil
}
