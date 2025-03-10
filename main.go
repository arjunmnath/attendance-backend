package main

import (
    "attendance-backend/handler"
	"fmt"
	"log"
)

func main() {
	for _, route := range handler.Engine.Routes() {
		fmt.Println(route.Method, route.Path)
	}
	log.Println("Server started at :8080")
	handler.Engine.Run(":8080")
}
