package calendar

import (
	"container/list"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type GCalendar struct {
}

// Get the end of the day (23:59) for the specified time
func LastTimeOfDay(t time.Time) time.Time {
	y, m, d := t.Date()

	return time.Date(y, m, d, 23, 59, 0, 0, t.Location())
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}

	return config.Client(context.Background(), tok)
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// Return events that exist on a calendar
func (gcalendar *GCalendar) GetEvents() *list.List {
	ctx := context.Background()
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	// Create new service to access GCalendar API
	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	//t := time.Now().Format(time.RFC3339)
	//events, err := srv.Events.List("primary").ShowDeleted(false).
	//	SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
	calendars, err := srv.CalendarList.List().Do()

	allEvents := list.New()

	startDate := LastTimeOfDay(time.Now()).Format(time.RFC3339)
	endDate := time.Now().Truncate(24*time.Hour).AddDate(0, 0, 1).Format(time.RFC3339)
	for _, calendar := range calendars.Items {
		events, err := srv.Events.List(calendar.Id).TimeMin(startDate).TimeMax(endDate).OrderBy("startTime").Do()
		if err != nil {
			log.Fatalf("Unable to retrieve next of the user's events: %v", err)
		}

		if len(events.Items) > 0 {
			for _, event := range events.Items {
				allEvents.PushBack(event)
			}

		} else {
			fmt.Println("No upcoming events found.")
		}
	}

	return allEvents
}

func (gcalendar *GCalendar) GetLocations(events *list.List) []string {
	type void struct{}
	var member void
	unrepeatedLocations := make(map[string]void)

	for e := events.Front(); e != nil; e = e.Next() {
		event := e.Value.(calendar.Event)

		if event.Location != "" {
			unrepeatedLocations[event.Location] = member
		}
	}

	var locations = make([]string, 0)
	for city := range unrepeatedLocations {
		locations = append(locations, city)
	}

	return locations
}
