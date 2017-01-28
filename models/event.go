package models

import (
	"context"
	"errors"
	"os"
	"time"

	"encoding/json"

	"github.com/markbates/pop"
	"github.com/markbates/validate"
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
	RegOpenDate  string    `db:"reg_open_date"         json:"reg_open_date"`
	RegCloseDate string    `db:"reg_close_date"        json:"reg_close_date"`
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
// This method is not required and may be deleted.
func (e *Event) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateSave gets run everytime you call "pop.ValidateSave" method.
// This method is not required and may be deleted.
func (e *Event) ValidateSave(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run everytime you call "pop.ValidateUpdate" method.
// This method is not required and may be deleted.
func (e *Event) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// geocode uses google maps to find the lat, lng of an event location
func geocode(loc string) (float32, float32, string, string, error) {
	var client *maps.Client
	var err error
	mapAPI := os.Getenv("GMAP_API")
	if mapAPI != "" {
		client, err = maps.NewClient(maps.WithAPIKey(mapAPI))
	} else {
		return 0, 0, "", "", errors.New("No geocoding API key")
	}
	r := &maps.GeocodingRequest{
		Address: loc,
	}
	resp, err := client.Geocode(context.Background(), r)
	if err != nil {
		// log.WithError(err).Debug("geocoding error")
		return 0.0, 0.0, "", "", err
	}
	lat := float32(resp[0].Geometry.Location.Lat)
	lng := float32(resp[0].Geometry.Location.Lng)
	state := resp[0].AddressComponents[2].ShortName
	location := resp[0].FormattedAddress
	// log.WithFields(log.Fields{
	// 	"response": resp,
	// 	"state":    state,
	// 	"lat":      lat,
	// 	"location": location,
	// 	"lng":      lng,
	// }).Debug("geocode result")
	return lat, lng, state, location, nil
}

// Upcoming finds future events
func Upcoming() pop.ScopeFunc {
	today := time.Now()
	return func(q *pop.Query) *pop.Query {
		return q.Where("start_date > ?", today)
	}
}

// Popular find events people have added to their list
func Popular() pop.ScopeFunc {
	// TODO
	return func(q *pop.Query) *pop.Query {
		return q.Order("created_at").Limit(10)
	}
}

// Updated finds recently updated events
func Updated() pop.ScopeFunc {
	return func(q *pop.Query) *pop.Query {
		return q.Order("updated_at")
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
