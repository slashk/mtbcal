package models

import (
	"time"

	"encoding/json"

	"github.com/markbates/pop"
	"github.com/markbates/validate"
	"github.com/markbates/validate/validators"
)

// User represents mtbcalendar site user
type User struct {
	ID            int       `json:"id" db:"id"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
	Login         string    `json:"login" db:"login"`
	Email         string    `json:"email" db:"email"`
	Hometown      string    `json:"hometown" db:"hometown"`
	Avatar        string    `json:"avatar" db:"avatar"`
	Active        bool      `json:"active" db:"active"`
	PublicProfile bool      `json:"public_profile" db:"public_profile"`
	Twitter       string    `json:"twitter" db:"twitter"`
	Facebook      string    `json:"facebook" db:"facebook"`
	Admin         bool      `json:"admin" db:"admin"`
}

// String is not required by pop and may be deleted
func (u User) String() string {
	b, _ := json.Marshal(u)
	return string(b)
}

// Users is not required by pop and may be deleted
type Users []User

// String is not required by pop and may be deleted
func (u Users) String() string {
	b, _ := json.Marshal(u)
	return string(b)
}

// Validate gets run everytime you call a "pop.Validate" method.
// This method is not required and may be deleted.
func (u *User) Validate() (*validate.Errors, error) {
	verrs := validate.Validate(
		&validators.StringIsPresent{Name: "Login", Field: u.Login},
		&validators.StringIsPresent{Name: "Email", Field: u.Email},
	)
	return verrs, nil
}

// ValidateSave gets run everytime you call "pop.ValidateSave" method.
// This method is not required and may be deleted.
func (u *User) ValidateSave() (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run everytime you call "pop.ValidateUpdate" method.
// This method is not required and may be deleted.
func (u *User) ValidateUpdate() (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ByLogin scopes a query to find by login name
func ByLogin(l string) pop.ScopeFunc {
	return func(q *pop.Query) *pop.Query {
		return q.Where("login = ?", l)
	}
}
