package calendar

import "container/list"

type CalendarFactory interface {
	GetEvents() *list.List
	GetLocations(events *list.List) []string
}

func GetCalendarFactory(calendarType string) CalendarFactory {
	switch calendarType {
	case "gcalendar":
		return &GCalendar{}
	}

	return nil
}
