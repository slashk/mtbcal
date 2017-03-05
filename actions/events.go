package actions

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/gorilla/schema"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"github.com/slashk/mtbcal/models"
)

// EventsResource covers events
type EventsResource struct {
	buffalo.Resource
}

func init() {
	var resource buffalo.Resource
	resource = &EventsResource{&buffalo.BaseResource{}}
	App().Resource("/events", resource)
}

// List shows the front page with popular, upcoming and new events
func (v *EventsResource) List(c buffalo.Context) error {
	var popular, upcoming, updated models.Events
	err := models.DB.Scope(models.Popular()).All(&popular)
	if err != nil {
		c.LogField("popular db error", err)
	}
	err = models.DB.Scope(models.Upcoming()).All(&upcoming)
	if err != nil {
		c.LogField("upcoming db error", err)
	}
	err = models.DB.Scope(models.Updated()).All(&updated)
	if err != nil {
		c.LogField("updated db error", err)
	}
	c.Set("popular", popular)
	c.Set("upcoming", upcoming)
	c.Set("updated", updated)
	c.Set("page", pageDefault)
	return c.Render(200, r.HTML("events/index.html"))
}

// Show displays a single event
func (v *EventsResource) Show(c buffalo.Context) error {
	e, err := findEventFromUUID(c)
	if err != nil {
		c.LogField("error", err)
		c.Flash().Add("danger", "Event could not be found")
		return c.Redirect(301, "/events")
	}
	races, err := models.FindRacesFromEvent(e)
	if err != nil {
		// TODO handle error gracefully
		c.LogField("error", err.Error())
		return c.Render(500, r.String(err.Error()))
	}
	setEventAndPage(c, &e, &pageDefault, &races)
	// c.Set("races", races)
	return c.Render(200, r.HTML("events/show.html"))
}

// New tees up a blank form page to enter a new event
func (v *EventsResource) New(c buffalo.Context) error {
	e := models.NewEmptyEvent()
	rs := models.NewEmptyRace()
	setEventAndPage(c, &e, &pageDefault, &rs)
	f := raceFormatSelect()
	c.Set("f", f)
	return c.Render(200, r.HTML("events/new.html"))
}

// Create stores a new event in the database
func (v *EventsResource) Create(c buffalo.Context) error {
	c.LogField("postform", c.Request().PostForm)
	e, err := customEventDecode(c)
	verrs, err := e.Validate()
	if err != nil {
		// TODO gracefully flash and return to new or index ?
		return errors.WithStack(err)
	}
	if verrs.HasAny() {
		c.Flash().Add("danger", verrs.String())
		return c.Render(422, r.HTML("events/new.html"))
	}
	// setEventAndPage(c, &e, &pageDefault)
	e.Active = true
	err = models.DB.Create(&e)
	if err != nil {
		c.Flash().Add("danger", err.Error())
		return c.Render(422, r.HTML("events/new.html"))
	}
	err = models.DB.Reload(&e)
	if err != nil {
		c.Flash().Add("danger", err.Error())
		return c.Render(422, r.HTML("events/new.html"))
	}
	races, err := decodeRacesFromPost(c)
	for race := range races {
		races[race].EventID = e.ID
		// TODO handle race decode errors
		// verrs, err := races[race].Validate()
		// if err != nil {
		// 	return errors.WithStack(err)
		// }
		// if verrs.HasAny() {
		// 	c.Flash().Add("danger", verrs.String())
		// 	return c.Render(422, r.HTML("events/new.html"))
		// }
		err = models.DB.Create(&races[race])
		if err != nil {
			c.Flash().Add("danger", err.Error())
			return c.Render(422, r.HTML("events/new.html"))
		}
	}
	c.Flash().Add("success", "Event created successfully")
	return c.Redirect(301, "/events/%s", e.ID.String())
}

// Edit returns a editing page with the event record filled in
func (v *EventsResource) Edit(c buffalo.Context) error {
	e, err := findEventFromUUID(c)
	if err != nil {
		c.Flash().Add("danger", err.Error())
		return c.Render(404, r.HTML("events/index.html"))
	}
	races, err := models.FindRacesFromEvent(e)
	if err != nil {
		// TODO handle error gracefully
		c.LogField("error", err.Error())
		return c.Render(500, r.String(err.Error()))
	}
	setEventAndPage(c, &e, &pageDefault, &races)
	c.Set("f", models.FormatbyID)
	return c.Render(200, r.HTML("events/edit.html"))
}

// Update changes a existing event record
func (v *EventsResource) Update(c buffalo.Context) error {
	e, err := findEventFromUUID(c)
	if err != nil {
		c.Flash().Add("danger", "Find failed")
		c.Flash().Add("danger", err.Error())
		return c.Render(422, r.HTML("events/index.html"))
	}
	e, err = customEventDecode(c)
	if c.Request().PostForm.Get("WebReg") == "" {
		e.WebReg = false
	}
	if err != nil {
		return errors.WithStack(err)
	}
	c.LogField("event", e)
	races, err := decodeRacesFromPost(c)
	// races, err := models.FindRacesFromEvent(e)
	if err != nil {
		// TODO handle error gracefully
		c.LogField("error", err.Error())
		return c.Render(500, r.String(err.Error()))
	}
	for race := range races {
		// races[race].EventID = e.ID
		// TODO fix validations
		// races[race].Validate()
		// if verrs.HasAny() {
		// 	c.Flash().Add("danger", verrs.String())
		// 	c.LogField("validation error", verrs)
		// 	return c.Render(422, r.HTML("events/edit.html"))
		// }
		// verrs, err := e.Validate()
		// if err != nil {
		// 	return errors.WithStack(err)
		// }
		c.LogField(fmt.Sprintf("race-%d", race), races[race].Description)
		races[race].Validate()
		err = models.DB.Update(&races[race])
		if err != nil {
			c.Flash().Add("danger", err.Error())
			return c.Render(422, r.HTML("events/new.html"))
		}
		// err = models.DB.Reload(&races[race])
		// if err != nil {
		// 	c.Flash().Add("danger", "Reload failed")
		// 	c.Flash().Add("danger", err.Error())
		// 	return c.Render(500, r.HTML("events/edit.html"))
		// }
	}
	e.Validate()
	setEventAndPage(c, &e, &pageDefault, &races)
	err = models.DB.Update(&e)
	if err != nil {
		c.Flash().Add("danger", "Update failed")
		c.Flash().Add("danger", err.Error())
		return c.Render(422, r.HTML("events/edit.html"))
	}
	err = models.DB.Reload(&e)
	if err != nil {
		c.Flash().Add("danger", "Reload failed")
		c.Flash().Add("danger", err.Error())
		return c.Render(500, r.HTML("events/edit.html"))
	}
	c.Flash().Add("success", "Event updated successfully")
	return c.Redirect(301, "/events/%s", e.ID)
}

// Destroy default implementation.
func (v *EventsResource) Destroy(c buffalo.Context) error {
	// TODO admin middleware
	e, err := findEventFromUUID(c)
	if err != nil {
		c.Flash().Add("danger", err.Error())
		return c.Render(404, r.HTML("events/index.html"))
	}
	e.Active = false
	// TODO don't actually delete in the DB, just mark inactive
	// err = models.DB.Destroy(&models.Event{ID: e.ID})
	err = models.DB.Update(&e)
	if err != nil {
		c.Flash().Add("danger", err.Error())
		return c.Render(500, r.HTML("events/index.html"))
	}
	c.Set("page", pageDefault)
	c.Flash().Add("success", "Event deleted")
	return c.Redirect(301, "/events")
}

func findEventFromUUID(c buffalo.Context) (models.Event, error) {
	var e models.Event
	err := models.DB.Find(&e, c.Param("event_id"))
	return e, err
}

func setEventAndPage(c buffalo.Context, e *models.Event, p *PageDefaults, r *models.Races) {
	c.Set("e", e)
	c.Set("page", p)
	if len((*r)) > 0 {
		c.Set("races", r)
	}
}

func customEventDecode(c buffalo.Context) (models.Event, error) {
	e := models.Event{}
	// Alternate to bind due to time.Time parsing
	// the usual would be to do `err = c.Bind(&e)`
	err := c.Request().ParseForm()
	if err != nil {
		return e, err
	}
	dec := schema.NewDecoder()
	dec.IgnoreUnknownKeys(true)
	dec.ZeroEmpty(true)
	// this is the money call that gets us a time parser
	dec.RegisterConverter(time.Time{}, ConvertFormDate)
	// this is the equivalent to Bind(&e)
	err = dec.Decode(&e, c.Request().PostForm)
	if err != nil {
		return e, err
	}
	// this makes sure WebReg bool gets set if not present
	if c.Request().PostForm.Get("WebReg") == "" {
		e.WebReg = false
	}
	c.LogField("event", e)
	return e, nil
}

func decodeRacesFromPost(c buffalo.Context) (models.Races, error) {
	var races models.Races
	for x := 0; true; x++ {
		r := models.Race{}
		// comes across as 'Race.0.Cost'
		// TODO refactor this to some kind of custom bind
		k := map[string]string{
			"ID":          fmt.Sprintf("Race.%d.ID", x),
			"EventID":     fmt.Sprintf("Race.%d.EventID", x),
			"Cost":        fmt.Sprintf("Race.%d.Cost", x),
			"Description": fmt.Sprintf("Race.%d.Description", x),
			"URL":         fmt.Sprintf("Race.%d.URL", x),
			"License":     fmt.Sprintf("Race.%d.License", x),
			"FormatID": fmt.Sprintf("Race.%d.FormatID", x),
		}
		r.ID, _ = uuid.FromString(c.Request().PostFormValue(k["ID"]))
		r.EventID, _ = uuid.FromString(c.Request().PostFormValue(k["EventID"]))
		r.FormatID, _ = strconv.Atoi(c.Request().PostFormValue(k["FormatID"]))
		r.Cost = c.Request().PostFormValue(k["Cost"])
		r.Description = c.Request().PostFormValue(k["Description"])
		r.URL = c.Request().PostFormValue(k["URL"])
		r.License = c.Request().PostFormValue(k["License"])
		// TODO this is bad -- check all variables ?
		if r.License == "" {
			break
		}
		// c.LogField(fmt.Sprintf("bare-%d", x), r.Description)
		races = append(races, r)
	}
	c.LogField("races", races)
	return races, nil
}

func raceFormatSelect() []models.RaceFormat {
	a := []models.RaceFormat{}
	for k := range models.FormatbyID {
		a = append(a, models.FormatbyID[k])
	}
	return a
}
