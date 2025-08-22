package template

import (
	"time"
)

// GetCurrentFormattedTime returns the current date and time in the format
func GetCurrentFormattedTime() string {
	currentTime := time.Now()
	return currentTime.Format("Monday, 02 Jan 2006 15:04")
}
