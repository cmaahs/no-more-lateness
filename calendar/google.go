package calendar

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	// "github.com/skratchdot/open-golang/open"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	calendar "google.golang.org/api/calendar/v3"
)

const googleCalendarDateTimeFormat = time.RFC3339

// https://www.google.com/url?q=https://teams.microsoft.com/l/meetup-join/19%253ameeting_MTU0YzY5MTgtNDZjNy00MjM4LTg5MzYtYTgxZDk0MzJjYzAx%2540thread.v2/0?context%3D%257b%2522Tid%2522%253a%2522e0793d39-0939-496d-b129-198edd916feb%2522%252c%2522Oid%2522%253a%25225a399f46-82c6-4606-84bf-a154d4b60548%2522%257d&sa=D&source=calendar&ust=1600608611212000&usg=AOvVaw2TYDQ2tagscOoJNxkpzKRw
var zoomURLRegexp = regexp.MustCompile(`https://.*?zoom\.us/(?:j/(\d+)|my/(\S+))`)
var teamsURLRegexp = regexp.MustCompile(`https://.*?teams.microsoft.com/.*`)
var webexURLRegexp = regexp.MustCompile(`https://.*?webex.com/.*j.php?.*>`)
var zoomURLRegexpPwd = regexp.MustCompile(`https://.*?zoom\.us/j/.*pwd=(.*)`)

// GoogleCal - Structure to hold stuff
type GoogleCal struct {
	Provider string
	Client   *http.Client
}

// GetClient - Apply DNS updates to Google DNS Hosted Zone
func (p *GoogleCal) GetClient() (bool, error) {

	usr, err := user.Current()
	if err != nil {
		fmt.Println("No current user")
		os.Exit(1)
	}

	directory := filepath.Join(usr.HomeDir, ".config", "google", "no-more-lateness.json")
	config, err := readGoogleClientConfigFromFile(directory)
	if err != nil && err.Error() != "oauth2/google: no credentials found" {
		fmt.Println(err)
		os.Exit(1)
	}

	client := getClient(config)

	p.Client = client

	return true, nil

}

// GetEvents - Return the next num events
func (p *GoogleCal) GetEvents(num int64) ([]MeetingEvent, error) {

	eventList := []MeetingEvent{}
	// ********* Use the client to do things... ****************
	srv, err := calendar.New(p.Client)
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	t := time.Now().Add(time.Hour * -1).Format(time.RFC3339)
	events, err := srv.Events.List("primary").ShowDeleted(false).
		SingleEvents(true).TimeMin(t).MaxResults(num).OrderBy("startTime").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
	}
	// fmt.Println("Upcoming events:")
	if len(events.Items) == 0 {
		fmt.Println("No upcoming events found.")
	} else {
		for _, item := range events.Items {
			// fmt.Printf("%s\n", item.Summary)
			date := item.Start.DateTime
			if date == "" {
				date = item.Start.Date
			}
			url, ok := MeetingURLFromEvent(item)
			if !ok {
				// do we care?
			} else {
				soon := IsMeetingSoon(item)
				// 	_ = open.Run(url.String())
				// 	fmt.Printf("%v,(%v),<%s>\n", item.Summary, date, url)
				startTime, err := MeetingStartTime(item)
				if err != nil {
					// no action needed
				} else {
					eventList = append(eventList, MeetingEvent{
						Description:     item.Summary,
						Start:           startTime,
						MeetingProvider: "No Need",
						MeetingLink:     *url,
						IsMeetingSoon:   soon,
					})
				}
			}

		}
	}
	return eventList, nil
}

// IsMeetingSoon returns true if the meeting is less than 5 minutes from now.
func IsMeetingSoon(event *calendar.Event) bool {
	startTime, err := MeetingStartTime(event)
	if err != nil {
		return false
	}
	minutesUntilStart := time.Until(startTime).Minutes()
	return -5 < minutesUntilStart && minutesUntilStart < 5
}

// MeetingStartTime returns the calendar event's start time.
func MeetingStartTime(event *calendar.Event) (time.Time, error) {
	if event == nil || event.Start == nil || event.Start.DateTime == "" {
		return time.Time{}, errors.New("event does not have a start datetime")
	}
	return time.Parse(googleCalendarDateTimeFormat, event.Start.DateTime)
}

// MeetingURLFromEvent returns a URL if the event is a Zoom meeting.
func MeetingURLFromEvent(event *calendar.Event) (*url.URL, bool) {
	input := event.Location + " " + event.Description
	if videoEntryPointURL, ok := conferenceVideoEntryPointURL(event); ok {
		input = videoEntryPointURL + " " + input
	}

	stringURL := ""
	haveMatch := false
	// ZOOM Matches
	matches := zoomURLRegexp.FindAllStringSubmatch(input, -1)
	if len(matches) == 0 || len(matches[0]) == 0 {
		// fmt.Println("No matches...")
		//return nil, false
	} else {

		haspass := zoomURLRegexpPwd.FindAllStringSubmatch(event.Description, -1)
		passcode := ""

		//fmt.Println("~~~~~~~~~~~~~~~~~~~")
		//fmt.Println(fmt.Sprintf("%#v", event.ConferenceData))
		if event.ConferenceData != nil {
			//fmt.Println(event.ConferenceData.ConferenceId)
			//fmt.Println(event.ConferenceData.Notes)
			if strings.Contains(event.ConferenceData.Notes, "Passcode") {
				passcode = strings.TrimSpace(strings.Split(event.ConferenceData.Notes, ":")[1])
				// fmt.Println(pc)
			}
		}
		//fmt.Println(event.ConferenceData)
		//fmt.Println("~~~~~~~~~~~~~~~~~~~")
		// By default, match the whole URL.
		stringURL = matches[0][0]

		// If we have a meeting ID in the URL, then use zoommtg:// instead of the HTTPS URL.
		if len(matches[0]) >= 2 {
			if _, err := strconv.Atoi(matches[0][1]); err == nil {
				if len(haspass) >= 1 {
					if len(haspass[0]) >= 2 {
						stringURL = "zoommtg://zoom.us/join?confno=" + matches[0][1] + "&pwd=" + haspass[0][1]
					} else {
						stringURL = "zoommtg://zoom.us/join?confno=" + matches[0][1]
					}
				} else {
					if passcode == "" {
						stringURL = "zoommtg://zoom.us/join?confno=" + matches[0][1]
					} else {
						stringURL = "zoommtg://zoom.us/join?confno=" + matches[0][1] + "&pwd=" + passcode
					}
				}
				haveMatch = true
			}
		}
	}
	if !haveMatch {
		// MS TEAMS Matches
		// fmt.Println("trying ms")
		matches = teamsURLRegexp.FindAllStringSubmatch(event.Description, -1)
		if len(matches) == 0 || len(matches[0]) == 0 {
			//fmt.Println("No matches...")
			//return nil, false
		} else {

			// haspass := zoomURLRegexpPwd.FindAllStringSubmatch(event.Description, -1)

			// By default, match the whole URL.
			stringURL = matches[0][0]
			stringURL = strings.TrimSuffix(stringURL, ">")

			haveMatch = true
			// If we have a meeting ID in the URL, then use zoommtg:// instead of the HTTPS URL.
			// if len(matches[0]) >= 2 {
			// 	if _, err := strconv.Atoi(matches[0][1]); err == nil {
			// 		//if len(haspass) >= 1 {
			// 		// if len(haspass[0]) >= 2 {
			// 		// 	stringURL = "zoommtg://zoom.us/join?confno=" + matches[0][1] + "&pwd=" + haspass[0][1]
			// 		// } else {
			// 		// 	stringURL = "zoommtg://zoom.us/join?confno=" + matches[0][1]
			// 		// }
			// 		//} else {
			// 		stringURL = "teams://teams.microsoft.com/join?confno=" + matches[0][1]
			// 		//}
			// 		haveMatch = true
			// 	}
			// }
		}
	}
	if !haveMatch {
		// Cisco WebEx Matches
		fmt.Println("trying webex")
		matches = webexURLRegexp.FindAllStringSubmatch(event.Description, -1)
		if len(matches) == 0 || len(matches[0]) == 0 {
			fmt.Println("No matches...")
			return nil, false
		} else {

			// haspass := zoomURLRegexpPwd.FindAllStringSubmatch(event.Description, -1)

			// By default, match the whole URL.
			stringURL = matches[0][0]
			stringURL = strings.TrimSuffix(stringURL, ">")

			haveMatch = true
			// If we have a meeting ID in the URL, then use zoommtg:// instead of the HTTPS URL.
			// if len(matches[0]) >= 2 {
			// 	if _, err := strconv.Atoi(matches[0][1]); err == nil {
			// 		//if len(haspass) >= 1 {
			// 		// if len(haspass[0]) >= 2 {
			// 		// 	stringURL = "zoommtg://zoom.us/join?confno=" + matches[0][1] + "&pwd=" + haspass[0][1]
			// 		// } else {
			// 		// 	stringURL = "zoommtg://zoom.us/join?confno=" + matches[0][1]
			// 		// }
			// 		//} else {
			// 		stringURL = "teams://teams.microsoft.com/join?confno=" + matches[0][1]
			// 		//}
			// 		haveMatch = true
			// 	}
			// }
		}
	}

	//fmt.Println(stringURL)
	parsedURL, err := url.Parse(stringURL)
	if err != nil {
		return nil, false
	}
	return parsedURL, haveMatch
}

// conferenceVideoEntryPointURL returns the URL for the video entrypoint if one exists.
func conferenceVideoEntryPointURL(event *calendar.Event) (string, bool) {
	if event.ConferenceData == nil {
		return "", false
	}

	for _, entryPoint := range event.ConferenceData.EntryPoints {
		fmt.Println(fmt.Sprintf("ep-%s, %s", entryPoint.EntryPointType, entryPoint.Uri))
		if entryPoint.EntryPointType == "video" && strings.Contains(entryPoint.Uri, "zoom") {
			return entryPoint.Uri, true
		}
		if entryPoint.EntryPointType == "video" && strings.Contains(entryPoint.Uri, "teams") {
			return entryPoint.Uri, true
		}
	}

	return "", false
}

// readGoogleClientConfigFromFile reads the content of a file and parses it as an *oauth2.Config
func readGoogleClientConfigFromFile(filepath string) (*oauth2.Config, error) {
	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	//fmt.Println(fmt.Sprintf("%#v", string(b[:])))
	// If modifying these scopes, delete your previously saved client_secret.json.
	conf, err := google.ConfigFromJSON(b, calendar.CalendarReadonlyScope)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return conf, nil
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	usr, err := user.Current()
	if err != nil {
		fmt.Println("No current user")
		os.Exit(1)
	}
	tokFile := filepath.Join(usr.HomeDir, ".config", "google", "no-more-lateness_token.json")
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
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
