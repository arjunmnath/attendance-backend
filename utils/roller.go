package utils

import (
	"container/heap"

	"github.com/google/uuid"
)

type DeviceHeap []Device

type Device struct {
	DeviceId    int
	SelectCount int
	HeapIndex   int
}

func (h DeviceHeap) Len() int { return len(h) }

func (h DeviceHeap) Less(i, j int) bool { return h[i].SelectCount < h[j].SelectCount }

func (h DeviceHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].HeapIndex = i
	h[j].HeapIndex = j
}

func (h *DeviceHeap) Push(x interface{}) {
	n := len(*h)
	device := x.(Device)
	device.HeapIndex = n
	*h = append(*h, device)
}

func (h *DeviceHeap) Pop() interface{} {
	old := *h
	n := len(old)
	device := old[n-1]
	*h = old[0 : n-1]
	return device
}

func SelectAggregators(devices map[int]int, numBroadcasters int) []int {
	h := &DeviceHeap{}
	heap.Init(h)

	// Add all devices to the heap
	for d, s := range devices {
		heap.Push(h, Device{DeviceId: d, SelectCount: s})
	}
	selected := []int{}

	for len(selected) < numBroadcasters && h.Len() > 0 {
		device := heap.Pop(h).(Device)
		selected = append(selected, device.DeviceId)
		(devices)[device.DeviceId]++
	}
	return selected

}

func InitializeEventDevices(eventID uuid.UUID) {
	if _, exists := EventDevices[eventID]; !exists {
		EventDevices[eventID] = DeviceCollection{
			Devices: make(map[int]int),
			Channel: make(map[int]chan string),
		}
	}

}
