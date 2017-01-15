package actions_test

import (
	"os"
	"testing"

	"github.com/markbates/willie"
	"github.com/slashk/mtbcal/actions"
	"github.com/slashk/mtbcal/models"
	"github.com/stretchr/testify/require"
)

// TestMain is the setup/teardown
func TestMain(m *testing.M) {
	models.DB.MigrateUp("../migrations")
	ret := m.Run()
	models.DB.MigrateDown("../migrations")
	os.Exit(ret)
}

func dummyUser() *models.User {
	return &models.User{
		Login:   "Mark",
		Email:   "mark@example.com",
		Twitter: "markb",
		Active:  true,
	}
}

func Test_UsersList(t *testing.T) {
	r := require.New(t)
	w := willie.New(actions.App())
	u := dummyUser()
	r.NoError(models.DB.Create(u))
	res := w.Request("/users").Get()
	r.Equal(200, res.Code)
	r.Contains(res.Body.String(), u.Login)
}

func Test_UsersShow(t *testing.T) {
	r := require.New(t)
	w := willie.New(actions.App())
	u := dummyUser()
	r.NoError(models.DB.Create(u))
	res := w.Request("/users/%d", u.ID).Get()
	r.Equal(200, res.Code)
	r.Contains(res.Body.String(), u.Email)
}

func Test_UsersCreate(t *testing.T) {
	r := require.New(t)
	w := willie.New(actions.App())
	ct, err := models.DB.Count("users")
	r.NoError(err)
	// r.Equal(0, ct)
	oct := ct
	u := dummyUser()
	res := w.Request("/users").Post(u)
	r.Equal(301, res.Code)
	ct, err = models.DB.Count("users")
	r.NoError(err)
	r.Equal(oct+1, ct)
}

func Test_UsersCreate_HandlesErrors(t *testing.T) {
	r := require.New(t)
	w := willie.New(actions.App())
	// ct, err := models.DB.Count("users")
	// r.NoError(err)
	// r.Equal(0, ct)
	u := &models.User{}
	res := w.Request("/users").Post(u)
	r.Equal(422, res.Code)
	// TODO implement flash messages to test
	// r.Contains(res.Body.String(), "valid")
}

func Test_UsersUpdate(t *testing.T) {
	r := require.New(t)
	w := willie.New(actions.App())
	u := dummyUser()
	r.NoError(models.DB.Create(u))
	res := w.Request("/users/%d", u.ID).Put(map[string]string{
		"email": "bates@example.com",
	})
	r.Equal(301, res.Code)
	r.NoError(models.DB.Reload(u))
	r.Equal("bates@example.com", u.Email)
}

func Test_UsersUpdate_HandlesErrors(t *testing.T) {
	r := require.New(t)
	w := willie.New(actions.App())
	u := dummyUser()
	r.NoError(models.DB.Create(u))
	res := w.Request("/users/%d", u.ID).Put(map[string]string{
		// "email": "bates@example.com",
		"email": "",
	})
	r.Equal(422, res.Code)
	// TODO implement flash messages to test
	// r.Contains(res.Body.String(), "First Name can not be blank.")
}

func Test_UsersDestroy(t *testing.T) {
	r := require.New(t)
	w := willie.New(actions.App())
	u := dummyUser()
	r.NoError(models.DB.Create(u))
	ct, err := models.DB.Count("users")
	r.NoError(err)
	oct := ct
	res := w.Request("/users/%d", u.ID).Delete()
	r.Equal(301, res.Code)
	ct, err = models.DB.Count("users")
	r.NoError(err)
	r.Equal(oct-1, ct)
}
