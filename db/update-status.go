package db

import (
	"encoding/json"
	"log"
	"time"

	"github.com/awesome-sphere/as-booking/kafka/interfaces"
	"github.com/awesome-sphere/as-booking/models"
)

func UpdateStatus(topic string, message []byte) {
	if topic == "booking" {
		updateBookingStatus(message)
	} else if topic == "canceling" {
		updateCancelingStatus(message)
	}
}

func updateBookingStatus(message []byte) {
	var seatInfo models.SeatInfo
	var querySet []models.SeatInfo

	var value interfaces.BookingWriterInterface
	err := json.Unmarshal(message, &value)

	if err != nil {
		log.Fatalf("Failed to unmarshal message: %v", err.Error())
		return
	}

	for _, seatNum := range value.SeatNumber {

		if err := models.DB.Model(&seatInfo).Where(
			"theater_id", value.TheaterID,
		).Find(
			&querySet,
		).Where(
			"time_slot_id = ?", value.TimeSlotID,
		).Where(
			"seat_number = ?", seatNum,
		).Updates(
			models.SeatInfo{
				Status:     "BOOKED",
				BookedTime: time.Now(),
				BookedBy:   value.UserID,
			},
		).Error; err != nil {
			log.Fatal(err.Error())
			return
		}
	}
}

func updateCancelingStatus(message []byte) {
	var seatInfo models.SeatInfo
	var querySet []models.SeatInfo

	var value interfaces.CancelingWriterInterface
	err := json.Unmarshal(message, &value)

	if err != nil {
		log.Fatalf("Failed to unmarshal message: %v", err.Error())
		return
	}

	for _, seatNum := range value.SeatNumber {

		if err := models.DB.Model(&seatInfo).Where(
			"theater_id", value.TheaterID,
		).Find(
			&querySet,
		).Where(
			"time_slot_id = ?", value.TimeSlotID,
		).Where(
			"seat_number = ?", seatNum,
		).Updates(
			models.SeatInfo{
				Status:     "AVAILABLE",
				BookedTime: time.Now(),
				BookedBy:   0,
			},
		).Error; err != nil {
			log.Fatal(err.Error())
			return
		}
	}
}
