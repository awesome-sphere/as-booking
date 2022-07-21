package db

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/awesome-sphere/as-booking/db/models"
	"github.com/awesome-sphere/as-booking/kafka/interfaces"
	"github.com/awesome-sphere/as-booking/serializer"
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
	var seatType models.SeatType

	var value interfaces.BookingWriterInterface
	err := json.Unmarshal(message, &value)

	if err != nil {
		log.Fatalf("Failed to unmarshal message: %v", err.Error())
		return
	}

	status := "BOOKED"
	totalPrice := 0

	for _, seatNum := range value.SeatNumber {

		var timeSlotQuerySet []models.SeatInfo
		var seatNumQuerySet models.SeatInfo

		if err := models.DB.Model(&seatInfo).Where(
			"theater_id", value.TheaterID,
		).Find(
			&timeSlotQuerySet, "time_slot_id = ?", value.TimeSlotID,
		).Find(
			&seatNumQuerySet, "seat_number = ?", seatNum,
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

		if err := models.DB.Find(&seatType, "id", seatNumQuerySet.SeatTypeID).Error; err != nil {
			log.Fatal(err.Error())
			return
		}

		totalPrice += seatType.Price

		updateRedisStatus(value.TheaterID, value.TimeSlotID, seatNum, status)
	}
	updatePaymentOrder(value.UserID, value.TheaterID, value.TimeSlotID, value.SeatNumber, totalPrice, true)
}

func updateCancelingStatus(message []byte) {
	var seatInfo models.SeatInfo
	var seatType models.SeatType

	var value interfaces.CancelingWriterInterface
	err := json.Unmarshal(message, &value)

	if err != nil {
		log.Fatalf("Failed to unmarshal message: %v", err.Error())
		return
	}

	status := "AVAILABLE"
	totalPrice := 0

	for _, seatNum := range value.SeatNumber {

		var timeSlotQuerySet []models.SeatInfo
		var seatNumQuerySet models.SeatInfo

		if err := models.DB.Model(&seatInfo).Where(
			"theater_id", value.TheaterID,
		).Find(
			&timeSlotQuerySet, "time_slot_id = ?", value.TimeSlotID,
		).Find(
			&seatNumQuerySet, "seat_number = ?", seatNum,
		).Updates(
			models.SeatInfo{
				Status:     models.SeatStatus(status),
				BookedTime: time.Now(),
				BookedBy:   -1,
			},
		).Error; err != nil {
			log.Fatal(err.Error())
			return
		}

		if err := models.DB.Find(&seatType, "id", seatNumQuerySet.SeatTypeID).Error; err != nil {
			log.Fatal(err.Error())
			return
		}

		totalPrice += seatType.Price

		updateRedisStatus(value.TheaterID, value.TimeSlotID, seatNum, status)
	}
	updatePaymentOrder(value.UserID, value.TheaterID, value.TimeSlotID, value.SeatNumber, totalPrice, false)
}

func updateRedisStatus(theaterID int, timeSlotID int, seatNum int, status string) {
	url := "http://localhost:9004/seating/update-status"

	input := serializer.UpdateStatusSerializer{
		TimeSlotID: timeSlotID,
		TheaterID:  theaterID,
		SeatID:     seatNum,
		Status:     status,
	}
	json, err := json.Marshal(input)

	if err != nil {
		log.Fatal(err.Error())
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(json))

	if err != nil {
		log.Fatal(err.Error())
		return
	}

	defer resp.Body.Close()
}

func updatePaymentOrder(userID int, theaterID int, timeSlotID int, seatNum []int, price int, order bool) {
	input := serializer.UpdatePaymentSerializer{
		UserID:     userID,
		TheaterID:  theaterID,
		SeatID:     seatNum,
		TimeSlotId: timeSlotID,
		Price:      price,
	}

	json, err := json.Marshal(input)

	if err != nil {
		log.Fatal(err.Error())
		return
	}

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
	}

	defer resp.Body.Close()
}
