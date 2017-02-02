package models_test

import (
	"testing"
	"time"

	"github.com/slashk/mtbcal/models"
	"github.com/stretchr/testify/require"
)

func TestEventString(t *testing.T) {
	r := require.New(t)
	e := models.Event{
		Name:         "BigBoy",
		Location:     "Folsom, CA",
		State:        "CA",
		WebReg:       true,
		Active:       true,
		PublishedAt:  time.Now(),
		StartDate:    time.Now(),
		EndDate:      time.Now(),
		RegOpenDate:  time.Now(),
		RegCloseDate: time.Now(),
		Description:  "This is the bigboy event",
		URL:          "http://www.mtbcalendar.com",
		UserID:       "admin",
		Country:      "US",
		Lng:          -95.2353,
		Lat:          38.9717,
	}
	r.Contains(e.String(), "BigBoy")
}
