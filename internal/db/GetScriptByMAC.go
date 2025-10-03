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
item exit       Exit iPXE
item register   Register Device
item netbootxyz netboot.xyz
choose target && goto ${target}

:register
echo -n Hostname:
read hostname
chain --autofree http://${next-server}/api/new/host/${net0/mac}/${hostname}

:netbootxyz
chain --autofree http://boot.netboot.xyz/

:exit
exit
`

var registeredScript = `#!ipxe

menu Boot Menu - Registered as {hostname}
item --gap -- -------------------------------
item exit       Exit iPXE
item netbootxyz netboot.xyz
choose --default local --timeout 3000 target || goto local
goto ${target}

:register
read hostname
chain --autofree http://${next-server}/api/new-host/${net0/mac}/${hostname}

:netbootxyz
chain --autofree http://boot.netboot.xyz/

:exit
exit
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
	} else if log {
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

	var zero int = 0

	if !host.PermanentTask {
		EditHost(host.Name, host.Mac, &zero, false, host.ID, db)
	}
	return script, nil
}
