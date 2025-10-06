package db

import (
	"context"
	"fmt"
	"html/template"
	"time"

	"gorm.io/gorm"
)

func GetTotalRequestCount(db *gorm.DB) (totalRequests int64, err error) {
	ctx := context.Background()

	totalRequests, err = gorm.G[Request](db).Count(ctx, "ID")
	if err != nil {
		return 0, err
	}

	return
}

func GetRequestGraphData(date time.Time, db *gorm.DB) (registeredPerDay, unregisteredPerDay []int, graphDates []string, err error) {
	ctx := context.Background()

	year, month, _ := date.Date()
	location := date.Location()
	firstOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, location)
	firstOfNextMonth := firstOfMonth.AddDate(0, 1, 0)
	daysInMonth := firstOfNextMonth.Add(-time.Nanosecond).Day()

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
		return nil, nil, nil, err
	}

	for _, r := range requests {
		day := r.Time.Day() - 1
		if r.Registered {
			registeredPerDay[day]++
		} else {
			unregisteredPerDay[day]++
		}
	}

	return
}

func GetTotalHostCount(db *gorm.DB) (totalRequests int64, err error) {
	ctx := context.Background()

	totalRequests, err = gorm.G[Host](db).Count(ctx, "ID")
	if err != nil {
		return 0, err
	}

	return
}

func GetActiveTaskCount(db *gorm.DB) (totalRequests int64, err error) {
	ctx := context.Background()

	totalRequests, err = gorm.G[Host](db).Where("task_id IS NOT NULL").Count(ctx, "ID")
	if err != nil {
		return 0, err
	}

	return
}

func GetHostsAsHTML(db *gorm.DB) (hostsHtml template.HTML, err error) {
	ctx := context.Background()

	hosts, err := gorm.G[Host](db).Preload("Task", nil).Find(ctx)
	if err != nil {
		return "", err
	}

	var html string

	for _, u := range hosts {
		createdAt := u.CreatedAt.Format("2006-01-02 15:04:05")
		if u.TaskID == nil {
			html += fmt.Sprintf(`<tr>
				<td><a href="/hosts/edit/%d">%s</a></td>
				<td class="text-secondary">%s</td>
				<td class="text-secondary">%s</td>
				<td class="text-secondary">N/A</td>
			`, u.ID, u.Name, u.Mac, createdAt)
		} else {
			html += fmt.Sprintf(`<tr>
				<td><a href="/hosts/edit/%d">%s</a></td>
				<td class="text-secondary">%s</td>
				<td class="text-secondary">%s</td>
				<td class="text-secondary"><a href="/tasks/edit/%d">%s</a></td>
			`, u.ID, u.Name, u.Mac, createdAt, u.Task.ID, u.Task.Name)
		}
	}

	return template.HTML(html), err
}

func GetTasksAsHTML(db *gorm.DB) (tasksHtml template.HTML, err error) {
	ctx := context.Background()

	tasks, err := gorm.G[Task](db).Find(ctx)
	if err != nil {
		return "", err
	}

	var html string

	for _, u := range tasks {
		createdAt := u.CreatedAt.Format("2006-01-02 15:04:05")
		html += fmt.Sprintf(`<tr>
			<td><a href="/tasks/edit/%d">%s</a></td>
			<td class="text-secondary">%s</td>
		`, u.ID, u.Name, createdAt)
	}

	return template.HTML(html), err
}

func GetWifiKeysAsHTML(db *gorm.DB) (wifiHtml template.HTML, err error) {
	ctx := context.Background()

	var usedKeys []*WifiKey
	if err := db.
		Joins("LEFT JOIN hosts ON hosts.wifi_key_id = wifi_keys.id").
		Where("hosts.id IS NOT NULL").
		Find(&usedKeys).Error; err != nil {
		return "", err
	}

	usedMap := make(map[uint]bool, len(usedKeys))
	for _, k := range usedKeys {
		usedMap[k.ID] = true
	}

	keys, err := gorm.G[WifiKey](db).Find(ctx)
	if err != nil {
		return "", err
	}

	var html string
	for _, u := range keys {
		createdAt := u.CreatedAt.Format("2006-01-02 15:04:05")
		usedCol := ""
		if usedMap[u.ID] {
			usedCol = `<svg  xmlns="http://www.w3.org/2000/svg"  width="24"  height="24"  viewBox="0 0 24 24"  fill="none"  stroke="currentColor"  stroke-width="2"  stroke-linecap="round"  stroke-linejoin="round"  class="icon icon-tabler icons-tabler-outline icon-tabler-check"><path stroke="none" d="M0 0h24v24H0z" fill="none"/><path d="M5 12l5 5l10 -10" /></svg>`
		}

		html += fmt.Sprintf(`<tr>
			<td><a href="/wifikeys/edit/%d">%d</a></td>
			<td class="text-secondary">%s</td>
			<td class="text-secondary">%s</td>
		</tr>`, u.ID, u.ID, createdAt, usedCol)
	}

	return template.HTML(html), nil
}

func GetUnassignedWifiKeyCount(db *gorm.DB) (int, error) {
	var keys []*WifiKey
	db.
		Joins("LEFT JOIN hosts ON hosts.wifi_key_id = wifi_keys.id").
		Where("hosts.id IS NULL").
		Find(&keys)

	return len(keys), nil
}
