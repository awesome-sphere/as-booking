package models

import "time"

type SeatStatus string

const (
	Available SeatStatus = "AVAILABLE"
	Booked    SeatStatus = "BOOKED"
	Bought    SeatStatus = "BOUGHT"
)

type SeatInfo struct {
	ID         int64 `json:"id" gorm:"primaryKey;autoincrement;not null;"`
	TimeSlotId int   `json:"time_slot_id" gorm:"not null;"`
	// TimeSlot   TimeSlot   `gorm:"foreignKey:time_slot_id; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TheaterId  int        `json:"theater_id" gorm:"not null"`
	Theater    Theater    `gorm:"foreignKey:theater_id; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	SeatNumber int        `json:"seat_number" gorm:"not null"`
	SeatTypeId int        `json:"seat_type_id" gorm:"not null"`
	SeatType   SeatType   `gorm:"foreignKey:seat_type_id; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Status     SeatStatus `json:"status" gorm:"default:'AVAILABLE'" sql:"type:seat_status"`
	BookedTime time.Time  `json:"booked_time" gorm:"not null"`
	BookedBy   int        `json:"booked_by" gorm:"not null"`
}
