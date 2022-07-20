package service

import (
	"net/http"
	"time"

	"github.com/awesome-sphere/as-booking/db/models"
	"github.com/awesome-sphere/as-booking/serializer"
	"github.com/gin-gonic/gin"
)

type void struct{}

var member void

func GetTheaterID(c *gin.Context) (bool, []models.Theater) {
	var querySet []models.Theater
	err := models.DB.Model(&models.Theater{}).Find(&querySet).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return false, querySet
	}
	return true, querySet
}

func GetTimeSlots(c *gin.Context, theater_id int, movie_id int) (bool, []time.Time, []serializer.TimeSlotOutputSerializer) {
	var querySet []models.TimeSlot
	var filteredColumnQurySet []serializer.TimeSlotOutputSerializer
	date_slots := make(map[time.Time]void)
	date_slots_list := []time.Time{}
	err := models.DB.Model(&models.TimeSlot{}).Where(
		"theater_id", theater_id,
	).Find(&querySet).Where(
		"movie_id", movie_id,
	).Find(&filteredColumnQurySet).Error

	for _, time_slot := range filteredColumnQurySet {
		new_time_slot := time.Date(time_slot.Time.Year(), time_slot.Time.Month(), time_slot.Time.Day(), 0, 0, 0, 0, time.UTC)
		date_slots[new_time_slot] = member
	}
	for date, _ := range date_slots {
		date_slots_list = append(date_slots_list, date)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return false, date_slots_list, []serializer.TimeSlotOutputSerializer{}
	}
	return true, date_slots_list, filteredColumnQurySet
}

func GetAllTimeSlot(c *gin.Context) {
	if !models.DONE_SEEDING {
		c.JSON(http.StatusAccepted, gin.H{
			"message": "Service is not ready, please try later",
		})
		return
	}
	var inputSerializer serializer.TimeSlotInputSerializer
	if err := c.BindJSON(&inputSerializer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	status, all_theaters := GetTheaterID(c)
	date_slots := []time.Time{}
	var does_added bool
	if status {
		// theater id => time slot[]
		all_time_slots := []serializer.TimeSlotOutputSerializer{}
		for _, theater := range all_theaters {
			status, date_slot, time_slot_for_theater := GetTimeSlots(c, int(theater.ID), inputSerializer.MovieID)
			if !status {
				return
			}
			if !does_added {
				date_slots = date_slot
			}
			all_time_slots = append(all_time_slots, time_slot_for_theater...)

		}
		c.JSON(http.StatusOK, gin.H{
			"date_slot":  date_slots,
			"time_slots": all_time_slots,
		})
		return

	}
}
