package serializer

type UpdateStatusSerializer struct {
	TimeSlotID int    `json:"time_slot_id"`
	TheaterID  int    `json:"theater_id"`
	SeatID     int    `json:"seat_id"`
	Status     string `json:"status"`
}
