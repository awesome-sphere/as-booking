package serializer

type SeatsInputSerializer struct {
	TheaterID  int `json:"theater_id" binding:"required"`
	TimeSlotId int `json:"time_slot_id" binding:"required"`
}
