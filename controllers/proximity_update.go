package controllers

import (
	"attendance-backend/db"
	"attendance-backend/models"
	"attendance-backend/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProximityUpdateInput struct {
	Source  int            `json:"source" binding:"required"`
	EventID uuid.UUID      `json:"event_id" binding:"required"`
	Devices map[string]int `json:"devices" binding:"required"`
}

func ProximityUpdate(c *gin.Context) {

	// Validate input
	var input ProximityUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "message": "Invalid input"})
		return
	}

	// This transaction is not reverting the changes made in the graph
	// This is a bug that needs to be fixed
	// The graph should be updated only if the transaction is successful

	var edges []struct {
		Source      int
		Destination int
	}

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		var attendance models.Attendance

		// Check if device is participating in the event
		if err := tx.First(&attendance, "device_id=? AND event_id=?", input.Source, input.EventID).Error; err != nil {
			return fmt.Errorf("device is not participating in the event")
		}

		// Check poll count limit
		if attendance.PollCount > utils.Polling[input.EventID][0] {
			return fmt.Errorf("poll count exceeded, attendance not updated")
		}

		// Update proximity score for each detected device
		for deviceStr, distance := range input.Devices {
			// Convert device ID from string to int
			device, conversionError := strconv.Atoi(deviceStr)
			if conversionError != nil {
				return fmt.Errorf("invalid device ID")
			}

			// Skip devices with RSSI less than -100
			if distance < -100 {
				continue
			}

			var targetAttendance models.Attendance

			// Check if detected device is participating in the event
			if err := tx.First(&targetAttendance, "device_id=? AND event_id=?", device, input.EventID).Error; err != nil {
				continue
			}

			// Update proximity score for the detected device
			if err := tx.Model(&models.Attendance{}).
				Where("device_id=? AND event_id=?", device, input.EventID).
				Update("proximity_score", gorm.Expr("proximity_score + ?", 1)).Error; err != nil {
				return fmt.Errorf("failed to update attendance for device %d", device)
			}

			// Add edge between the devices in the memory graph
			// utils.AddEdge(c, input.EventID, input.Source, device)

			// Store edge in temporary list
			edges = append(edges, struct {
				Source      int
				Destination int
			}{input.Source, device})

			// Update the poll count
			if err := tx.Model(&models.Attendance{}).
				Where("id=?", attendance.ID).
				Update("poll_count", gorm.Expr("poll_count + ?", 1)).Error; err != nil {
				return fmt.Errorf("failed to update poll count")
			}
		}

		return nil

	})

	// Handle transaction result
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, edge := range edges {
		utils.AddEdge(c, input.EventID, edge.Source, edge.Destination)

	}
	c.JSON(http.StatusOK, gin.H{"message": "Proximity updated successfully"})

}

// old code without transaction
/*



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
		if attendance.PollCount > utils.Polling[input.EventID][0] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Poll count exceeded, attendance not updated"})
			return
		}

		// Update proximity score for each detected device
		for deviceStr, distance := range input.Devices {
			// Skip devices with RSSI less than -100

			device, conversionError := strconv.Atoi(deviceStr)
			if conversionError != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device ID "})
				return
			}
			if distance < -100 {
				continue
			}
			var targetAttendance models.Attendance

			// Check if device is participating in the event
			if err := db.DB.First(&targetAttendance, "device_id=? AND event_id=?", device, input.EventID).Error; err != nil {
				continue
			}

			// Update proximity score for the device
			err := db.DB.Model(&models.Attendance{}).Where("device_id=? AND event_id=?", device, input.EventID).Update("proximity_score", gorm.Expr("proximity_score + ?", 1)).Error

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update attendance"})
				return
			}

			// Add edge between the devices in the memory graph
			utils.AddEdge(c, input.EventID, input.Source, device)

			// update the poll count
			pcerr := db.DB.Model(&models.Attendance{}).Where("id=?", attendance.ID).Update("poll_count", gorm.Expr("poll_count + ?", 1)).Error
			if pcerr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update poll count"})
				return
			}
			fmt.Println(utils.AttendanceGraph[input.EventID].Nodes)

		}
		c.JSON(http.StatusOK, gin.H{"message": "Attendance updated successfully"})

	}



*/
