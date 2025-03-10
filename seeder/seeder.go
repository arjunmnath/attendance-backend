package seeder

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
		{UserID: 104, DeviceID: 5004, CreatedAt: time.Now()},
		{UserID: 105, DeviceID: 5005, CreatedAt: time.Now()},
		{UserID: 106, DeviceID: 5006, CreatedAt: time.Now()},
		{UserID: 107, DeviceID: 5007, CreatedAt: time.Now()},
		{UserID: 108, DeviceID: 5008, CreatedAt: time.Now()},
		{UserID: 109, DeviceID: 5009, CreatedAt: time.Now()},
		{UserID: 110, DeviceID: 5010, CreatedAt: time.Now()},
	}
	for _, device := range devices {
		err := db.DB.Create(&device).Error
		if err != nil {
			log.Println("Failed to insert:", err)
		}
	}
	log.Println("Dummy data inserted successfully!")
}