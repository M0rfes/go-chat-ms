package models

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	ID        uint      `gorm:"primaryKey"`
	From      string    `gorm:"not null"`
	Content   string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
