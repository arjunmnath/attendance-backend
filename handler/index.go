package handler

import (
	"attendance-backend/controllers"
	"attendance-backend/db"
	"attendance-backend/seeder"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// Declare the Gin engine globally
var Engine *gin.Engine

func init() {
	godotenv.Load()
	db.ConnectDatabase()

	// For production only, drops the preexisting tables and creates new ones
    /*
	db.DB.Migrator().DropTable(&models.Device{}, &models.CurrentEvents{}, &models.Attendance{})
	err := db.DB.AutoMigrate(&models.Device{}, &models.CurrentEvents{}, &models.Attendance{})

	if err != nil {
		log.Fatal("Failed to migrate the database", err)
	}
	log.Println("Database migrated successfully")
    */

	seeder.Seeder()
    Engine = gin.Default()
	// routes.RegisterDeviceRoutes(engine)

	// Routes
	Engine.POST("/create-event", controllers.InitaiteEvent)
	Engine.GET("/events", controllers.GetEvents)
	Engine.GET("/event/:id", controllers.GetDevicesInEvent)
	Engine.DELETE("/delete-event/:id", controllers.DeleteEvent)
	Engine.GET("/event/graph/:id", controllers.GetEventGraph)

	Engine.POST("/device/participate", controllers.Participate)
	Engine.POST("/device/proximity-update", controllers.ProximityUpdate)
	Engine.GET("/connectSSE/:event_id", controllers.ConnectSSE)
}

// **Exported Handler function for Vercel**
func Handler(w http.ResponseWriter, r *http.Request) {
	Engine.ServeHTTP(w, r)
}
