package main

import (
	"attendance-backend/db"
	"attendance-backend/models"
	"log"
	"time"
)

func Seeder() {
	devices := []models.Device{
		{UserID: 101, DeviceID: 5001, CreatedAt: time.Now()},
		{UserID: 102, DeviceID: 5002, CreatedAt: time.Now()},
		{UserID: 103, DeviceID: 5003, CreatedAt: time.Now()},
	}
	for _, device := range devices {
		err := db.DB.Create(&device).Error
		if err != nil {
			log.Println("Failed to insert:", err)
		}
	}
	log.Println("Dummy data inserted successfully!")
}
