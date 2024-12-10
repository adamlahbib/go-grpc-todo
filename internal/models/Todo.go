package models

import (
	"time"

	"gorm.io/gorm"
)

type Todo struct {
	gorm.Model

	Id          int `gorm:"primaryKey,autoIncrement"`
	Title       string
	Description string
	Deadline    time.Time
}
