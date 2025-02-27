package controllers

import (
	"attendance-backend/db"
	"attendance-backend/models"
	"attendance-backend/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProximityUpdateInput struct {
	Source  int         `json:"source" binding:"required"`
	EventID uuid.UUID   `json:"event_id" binding:"required"`
	Devices map[int]int `json:"devices" binding:"required"`
}

func ProximityUpdate(c *gin.Context) {

	// Validate input
	var input ProximityUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "Invalid input"})
		return
	}

	var attendance models.Attendance

	// Check if device is participating in the event
	if err := db.DB.First(&attendance, "device_id=? AND event_id=?", input.Source, input.EventID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Device is not participating in the event"})
		return
	}

	// Update proximity score for each detected device
	for device, distance := range input.Devices {
		// Skip devices with RSSI less than -100
		if distance < -100 {
			continue
		}
		var targetAttendance models.Attendance

		// Check if device is participating in the event
		if err := db.DB.First(&targetAttendance, "device_id=? AND event_id=?", device, input.EventID).Error; err != nil {
			continue
		}

		// Update proximity score for the device
		err := db.DB.Model(&models.Attendance{}).Where("device_id=? AND event_id=?", device, input.EventID).Update("proximity_score ", gorm.Expr("proximity_score + ?", 1)).Error

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update attendance"})
			return
		}

		// Add edge between the devices in the memory graph
		utils.AddEdge(c, input.EventID, input.Source, device)
		fmt.Println(utils.AttendanceGraph[input.EventID].Nodes)

	}
	c.JSON(http.StatusOK, gin.H{"message": "Attendance updated successfully"})

}
