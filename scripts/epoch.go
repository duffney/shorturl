package main

import (
	"fmt"
	"time"
)

func main() {
	// Go release day: November 10, 2009
	goReleaseDate := time.Date(2009, time.November, 10, 0, 0, 0, 0, time.UTC)

	// Current time
	now := time.Now()

	// Calculate the duration between Go release date and current time
	duration := now.Sub(goReleaseDate)

	// Convert the duration to milliseconds
	milliseconds := duration.Milliseconds()

	fmt.Println("Milliseconds since Go release day:", milliseconds)
}
