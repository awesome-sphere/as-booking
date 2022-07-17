package service

import (
	"net/http"

	"github.com/awesome-sphere/as-booking/jwt"
	"github.com/awesome-sphere/as-booking/kafka"
	"github.com/awesome-sphere/as-booking/kafka/writer_interface"
	"github.com/awesome-sphere/as-booking/serializer"
	"github.com/gin-gonic/gin"
)

func BookSeat(c *gin.Context) {
	is_valid, claimed_token := jwt.IsValidJWT(c)
	if is_valid {
		var input_serializer serializer.InputSerializer
		if err := c.BindJSON(&input_serializer); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		user_id, ok := claimed_token["user_id"].(float64)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Error converting User ID (User ID is not given or wrong type)",
			})
			return
		}
		kafka_message := &writer_interface.BookingWriterInterface{
			UserID:     int(user_id),
			TimeSlotId: input_serializer.TimeSlotId,
			TheaterId:  input_serializer.TheaterID,
			SeatNumber: input_serializer.SeatID,
		}

		result := make(chan kafka.Result)
		go kafka.Produce(kafka_message, result)
		r := <-result
		if r.IsCompleted {
			c.JSON(http.StatusOK, gin.H{
				"status": "Submitted",
			})
			return
		}
		kafka.Consume(kafka.BOOKING_TOPIC)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Failed with error",
			"error":  r.Err.Error(),
		})
	}
}
