package models

import (
	"time"
)

type TimeSlot struct {
	// gorm.Model
	ID        int64     `json:"id" gorm:"primaryKey;autoincrement;not null"`
	MovieID   int       `json:"movie_id" gorm:"not null"`
	Time      time.Time `json:"time" gorm:"not null"`
	TheaterID int       `json:"theater_id" gorm:"not null"`
	Theater   Theater   `gorm:"foreignKey:theater_id; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
