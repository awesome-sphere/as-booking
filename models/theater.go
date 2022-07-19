package models

// type Seat struct {
// 	Number    int64    `json:"seat_number" gorm:"primaryKey;autoincrement;not null"`
// 	Type      SeatType `json:"seat_type" gorm:"foreignKey:TypeID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
// 	TypeID    int      `json:"seat_type_id" gorm:"not null"`
// 	Theater   Theater  `json:"theater" gorm:"foreignKey:TheaterID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
// 	TheaterID int      `json:"theater_id" gorm:"not null"`
// }

type Theater struct {
	ID       int64  `json:"id" gorm:"primaryKey;autoincrement;not null"`
	Location string `json:"location"`
}
