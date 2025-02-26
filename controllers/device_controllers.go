package controllers

import (
	"attendance-backend/db"
	"attendance-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetDevice(c *gin.Context) {
	id := c.Param("id")
	var device models.Device

	if err := db.DB.First(&device, "id=?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}
	c.JSON(http.StatusOK, device)
}
