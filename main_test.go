package main

import (
	"attendance-backend/utils"
	"fmt"
	"testing"
)

func TestSelectDevices(t *testing.T) {

	devices := map[int]int{
		5001: 0,
		5002: 0,
		5003: 0,
		5004: 0,
		5005: 0,
		5006: 0,
		5007: 0,
	}

	for i := range 4 {
		selected := utils.SelectAggregators(&devices, 3)
		fmt.Printf("Round %d Selected devices: %v\n", i, selected)

	}
	fmt.Println(devices)

}
