package serializer

import "time"

type TimeSlotInputSerializer struct {
	MovieID int `json:"movie_id" binding:"required"`
}

type TimeSlotOutputSerializer struct {
	ID        int64     `json:"id" gorm:"primaryKey;autoincrement;not null"`
	MovieID   int       `json:"movie_id" gorm:"not null"`
	Time      time.Time `json:"time" gorm:"not null"`
	TheaterID int       `json:"theater_id" gorm:"not null"`
}
