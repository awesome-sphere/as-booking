package serializer

type InputSerializer struct {
	TheaterID  int   `json:"theater_id" binding:"required"`
	SeatID     []int `json:"seat_id" binding:"required"`
	TimeSlotID int   `json:"time_slot_id" binding:"required"`
}
