package db

import (
	"context"

	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	ID     int
	Name   string
	Script string `gorm:"type:longtext"`
}

func CreateTask(name, script string, db *gorm.DB) error {
	ctx := context.Background()

	err := gorm.G[Task](db).Create(ctx, &Task{Name: name, Script: script})
	if err != nil {
		return err
	}

	return nil
}

func EditTask(name, script, id string, db *gorm.DB) error {
	ctx := context.Background()

	_, err := gorm.G[Task](db).Where("id = ?", id).Updates(ctx, Task{Name: name, Script: script})
	if err != nil {
		return err
	}

	return nil
}

func DeleteTask(id string, db *gorm.DB) error {
	ctx := context.Background()

	_, err := gorm.G[Task](db).Where("id = ?", id).Delete(ctx)
	if err != nil {
		return err
	}

	return nil
}

func GetTaskByID(id string, db *gorm.DB) (*Task, error) {
	ctx := context.Background()

	task, err := gorm.G[Task](db).Where("id = ?", id).First(ctx)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func GetTasks(db *gorm.DB) ([]Task, error) {
	ctx := context.Background()

	tasks, err := gorm.G[Task](db).Find(ctx)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}
