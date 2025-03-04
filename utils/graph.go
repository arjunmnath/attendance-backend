package utils

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EventGraph struct {
	Nodes map[int]map[int]int

	/*
		1: {2: 1, 3: 1}
		2: {1: 1}
		3: {1: 1}
	*/
}

var (
	GraphMutex      = make(map[uuid.UUID]*sync.Mutex)
	AttendanceGraph = make(map[uuid.UUID]*EventGraph)
	Polling         = make(map[uuid.UUID][]int)
)

func InitializeGraph(eventID uuid.UUID) {
	if _, exists := GraphMutex[eventID]; !exists {
		GraphMutex[eventID] = &sync.Mutex{}
	}

	GraphMutex[eventID].Lock()
	defer GraphMutex[eventID].Unlock()

	if _, exists := AttendanceGraph[eventID]; !exists {
		AttendanceGraph[eventID] = &EventGraph{
			Nodes: make(map[int]map[int]int),
		}

		log.Println("Graph initialized")
	}

}

func AddEdge(c *gin.Context, eventID uuid.UUID, source, destination int) {

	mutex, mutexExists := GraphMutex[eventID]
	if !mutexExists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Event mutex not initialized"})
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	graph, graphExists := AttendanceGraph[eventID]

	if !graphExists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Event does not exist"})
		return
	}

	if _, exists := graph.Nodes[source]; !exists {
		graph.Nodes[source] = make(map[int]int)
	}
	if _, exists := graph.Nodes[destination]; !exists {
		graph.Nodes[destination] = make(map[int]int)
	}
	log.Println("Edge added between", source, "and", destination)

	graph.Nodes[source][destination]++
	graph.Nodes[destination][source]++

}

func StartEventPolling(eventID uuid.UUID) {
	ticker := time.NewTicker(2 * time.Minute)
	defer ticker.Stop()
	fmt.Println("Polling started")

	for {
		if Polling[eventID][0] > Polling[eventID][1] {
			fmt.Println("Polling completed")
			return
		}
		fmt.Println("Incrementing poll count")
		Polling[eventID][0]++
		<-ticker.C

		fmt.Println("Polling", Polling[eventID][0], "out of", Polling[eventID][1])

	}
}
