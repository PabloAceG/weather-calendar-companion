// First Go program
package main

import (
	calendarFactory "weather-companion/calendar"

	"fmt"
	"time"
)

// Main function
func main() {
	fmt.Println(time.Now().Format(time.RFC3339))
	fmt.Println(time.Now().Truncate(24 * time.Hour).Format(time.RFC3339))
	fmt.Println(time.Now().Truncate(24*time.Hour).AddDate(0, 0, 1).Format(time.RFC3339))

	calendar := calendarFactory.GetCalendarFactory("gcalendar")
	events := calendar.GetEvents()
	locations := calendar.GetLocations(events)

	fmt.Printf("%v", locations)

	// TODO: get weather
	// TODO: plot information
	// TODO: send telegram message

}
