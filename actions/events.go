package actions

import (
	"github.com/gobuffalo/buffalo"
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
	var u models.Event
	err := models.DB.Find(&u, c.Param("event_id"))
	return u, err
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
	return c.Render(200, r.String("Events#New"))
}

// Create default implementation.
func (v *EventsResource) Create(c buffalo.Context) error {
	e := models.Event{}
	c.Set("e", e)
	c.Set("page", pageDefault)
	return c.Render(200, r.String("Events#Create"))
}

// Edit default implementation.
func (v *EventsResource) Edit(c buffalo.Context) error {
	e, err := findEventFromUUID(c)
	if err != nil {
		return c.Render(500, r.String("Event id not found"))
	}
	c.Set("e", e)
	c.Set("page", pageDefault)
	return c.Render(200, r.HTML("events/edit.html"))
}

// Update default implementation.
func (v *EventsResource) Update(c buffalo.Context) error {
	return c.Render(200, r.String("Events#Update"))
}

// Destroy default implementation.
func (v *EventsResource) Destroy(c buffalo.Context) error {
	return c.Render(200, r.String("Events#Destroy"))
}
