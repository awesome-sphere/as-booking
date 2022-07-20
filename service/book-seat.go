package service

import (
	"net/http"

	"github.com/awesome-sphere/as-booking/jwt"
	"github.com/awesome-sphere/as-booking/kafka"
	"github.com/awesome-sphere/as-booking/kafka/interfaces"
	"github.com/awesome-sphere/as-booking/models"
	"github.com/awesome-sphere/as-booking/serializer"
	"github.com/gin-gonic/gin"
)

func BookSeat(c *gin.Context) {
	if !models.DONE_SEEDING {
		c.JSON(http.StatusAccepted, gin.H{
			"message": "Service is not ready, please try later",
		})
		return
	}
	isValid, claimedToken := jwt.IsValidJWT(c)
	if isValid {
		var inputSerializer serializer.InputSerializer
		if err := c.BindJSON(&inputSerializer); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userID, ok := claimedToken["user_id"].(float64)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Error converting User ID (User ID is not given or wrong type)",
			})
			return
		}
		message := &interfaces.BookingWriterInterface{
			UserID:     int(userID),
			TimeSlotID: inputSerializer.TimeSlotID,
			TheaterID:  inputSerializer.TheaterID,
			SeatNumber: inputSerializer.SeatID,
		}

		result := make(chan kafka.Result)
		go kafka.ProduceBooking(message, result)
		r := <-result
		if r.IsCompleted {
			c.JSON(http.StatusOK, gin.H{
				"status": "Submitted",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Failed with error",
			"error":  r.Err.Error(),
		})
	}
}
