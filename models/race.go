package models

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/markbates/validate"
	"github.com/markbates/validate/validators"
	"github.com/satori/go.uuid"
)

// Race details a specific competition at an event
type Race struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	Cost        string    `json:"cost" db:"cost"`
	License     string    `json:"license" db:"license"`
	Description string    `db:"description,size:1024" json:"description,size:1024"`
	URL         string    `json:"url" db:"url"`
	EventID     uuid.UUID `json:"event_id" db:"event_id"`
	FormatID    int       `json:"format_id" db:"format_id"`
}

// String is not required by pop and may be deleted
func (r Race) String() string {
	b, _ := json.Marshal(r)
	return string(b)
}

// Races is not required by pop and may be deleted
type Races []Race

// String is not required by pop and may be deleted
func (r Races) String() string {
	b, _ := json.Marshal(r)
	return string(b)
}

// Validate gets run everytime you call a "pop.Validate" method.
func (r *Race) Validate() (*validate.Errors, error) {
	var v *validate.Errors
	// ensure that format_id exists
	// if raceFormatExists(r.FormatID) == false {
	// 	v.Add("raceformat", "race format does not exist")
	// }
	// ensure that event_id exists
	v.Append(validate.Validate(
		&validators.StringIsPresent{Field: r.License, Name: "License"},
	))
	return v, nil
}

// NewEmptyRace creates a valid new Event
func NewEmptyRace() Races {
	return Races{
		Race{
			Cost:        "",
			License:     "",
			Description: "",
			URL:         "",
			FormatID:    1,
		},
	}
}

// FindRacesFromEvent returns array of races attached to an event
func FindRacesFromEvent(e Event) (Races, error) {
	var r Races
	err := DB.BelongsTo(&e).All(&r)
	return r, err
}

// SaveRacesFromEventUUID returns array of races attached to an event
func SaveRacesFromEventUUID(r Races) error {
	var err error
	// TODO maybe this should be single race saves ?
	for x := range r {
		// validate these first
		sErr := DB.Save(r[x])
		if sErr != nil {
			// c.Render(200, r.String("Shit problems"))
		}
	}
	return err
}

func raceFormatExists(f int) bool {
	if _, ok := FormatbyID[f]; ok {
		return true
	}
	return false
}

// FormatString returns a comma delimited string of race formats
func (r Races) FormatString() string {
	if len(r) == 0 {
		return ""
	}
	var a []string
	for x := range r {
		a = append(a, FormatbyID[r[x].FormatID].US)
	}
	return strings.Join(a, ",")
}

// RaceFormat defines races formats
type RaceFormat struct {
	ID          int
	UCI         string
	US          string
	Name        string
	Description string
}

// FormatbyID is a map of int to Format
var FormatbyID = map[int]RaceFormat{
	1:  RaceFormat{ID: 1, UCI: "XC", US: "MXC", Name: "Cross Country", Description: ""},
	2:  RaceFormat{ID: 2, UCI: "DHI", US: "MDH", Name: "Downhill", Description: ""},
	3:  RaceFormat{ID: 3, Name: "Super D", UCI: "SD", US: "SD", Description: ""},
	4:  RaceFormat{ID: 4, Name: "MountainCross / 4X", UCI: "MTX", US: "4X", Description: ""},
	5:  RaceFormat{ID: 5, Name: "Dual Slalom", UCI: "DSL", US: "DS", Description: ""},
	6:  RaceFormat{ID: 6, Name: "Cyclocross", UCI: "CX", US: "CX", Description: ""},
	7:  RaceFormat{ID: 7, Name: "Marathon", UCI: "MAR", US: "MA", Description: ""},
	8:  RaceFormat{ID: 8, Name: "Short Track", UCI: "STX", US: "ST", Description: ""},
	9:  RaceFormat{ID: 9, Name: "Eight Hour", UCI: "8HR", US: "8H", Description: ""},
	10: RaceFormat{ID: 10, Name: "Twentyfour Hour", UCI: "24HR", US: "24H", Description: ""},
	11: RaceFormat{ID: 11, Name: "Time Trial", UCI: "MTT", US: "TT", Description: ""},
	12: RaceFormat{ID: 12, Name: "MTB Stage", UCI: "MSTG", US: "MST", Description: ""},
	13: RaceFormat{ID: 13, Name: "Singlespeed", UCI: "MSS", US: "SS", Description: ""},
	14: RaceFormat{ID: 14, Name: "Six Hour", UCI: "6HR", US: "6H", Description: ""},
	15: RaceFormat{ID: 15, Name: "Twelve Hour", UCI: "12HR", US: "12H", Description: ""},
	16: RaceFormat{ID: 16, Name: "Enduro", UCI: "END", US: "EN", Description: ""},
	17: RaceFormat{ID: 17, Name: "Slopestyle", UCI: "SLP", US: "SS", Description: ""},
	18: RaceFormat{ID: 18, Name: "Duathlon", UCI: "DUT", US: "DL", Description: ""},
	20: RaceFormat{ID: 20, Name: "Giant Slalom", UCI: "GIS", US: "GS", Description: ""},
	21: RaceFormat{ID: 21, Name: "Trials", UCI: "TRI", US: "TR", Description: ""},
	22: RaceFormat{ID: 22, Name: "Hill Climb", UCI: "HILL", US: "HC", Description: ""},
	23: RaceFormat{ID: 23, Name: "Other", UCI: "OTH", US: "OT", Description: ""},
	24: RaceFormat{ID: 24, Name: "MTB Triathlon", UCI: "MTRI", US: "MT", Description: ""},
}
