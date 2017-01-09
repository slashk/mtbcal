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
	var u models.User
	k, err := c.ParamInt("user_id")
	if err != nil {
		// if no ID (int) found, lookup by login name (unique)
		k := c.Param("user_id")
		err = models.DB.Scope(models.ByLogin(k)).First(&u)
	} else {
		err = models.DB.Find(&u, k)
	}
	if err != nil {
		return c.Render(500, r.String("User id not found"))
	}
	c.Set("user", u)
	return c.Render(200, r.HTML("users/show.html"))
}

// New creates an empty user for creation
func (v *UsersResource) New(c buffalo.Context) error {
	return c.Render(200, r.String("Users#New"))
}

// Create default implementation.
func (v *UsersResource) Create(c buffalo.Context) error {
	var u models.User
	err := c.Bind(&u)
	if err != nil {
		return c.Render(500, r.String("bad user data rejected"))
	}
	err = models.DB.Save(&u)
	if err != nil {
		return c.Render(500, r.String("user cannot be saved to DB"))
	}
	return c.Render(200, r.String(u.String()))
}

// Edit default implementation.
func (v *UsersResource) Edit(c buffalo.Context) error {
	// TODO admin middleware
	var u models.User
	k, err := c.ParamInt("user_id")
	if err != nil {
		// if no ID (int) found, lookup by login name (unique)
		k := c.Param("user_id")
		err = models.DB.Scope(models.ByLogin(k)).First(&u)
	} else {
		err = models.DB.Find(&u, k)
	}
	if err != nil {
		return c.Render(500, r.String("User id not found"))
	}
	// TODO show form passing u into it
	return c.Render(200, r.String("Users#Edit"))
}

// Update default implementation.
func (v *UsersResource) Update(c buffalo.Context) error {
	return c.Render(200, r.String("Users#Update"))
}

// Destroy user
func (v *UsersResource) Destroy(c buffalo.Context) error {
	// TODO admin middleware
	k, err := c.ParamInt("user_id")
	if err != nil {
		return c.Render(404, r.String("user cannot be found"))
	}
	err = models.DB.Destroy(&models.User{ID: k})
	if err != nil {
		return c.Render(500, r.String("user cannot be deleted from DB"))
	}
	return c.Render(200, r.String("Users#Destroy"))
}
