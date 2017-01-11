package models_test

import (
	"testing"

	"github.com/slashk/mtbcal/models"
	"github.com/stretchr/testify/require"
)

func TestUserString(t *testing.T) {
	r := require.New(t)
	u := models.User{
		Login:   "Mark",
		Email:   "mark@example.com",
		Twitter: "markb",
	}
	r.Contains(u.String(), "mark@example.com")
}

func TestUsersString(t *testing.T) {
	r := require.New(t)
	u1 := models.User{
		Login:   "Mark",
		Email:   "mark@example.com",
		Twitter: "markb",
	}
	u2 := models.User{
		Login:   "Ken",
		Email:   "ken@example.com",
		Twitter: "kenp",
	}
	users := models.Users{u1, u2}
	r.Contains(users.String(), "mark@example.com")
	r.Contains(users.String(), "ken@example.com")
}

func TestUserValidateWValidUser(t *testing.T) {
	r := require.New(t)
	u := models.User{
		Login:   "Mark",
		Email:   "mark@example.com",
		Twitter: "markb",
	}
	verrs, err := u.Validate()
	r.NoError(err)
	if verrs.HasAny() {
		r.Fail("user did not validate")
	}
}

func TestUserValidateWNoLogin(t *testing.T) {
	r := require.New(t)
	u := models.User{
		Login:   "",
		Email:   "mark@example.com",
		Twitter: "markb",
	}
	verrs, err := u.Validate()
	r.NoError(err)
	// Users are require to have login
	if verrs.Count() != 1 {
		r.Fail("user validates with no login")
	}
}
