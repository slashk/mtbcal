package actions

import (
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/gorilla/schema"
	"github.com/pkg/errors"
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

func findEventFromUUID(c buffalo.Context) (models.Event, error) {
	var e models.Event
	err := models.DB.Find(&e, c.Param("event_id"))
	return e, err
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
	// end alternate Bind
	return e, nil
}

// List default implementation.
func (v *EventsResource) List(c buffalo.Context) error {
	var popular, upcoming, updated models.Events
	err := models.DB.Scope(models.Popular()).All(&popular)
	if err != nil {
		// TODO handle error
		// c.Render(200, r.String("Shit problems"))
	}
	err = models.DB.Scope(models.Upcoming()).All(&upcoming)
	if err != nil {
		// TODO handle error
		// c.Render(200, r.String("Shit problems"))
	}
	err = models.DB.Scope(models.Updated()).All(&updated)
	if err != nil {
		// TODO handle error
		// c.Render(200, r.String("Shit problems"))
	}
	c.Set("popular", popular)
	c.Set("upcoming", upcoming)
	c.Set("updated", updated)
	return c.Render(200, r.HTML("events/index.html"))
}

// Show default implementation.
func (v *EventsResource) Show(c buffalo.Context) error {
	e, err := findEventFromUUID(c)
	if err != nil {
		// TODO handle error gracefully
		return c.Render(500, r.String("Event id not found"))
	}
	races, err := models.FindRacesFromEvent(e)
	if err != nil {
		// TODO handle error gracefully
		c.LogField("error", err.Error())
		return c.Render(500, r.String(err.Error()))
	}
	c.Set("e", e)
	c.Set("races", races)
	c.Set("page", pageDefault)
	return c.Render(200, r.HTML("events/show.html"))
}

// New default implementation.
func (v *EventsResource) New(c buffalo.Context) error {
	e := models.NewEmptyEvent()
	c.Set("e", e)
	c.Set("page", pageDefault)
	return c.Render(200, r.HTML("events/new.html"))
}

// Create default implementation.
func (v *EventsResource) Create(c buffalo.Context) error {
	e := models.Event{}
	// Alternate to bind due to time.Time parsing
	// the usual would be to do `err = c.Bind(&e)`
	err := c.Request().ParseForm()
	if err != nil {
		return errors.WithStack(err)
	}
	dec := schema.NewDecoder()
	dec.IgnoreUnknownKeys(true)
	dec.ZeroEmpty(true)
	// this is the money call that gets us a time parser
	dec.RegisterConverter(time.Time{}, ConvertFormDate)
	// this is the equivalent to Bind(&e)
	err = dec.Decode(&e, c.Request().PostForm)
	c.LogField("event", e)
	// end alternate Bind
	// e.Active = true
	verrs, err := e.Validate()
	if err != nil {
		return errors.WithStack(err)
	}
	if verrs.HasAny() {
		c.Set("e", e)
		c.Set("errors", verrs.Errors)
		return c.Render(422, r.HTML("events/new.html"))
	}
	// c.LogField("event", e)
	e.Active = true
	err = models.DB.Create(&e)
	if err != nil {
		c.Set("e", e)
		c.Set("errors", "Database save error")
		return c.Render(422, r.HTML("events/new.html"))
	}
	err = models.DB.Reload(&e)
	if err != nil {
		c.Set("e", e)
		c.Set("errors", "Database reload error")
		return c.Render(422, r.HTML("events/new.html"))
	}
	c.Set("e", e)
	c.Set("page", pageDefault)
	c.LogField("new event id", e.ID)
	c.Flash().Add("success", "Event created successfully")
	return c.Redirect(301, "/events/%s", e.ID.String())
}

// Edit returns a editing page with the event record filled in
func (v *EventsResource) Edit(c buffalo.Context) error {
	e, err := findEventFromUUID(c)
	if err != nil {
		// TODO handle error
		return c.Render(500, r.String("Event id not found"))
	}
	c.Set("e", e)
	c.Set("page", pageDefault)
	return c.Render(200, r.HTML("events/edit.html"))
}

// Update changes a existing event record
func (v *EventsResource) Update(c buffalo.Context) error {
	e, err := findEventFromUUID(c)
	if err != nil {
		c.Set("e", e)
		c.Set("errors", "Event not found in database")
		return c.Render(422, r.HTML("events/index.html"))
	}
	// Alternate to bind due to time.Time parsing
	// the usual would be to do `err = c.Bind(&e)`
	err = c.Request().ParseForm()
	if err != nil {
		// TODO handle error
		return errors.WithStack(err)
	}
	dec := schema.NewDecoder()
	dec.IgnoreUnknownKeys(true)
	dec.ZeroEmpty(true)
	// this is the money call that gets us a time parser
	dec.RegisterConverter(time.Time{}, ConvertFormDate)
	// this is the equivalent to Bind(&e)
	err = dec.Decode(&e, c.Request().PostForm)
	// end alternate Bind
	if c.Request().PostForm.Get("WebReg") == "" {
		e.WebReg = false
	}
	c.LogField("event", e)
	if err != nil {
		return errors.WithStack(err)
	}
	verrs, err := e.Validate()
	if err != nil {
		return errors.WithStack(err)
	}
	if verrs.HasAny() {
		c.Set("e", e)
		c.Set("errors", verrs.Errors)
		return c.Render(422, r.HTML("events/edit.html"))
	}
	err = models.DB.Update(&e)
	if err != nil {
		// TODO should this be a 500 error ?
		c.Set("e", e)
		c.Set("errors", "Cannot reload event from database")
		return c.Render(422, r.HTML("events/edit.html"))
	}
	err = models.DB.Reload(&e)
	if err != nil {
		// TODO should this be a 500 error ?
		c.Set("e", e)
		c.Set("errors", "Cannot reload event from database")
		return c.Render(500, r.HTML("events/edit.html"))
	}
	c.Set("e", e)
	c.Set("page", pageDefault)
	return c.Redirect(301, "/events/%s", e.ID)
}

// Destroy default implementation.
func (v *EventsResource) Destroy(c buffalo.Context) error {
	// TODO admin middleware
	e, err := findEventFromUUID(c)
	if err != nil {
		return c.Render(404, r.String("event cannot be found"))
	}
	e.Active = false
	// don't actually delete in the DB, just mark inactive
	// err = models.DB.Destroy(&models.Event{ID: e.ID})
	err = models.DB.Update(&e)
	if err != nil {
		// TODO add flash
		// return c.Render(422, r.HTML("events/edit.html"))
		return c.Render(422, r.String("event cannot be updated in DB"))
	}
	c.Set("page", pageDefault)
	return c.Redirect(301, "/events")
}
