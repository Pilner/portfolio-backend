package main

import (
	"fmt"
	"time"
)

func main() {
	// Load Timezone to calibrate server time
	loadTimezone()
}

func loadTimezone() {
	// Set timezone
	loc, err := time.LoadLocation("Asia/Manila")
	if err != nil {
		fmt.Println("Error loading time location", err)
	} else {
		time.Local = loc
	}

	currentTime := time.Now()
	fmt.Println("Loaded Timezone: ", time.Local.String())
	fmt.Println("Current Server Time: ", currentTime.Format(time.RFC3339))
}
