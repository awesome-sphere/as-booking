package interfaces

type CancelingWriterInterface struct {
	UserID     int   `json:"user_id"`
	TimeSlotID int   `json:"time_slot_id"`
	TheaterID  int   `json:"theater_id"`
	SeatNumber []int `json:"seat_number"`
}
