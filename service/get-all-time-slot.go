package service

import (
	"net/http"

	"github.com/awesome-sphere/as-booking/models"
	"github.com/awesome-sphere/as-booking/serializer"
	"github.com/gin-gonic/gin"
)

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

func GetTimeSlots(c *gin.Context, theater_id int, movie_id int) (bool, []serializer.TimeSlotOutputSerializer) {
	var querySet []models.TimeSlot
	var filteredColumnQurySet []serializer.TimeSlotOutputSerializer
	err := models.DB.Model(&models.TimeSlot{}).Where(
		"theater_id", theater_id,
	).Find(&querySet).Where(
		"movie_id", movie_id,
	).Find(&filteredColumnQurySet).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return false, []serializer.TimeSlotOutputSerializer{}
	}
	return true, filteredColumnQurySet
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
	if status {
		// theater id => time slot[]
		all_time_slots := []serializer.TimeSlotOutputSerializer{}
		for _, theater := range all_theaters {
			status, time_slot_for_theater := GetTimeSlots(c, int(theater.ID), inputSerializer.MovieID)
			if !status {
				return
			}
			all_time_slots = append(all_time_slots, time_slot_for_theater...)

		}
		c.JSON(http.StatusOK, gin.H{
			"time_slots": all_time_slots,
		})
		return

	}
}
