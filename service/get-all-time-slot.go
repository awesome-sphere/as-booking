package service

import (
	"net/http"

	"github.com/awesome-sphere/as-booking/models"
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

func GetTimeSlots(theater_id int) (bool, []models.TimeSlot) {
	var querySet []models.TimeSlot
	err := models.DB.Model(&models.TimeSlot{}).Where("theater_id", theater_id).Find(&querySet).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return false, querySet
	}
	return true, querySet
}

func GetAllTimeSlot(c *gin.Context) {
	if !models.DONE_SEEDING {
		c.JSON(http.StatusAccepted, gin.H{
			"message": "Service is not ready, please try later",
		})
		return
	}
	status, all_theaters := GetTheaterID(c)
	if status {
		// theater id => time slot[]
		all_time_slots := make(map[int][]models.TimeSlot)
		for _, theater := range all_theaters {
			status, time_slot_for_theater := GetTimeSlots(int(theater.ID))
			if !status {
				return
			}

		}
		return

	}
}
