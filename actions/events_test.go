package actions_test

import (
	"testing"
	"time"

	"github.com/markbates/willie"
	"github.com/slashk/mtbcal/actions"
	"github.com/slashk/mtbcal/models"
	"github.com/stretchr/testify/require"
)

func dummyEvent() *models.Event {
	return &models.Event{
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
		Dupe:         false,
	}
}

func Test_EventsResource_List(t *testing.T) {
	r := require.New(t)
	w := willie.New(actions.App())
	e := dummyEvent()
	r.NoError(models.DB.Create(e))
	res := w.Request("/events").Get()
	r.Equal(200, res.Code)
	r.Contains(res.Body.String(), e.Name)
}

func Test_EventsResource_Show(t *testing.T) {
	r := require.New(t)
	w := willie.New(actions.App())
	e := dummyEvent()
	r.NoError(models.DB.Create(e))
	models.DB.Reload(e)
	r.NotEmpty(e.ID, "event id is empty")
	res := w.Request("/events/" + e.ID.String()).Get()
	r.Equal(200, res.Code, res.Body.String())
	r.Contains(res.Body.String(), e.Name)
}

func Test_EventsResource_New(t *testing.T) {
	r := require.New(t)
	w := willie.New(actions.App())
	res := w.Request("/events/new").Get()
	r.Equal(200, res.Code)
	r.Contains(res.Body.String(), "Name")
}

// func Test_EventsResource_Create(t *testing.T) {
// 	r := require.New(t)
// 	r.Fail("Not Implemented!")
// }
//
func Test_EventsResource_Edit(t *testing.T) {
	r := require.New(t)
	w := willie.New(actions.App())
	e := dummyEvent()
	r.NoError(models.DB.Create(e))
	res := w.Request("/events/" + e.ID.String() + "/edit").Get()
	r.Equal(200, res.Code)
	r.Contains(res.Body.String(), "Name")
}

// func Test_EventsResource_Update(t *testing.T) {
// 	r := require.New(t)
// 	r.Fail("Not Implemented!")
// }

func Test_EventsResource_Destroy(t *testing.T) {
	var x models.Event
	r := require.New(t)
	w := willie.New(actions.App())
	e := dummyEvent()
	r.NoError(models.DB.Create(e))
	res := w.Request("/events/" + e.ID.String()).Delete()
	r.Equal(301, res.Code)
	err := models.DB.Find(&x, e.ID)
	if err != nil {
		r.Fail("event appears to have been deleted")
	}
	if x.Active == true {
		r.Fail("event was not marked inactive on destroy")
	}
}
