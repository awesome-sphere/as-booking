package serializer

type CheckSeatInputSerializer struct {
	TheaterID int   `json:"theater_id" binding:"required"`
	SeatID    int64 `json:"seat_id" binding:"required"`
}
