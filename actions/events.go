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

// List default implementation.
func (v *EventsResource) List(c buffalo.Context) error {
	var popular, upcoming, updated models.Events
	err := models.DB.Scope(models.Popular()).All(&popular)
	if err != nil {
		// c.Render(200, r.String("Shit problems"))
	}
	err = models.DB.Scope(models.Upcoming()).All(&upcoming)
	if err != nil {
		// c.Render(200, r.String("Shit problems"))
	}
	err = models.DB.Scope(models.Updated()).All(&updated)
	if err != nil {
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
		return c.Render(500, r.String("Event id not found"))
	}
	c.Set("e", e)
	c.Set("page", pageDefault)
	return c.Render(200, r.HTML("events/show.html"))
}

// New default implementation.
func (v *EventsResource) New(c buffalo.Context) error {
	e := models.Event{}
	c.Set("e", e)
	c.Set("page", pageDefault)
	return c.Render(200, r.String("events/new.html"))
}

// Create default implementation.
func (v *EventsResource) Create(c buffalo.Context) error {
	e := models.Event{}
	err := c.Bind(&e)
	if err != nil {
		return c.Render(422, r.String("new event not validated"))
	}
	verrs, err := e.Validate()
	if err != nil {
		return errors.WithStack(err)
	}
	if verrs.HasAny() {
		c.Set("verrs", verrs.Errors)
		c.Set("user", e)
		return c.Render(422, r.HTML("users/new.html"))
	}
	err = models.DB.Create(&e)
	if err != nil {
		return c.Render(422, r.String("new event cannot be saved to DB"))
	}
	c.Set("e", e)
	c.Set("page", pageDefault)
	return c.Redirect(301, "/events/%d", e.ID)
}

// Edit returns a editing page with the event record filled in
func (v *EventsResource) Edit(c buffalo.Context) error {
	e, err := findEventFromUUID(c)
	if err != nil {
		return c.Render(500, r.String("Event id not found"))
	}
	c.Set("e", e)
	c.Set("page", pageDefault)
	return c.Render(200, r.HTML("events/edit.html"))
}

// Update chnages a existing event record
func (v *EventsResource) Update(c buffalo.Context) error {
	e, err := findEventFromUUID(c)
	if err != nil {
		c.LogField("error", err)
		return c.Render(500, r.String("Event id not found"))
	}
	// Alternate to bind due to time.Time parsing
	// the usual would be to do `err = c.Bind(&e)`
	err = c.Request().ParseForm()
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
	// end alternate Bind
	if err != nil {
		return errors.WithStack(err)
	}
	verrs, err := e.Validate()
	if err != nil {
		return errors.WithStack(err)
	}
	if verrs.HasAny() {
		c.Set("verrs", verrs.Errors)
		c.Set("user", e)
		return c.Render(422, r.HTML("events/edit.html"))
	}
	err = models.DB.Update(&e)
	if err != nil {
		// TODO add flash
		// return c.Render(422, r.HTML("events/edit.html"))
		return c.Render(422, r.String("event cannot be updated in DB"))
	}
	err = models.DB.Reload(&e)
	if err != nil {
		return c.Render(500, r.String("cannot reload event object"))
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
	err = models.DB.Destroy(&models.Event{ID: e.ID})
	if err != nil {
		return c.Render(500, r.String("event cannot be deleted from DB"))
	}
	c.Set("page", pageDefault)
	return c.Redirect(301, "/events")
}
