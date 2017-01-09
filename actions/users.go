package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/slashk/mtbcal/models"
)

// UsersResource is the resource created for the User controller
type UsersResource struct {
	buffalo.Resource
}

func init() {
	var resource buffalo.Resource
	resource = &UsersResource{&buffalo.BaseResource{}}
	App().Resource("/users", resource)
}

func findUserFromParam(c buffalo.Context) (models.User, error) {
	var u models.User
	k, err := c.ParamInt("user_id")
	if err != nil {
		// if no ID (int) found, lookup by login name (unique)
		k := c.Param("user_id")
		err = models.DB.Scope(models.ByLogin(k)).First(&u)
	} else {
		err = models.DB.Find(&u, k)
	}
	return u, err
}

// List finds all registered users
func (v *UsersResource) List(c buffalo.Context) error {
	// limit this to admin
	var u models.Users
	err := models.DB.All(&u)
	if err != nil {
		return c.Render(500, r.String("User query error"))
	}
	c.Set("users", u)
	// TODO render a real user list page
	return c.Render(200, r.HTML("users/index.html"))
}

// Show default implementation.
func (v *UsersResource) Show(c buffalo.Context) error {
	u, err := findUserFromParam(c)
	if err != nil {
		return c.Render(500, r.String("User id not found"))
	}
	c.Set("user", u)
	return c.Render(200, r.HTML("users/show.html"))
}

// New creates an empty user for creation
func (v *UsersResource) New(c buffalo.Context) error {
	c.Set("user", models.User{})
	return c.Render(200, r.HTML("users/new.html"))
}

// Create default implementation.
func (v *UsersResource) Create(c buffalo.Context) error {
	var u models.User
	err := c.Bind(&u)
	if err != nil {
		return c.Render(500, r.String("bad user data rejected"))
	}
	err = models.DB.Create(&u)
	if err != nil {
		return c.Render(500, r.String("user cannot be saved to DB"))
	}
	return c.Redirect(301, "/users/%d", &u.ID)
}

// Edit default implementation.
func (v *UsersResource) Edit(c buffalo.Context) error {
	// TODO admin middleware
	u, err := findUserFromParam(c)
	if err != nil {
		return c.Render(500, r.String("User id not found"))
	}
	c.Set("user", u)
	// TODO show form passing u into it
	return c.Render(200, r.HTML("users/edit.html"))
	// return c.Render(200, r.String("Users#Edit"))
}

// Update default implementation.
func (v *UsersResource) Update(c buffalo.Context) error {
	u, err := findUserFromParam(c)
	if err != nil {
		return c.Render(404, r.String("user cannot be found"))
	}
	err = c.Bind(&u)
	if err != nil {
		return c.Render(500, r.String("bad user data rejected"))
	}
	err = models.DB.Update(&u)
	if err != nil {
		// return errors.WithStack(err)
	}
	err = models.DB.Reload(&u)
	if err != nil {
		// return errors.WithStack(err)
	}
	// return c.Render(200, r.String("Users#Update"))
	return c.Redirect(301, "/users/%d", u.ID)
}

// Destroy user
func (v *UsersResource) Destroy(c buffalo.Context) error {
	// TODO admin middleware
	u, err := findUserFromParam(c)
	if err != nil {
		return c.Render(404, r.String("user cannot be found"))
	}
	err = models.DB.Destroy(&models.User{ID: u.ID})
	if err != nil {
		return c.Render(500, r.String("user cannot be deleted from DB"))
	}
	return c.Render(200, r.String("Users#Destroy"))
}
