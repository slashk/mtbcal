package models

import (
	"context"
	"errors"
	"net/url"
	"os"
	"time"

	"encoding/json"

	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/markbates/validate/validators"
	"github.com/satori/go.uuid"
	"googlemaps.github.io/maps"
)

// Event represents a group of races
type Event struct {
	ID uuid.UUID `json:"id" db:"id"`

	LegacyID     int       `db:"legacy_id"`
	Name         string    `db:"name"                  json:"name"`
	Description  string    `db:"description,size:1024" json:"description,size:1024"`
	Location     string    `db:"location"              json:"location"`
	URL          string    `db:"url"                   json:"url"`
	StartDate    time.Time `db:"start_date"            json:"start_date"`
	EndDate      time.Time `db:"end_date"              json:"end_date"`
	WebReg       bool      `db:"web_reg"               json:"web_reg"`
	RegOpenDate  time.Time `db:"reg_open_date"         json:"reg_open_date"`
	RegCloseDate time.Time `db:"reg_close_date"        json:"reg_close_date"`
	UserID       string    `db:"user_id"               json:"user_id"`
	Region       string    `db:"region"                json:"region"`
	RegionID     int       `db:"region_id"             json:"region_id"`
	Lat          float32   `db:"lat"                   json:"lat"`
	Lng          float32   `db:"lng"                   json:"lng"`
	PublishedAt  time.Time `db:"published_at"          json:"published_at"`
	Permalink    string    `db:"permalink"             json:"permalink"`
	State        string    `db:"state"                 json:"state"`
	Active       bool      `db:"active"                json:"active"`
	Dupe         bool      `db:"dupe"                  json:"dupe"`
	Twin         string    `db:"twin"                  json:"twin"`
	Country      string    `db:"country"               json:"country"`

	// Races []Race

	CreatedAt time.Time `json:"created_at" db:"created_at"`

	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (e Event) String() string {
	b, _ := json.Marshal(e)
	return string(b)
}

// Events is not required by pop and may be deleted
type Events []Event

// String is not required by pop and may be deleted
func (e Events) String() string {
	b, _ := json.Marshal(e)
	return string(b)
}

// Validate gets run everytime you call a "pop.Validate" method.
func (e *Event) Validate() (*validate.Errors, error) {
	v := validate.NewErrors()
	err := e.Geocode()
	if err != nil {
		v.Add("geocode", err.Error())
	}
	v.Append(validate.Validate(&validators.TimeIsBeforeTime{
		FirstName:  "StartDate",
		FirstTime:  e.StartDate,
		SecondName: "EndDate",
		SecondTime: e.EndDate,
	}))
	if e.WebReg {
		v.Append(validate.Validate(&validators.TimeIsBeforeTime{
			FirstName:  "RegOpenDate",
			FirstTime:  e.RegOpenDate,
			SecondName: "RegCloseDate",
			SecondTime: e.RegCloseDate,
		}))
	}
	v.Append(validate.Validate(
		&validators.StringIsPresent{Field: e.Name, Name: "Name"},
		&validators.StringIsPresent{Field: e.Location, Name: "Location"},
	))
	return v, nil
}

// Geocode method find lat, lng, location and state for event location strings
func (e *Event) Geocode() error {
	var client *maps.Client
	var err error
	mapAPI := os.Getenv("GMAP_API")
	if mapAPI == "" {
		return errors.New("No GMAP_API variable set for geocoding")
	}
	r := &maps.GeocodingRequest{
		Address: e.Location,
	}
	client, err = maps.NewClient(maps.WithAPIKey(mapAPI))
	resp, err := client.Geocode(context.Background(), r)
	if err != nil {
		return err
	}
	e.Lat = float32(resp[0].Geometry.Location.Lat)
	e.Lng = float32(resp[0].Geometry.Location.Lng)
	e.State = resp[0].AddressComponents[2].ShortName
	e.Location = resp[0].FormattedAddress
	return nil
}

// ValidURL checks that the event URL is valid
func (e *Event) ValidURL() bool {
	_, err := url.ParseRequestURI(e.URL)
	if err != nil {
		return false
	}
	return true
}

// Upcoming finds future events
func Upcoming() pop.ScopeFunc {
	today := time.Now()
	return func(q *pop.Query) *pop.Query {
		return q.Where("start_date >= ?", today).Where("active = ?", true).Order("start_date asc")
	}
}

// Popular find events people have added to their list
func Popular() pop.ScopeFunc {
	// TODO
	return func(q *pop.Query) *pop.Query {
		return q.Where("active = ?", true).Order("start_date asc").Limit(10)
	}
}

// Updated finds recently updated events
func Updated() pop.ScopeFunc {
	return func(q *pop.Query) *pop.Query {
		return q.Order("updated_at").Order("updated_at asc").Limit(10)
	}
}

// Active finds events which are still active
func Active() pop.ScopeFunc {
	return func(q *pop.Query) *pop.Query {
		return q.Where("active = ?", true)
	}
}

// Unique finds non-dupe events
func Unique() pop.ScopeFunc {
	return func(q *pop.Query) *pop.Query {
		return q.Where("dupe = ?", false)
	}
}

// NewEmptyEvent creates a valid new Event
func NewEmptyEvent() Event {
	return Event{
		StartDate:    time.Now(),
		EndDate:      time.Now(),
		RegOpenDate:  time.Now(),
		RegCloseDate: time.Now(),
		WebReg:       false,
		Active:       true,
		Dupe:         false,
		Lat:          0.0,
		Lng:          0.0,
	}
}
