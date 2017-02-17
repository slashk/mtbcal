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

func TestValidDateRange(t *testing.T) {
	r := require.New(t)

	eValid := models.Event{
		StartDate: time.Now().Add(20 * 24 * time.Hour),
		EndDate:   time.Now().Add(24 * 24 * time.Hour),
	}
	r.Equal(eValid.ValidDates(), true, "EndDates after StartDate should be valid")

	eInvalid := models.Event{
		StartDate: time.Now().Add(20 * 24 * time.Hour),
		EndDate:   time.Now().Add(4 * 24 * time.Hour),
	}
	r.Equal(eInvalid.ValidDates(), false, "EndDate before StartDate should be invalid")

	n := time.Now()
	eEqual := models.Event{
		StartDate: n,
		EndDate:   n,
	}
	r.Equal(eEqual.ValidDates(), true, "Equal dates should be valid")
}

func TestValidWebReg(t *testing.T) {
	r := require.New(t)

	eValid := models.Event{
		EndDate:      time.Now().Add(26 * 24 * time.Hour),
		RegOpenDate:  time.Now().Add(20 * 24 * time.Hour),
		RegCloseDate: time.Now().Add(24 * 24 * time.Hour),
	}
	r.Equal(eValid.ValidWebReg(), true, "RegClose after RegOpen should be valid")

	eInvalid := models.Event{
		EndDate:      time.Now().Add(26 * 24 * time.Hour),
		RegOpenDate:  time.Now().Add(20 * 24 * time.Hour),
		RegCloseDate: time.Now().Add(4 * 24 * time.Hour),
	}
	r.Equal(eInvalid.ValidWebReg(), false, "RegClose before RegOpen should be invalid")

	n := time.Now()
	eEqual := models.Event{
		EndDate:      time.Now().Add(26 * 24 * time.Hour),
		RegOpenDate:  n,
		RegCloseDate: n,
	}
	r.Equal(eEqual.ValidWebReg(), true, "Equal Reg should be valid")
}

func TestValidURL(t *testing.T) {
	r := require.New(t)
	valid := models.Event{URL: "http://google.com"}
	r.True(valid.ValidURL(), "should be valid url")
	invalid := models.Event{URL: "sdfsdfsf"}
	r.False(invalid.ValidURL(), "should be invalid url")
	a := models.Event{URL: "http//google.com"}
	r.False(a.ValidURL(), "should be invalid url")
}
