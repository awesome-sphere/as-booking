package db

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/awesome-sphere/as-booking/db/models"
	"github.com/awesome-sphere/as-booking/kafka/interfaces"
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

	totalPrice := 0

	for _, seatNum := range value.SeatNumber {

		status := "BOOKED"

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
				Status:     models.SeatStatus(status),
				BookedTime: time.Now(),
				BookedBy:   value.UserID,
			},
		).Error; err != nil {
			log.Fatal(err.Error())
			return
		}

		totalPrice += seatInfo.SeatType.Price

		updateRedisStatus(value.TheaterID, value.TimeSlotID, seatNum, status)
	}
	updatePaymentOrder(value.UserID, value.TheaterID, value.TimeSlotID, value.SeatNumber, totalPrice, true)
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

	totalPrice := 0

	for _, seatNum := range value.SeatNumber {

		status := "AVAILABLE"

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
				Status:     models.SeatStatus(status),
				BookedTime: time.Now(),
				BookedBy:   0,
			},
		).Error; err != nil {
			log.Fatal(err.Error())
			return
		}

		totalPrice += seatInfo.SeatType.Price

		updateRedisStatus(value.TheaterID, value.TimeSlotID, seatNum, status)
	}
	updatePaymentOrder(value.UserID, value.TheaterID, value.TimeSlotID, value.SeatNumber, totalPrice, false)
}

func updateRedisStatus(theaterID int, timeSlotID int, seatNum int, status string) {
	url := "http://localhost:9004/seating/update-status"

	jsonString := fmt.Sprintf(
		`{
			"theater_id": %d,
			"time_slot_id": %d,
			"seat_number": %d,
			"seat_status": "%s"
		}`, theaterID, timeSlotID, seatNum, status,
	)
	json := []byte(jsonString)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(json))

	if err != nil {
		log.Fatal(err.Error())
		return
	}

	defer resp.Body.Close()
}

func updatePaymentOrder(userID int, theaterID int, timeSlotID int, seatNum []int, price int, order bool) {

	jsonString := fmt.Sprintf(
		`{
			"user_id": %d
			"theater_id": %d,
			"time_slot_id": %d,
			"seat_number": %d,
			"price": %d
		}`, userID, theaterID, timeSlotID, seatNum, price,
	)
	json := []byte(jsonString)

	url := "http://localhost:9003/payment/"
	if order {
		url += "add-order"
	} else {
		url += "cancel-order"
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(json))

	if err != nil {
		log.Fatal(err.Error())
		return
	} else {
		log.Println(resp.Body)
	}

	defer resp.Body.Close()
}
