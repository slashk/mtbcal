package models

import (
	"encoding/json"
	"time"

	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/satori/go.uuid"
)

// Race details a specific competition at an event
type Race struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`

	Cost        string    `json:"cost" db:"cost"`
	License     string    `json:"license" db:"license"`
	Description string    `json:"description" db:"description"`
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
func (r *Race) Validate(tx *pop.Connection) (*validate.Errors, error) {
	// ensure that format_id exists
	// ensure that event_id exists
	return validate.NewErrors(), nil
}

// ValidateSave gets run everytime you call "pop.ValidateSave" method.
func (r *Race) ValidateSave(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run everytime you call "pop.ValidateUpdate" method.
func (r *Race) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
