package service

import (
	"net/http"

	"github.com/awesome-sphere/as-booking/db/models"
	"github.com/awesome-sphere/as-booking/serializer"
	"github.com/gin-gonic/gin"
)

func BuySeat(c *gin.Context) {
	if !models.DONE_SEEDING {
		c.JSON(http.StatusAccepted, gin.H{
			"message": "Service is not ready, please try later",
		})
		return
	}
	var inputSerializer serializer.InputSerializer
	if err := c.BindJSON(&inputSerializer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var seatInfo models.SeatInfo
	var timeSlot models.TimeSlot

	var theater models.Theater
	var locationResult models.SeatInfo
	var timeSlotResult models.TimeSlot

	status := "BOUGHT"

	for _, seat := range inputSerializer.SeatID {

		var timeSlotQuerySet []models.SeatInfo
		var seatNumQuerySet models.SeatInfo

		if err := models.DB.Model(&seatInfo).Where(
			"theater_id", inputSerializer.TheaterID,
		).Find(
			&timeSlotQuerySet, "time_slot_id = ?", inputSerializer.TimeSlotID,
		).Find(
			&seatNumQuerySet, "seat_number = ?", seat,
		).Updates(
			models.SeatInfo{
				Status: models.SeatStatus(status),
			},
		).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "Failed with error",
				"error":  err.Error(),
			})
			return
		}

		locationResult = seatNumQuerySet

		var querySet models.TimeSlot

		if err := models.DB.Model(&timeSlot).Where(
			"theater_id", inputSerializer.TheaterID,
		).Find(
			&querySet, "id = ?", int64(inputSerializer.TimeSlotID),
		).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "Failed with error",
				"error":  err.Error(),
			})
			return
		}

		timeSlotResult = querySet
	}

	if err := models.DB.Find(&theater, "id", locationResult.TheaterID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Failed with error",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"location": theater.Location,
		"date":     timeSlotResult.Time,
		"movie_id": timeSlotResult.MovieID,
	})
}
