package calendar

//import errors to log errors when they occur
import (
	"errors"
	"net/url"
	"time"

	"golang.org/x/oauth2"
)

// MeetingEvent - Structure that we care about returning
type MeetingEvent struct {
	Description     string    `json:"description"`
	Start           time.Time `json:"start"`
	End             time.Time `json:"end"`
	MeetingProvider string    `json:"meetingprovider"`
	MeetingLink     url.URL   `json:"meetinglink"`
	IsMeetingSoon   bool      `json:"ismeetingsoon"`
	MeetingResponse string    `json:"meetingresponse"`
}

// Provider = The main interface used to describe appliances
type Provider interface {
	GetClient() (bool, error)
	GetEvents(num int64, attendee string) ([]MeetingEvent, error)
	GetAuthURL() string
	GetToken() (*oauth2.Token, error)
}

//Our appliance types
const (
	GOOGLE = "google"
)

// GetProvider - Function to create the appliances
func GetProvider(t string) (Provider, error) {
	//Use a switch case to switch between types, if a type exist then error is nil (null)
	switch t {
	case GOOGLE:
		return new(GoogleCal), nil
	default:
		//if type is invalid, return an error
		return nil, errors.New("Unsupported Provider")
	}
}
