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

	// For production only, drops the preexisting tables and creates new ones
	db.DB.Migrator().DropTable(&models.Device{}, &models.CurrentEvents{}, &models.Attendance{})
	err := db.DB.AutoMigrate(&models.Device{}, &models.CurrentEvents{}, &models.Attendance{})

	if err != nil {
		log.Fatal("Failed to migrate the database", err)
	}
	Seeder()
	log.Println("Database migrated successfully")

	r := gin.Default()
	// routes.RegisterDeviceRoutes(r)

	// Routes

	r.POST("/create-event", controllers.InitaiteEvent)
	r.GET("/events", controllers.GetEvents)
	r.GET("/event/:id", controllers.GetDevicesInEvent)
	r.DELETE("/delete-event/:id", controllers.DeleteEvent)
	r.GET("/event/graph/:id", controllers.GetEventGraph)

	r.POST("/device/participate", controllers.Participate)
	r.POST("/device/proximity-update", controllers.ProximityUpdate)

	for _, route := range r.Routes() {
		fmt.Println(route.Method, route.Path)
	}
	log.Println("Server started at :8080")
	r.Run(":8080")
}
