package controllers

import (
	"attendance-backend/db"
	"attendance-backend/models"
	"attendance-backend/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DeviceInput struct {
	DeviceID int       `json:"device_id" binding:"required"`
	EventID  uuid.UUID `json:"event_id" binding:"required"`
}

func GetDevice(deviceID int) (models.Device, error) {
	var device models.Device
	if err := db.DB.First(&device, "device_id=?", deviceID).Error; err != nil {
		return device, err
	}
	return device, nil
}

func Participate(c *gin.Context) {

	var input DeviceInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "Invalid input"})
		return
	}

	eventUUID := input.EventID

	if _, exists := utils.AttendanceGraph[eventUUID]; !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Event does not exist"})
		return
	}

	// Checks if device exists

	if _, err := GetDevice(input.DeviceID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Device does not exist"})
		return
	}

	// Checks if device is already participating in an event
	if alreadyExists := db.DB.First(&models.Attendance{}, "device_id=?", input.DeviceID); alreadyExists.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Device already participating in an event"})
		return
	}


	attendance := models.Attendance{
		ID:             uuid.New(),
		DeviceID:       input.DeviceID,
		EventID:        input.EventID,
		ProximityScore: 0,
	}

	if err := db.DB.Create(&attendance).Error; err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to participate in event"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Participation successful"})

}
