package controllers

import (
	"attendance-backend/db"
	"attendance-backend/models"
	"attendance-backend/utils"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type InitateEventInput struct {
	StartTime time.Time `json:"start_time" binding:"required"`
	EndTime   time.Time `json:"end_time" binding:"required"`
	Location  string    `json:"location" binding:"required"`
}

func InitaiteEvent(c *gin.Context) {
	var input InitateEventInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if alreadyExists := db.DB.First(&models.CurrentEvents{}, "start_time=? AND end_time=? AND location=?", input.StartTime, input.EndTime, input.Location); alreadyExists.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Event already exists"})
		return
	}

	event := models.CurrentEvents{
		ID:        uuid.New(),
		StartTime: input.StartTime,
		EndTime:   input.EndTime,
		Location:  input.Location,
	}
	maxPolls := int(event.EndTime.Sub(event.StartTime).Minutes() / 2)
	if maxPolls <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event duration"})
		return
	}

	// Create CurrentPolling entry for the event
	utils.CurrentPolling[event.ID] = [2]int{0, maxPolls}
	log.Println(utils.CurrentPolling[event.ID])
	// Initialize the graph for the event
	utils.InitializeGraph(event.ID)

	if err := db.DB.Create(&event).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create event"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Event created successfully", "event_id": event.ID})

}

func GetEvents(c *gin.Context) {
	var events []models.CurrentEvents

	if err := db.DB.Find(&events).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch events"})
		return
	}
	c.JSON(http.StatusOK, events)
}

func DeleteEvent(c *gin.Context) {
	id := c.Param("id")
	result := db.DB.Delete(&models.CurrentEvents{}, "id=?", id)

	if result.Error != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete event"})
		return
	} else if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return

	} else {
		// Parse the UUID, because delete expects a UUID
		uid, err := uuid.Parse(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
			return
		}
		delete(utils.CurrentPolling, uid)
		delete(utils.AttendanceGraph, uid)
		delete(utils.GraphMutex, uid)

		c.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully"})

	}

}
