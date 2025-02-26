package routes

import (
	"attendance-backend/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterDeviceRoutes(router *gin.Engine){
	deviceGroup:=router.Group("/device")
	{
		deviceGroup.GET("/:id",controllers.GetDevice)

	}
}