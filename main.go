package main

import (
	"attendance-backend/controllers"
	"attendance-backend/db"
	"attendance-backend/models"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()
	db.ConnectDatabase()
	// Seeder()

	err := db.DB.AutoMigrate(&models.Device{}, &models.CurrentEvents{}, &models.Attendance{})

	if err != nil {
		log.Fatal("Failed to migrate the database", err)
	}
	log.Println("Database migrated successfully")

	r := gin.Default()
	// routes.RegisterDeviceRoutes(r)

	// Routes

	r.POST("/create-event", controllers.InitaiteEvent)
	r.GET("/events", controllers.GetEvents)
	r.DELETE("/delete-event/:id", controllers.DeleteEvent)

	r.POST("/device/participate", controllers.Participate)

	for _, route := range r.Routes() {
		fmt.Println(route.Method, route.Path)
	}
	log.Println("Server started at :8080")
	r.Run(":8080")
}
