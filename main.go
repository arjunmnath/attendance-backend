package main

import (
	"attendance-backend/db"
	"attendance-backend/models"
	"attendance-backend/routes"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()
	db.ConnectDatabase()
	// Seeder()

	err := db.DB.AutoMigrate(&models.Device{})

	if err != nil {
		log.Fatal("Failed to migrate the database", err)
	}
	log.Println("Database migrated successfully")

	r := gin.Default()
	routes.RegisterDeviceRoutes(r)

	log.Println("Server started at :8080")
	r.Run(":8080")
}
