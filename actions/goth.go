package actions

import (
	"fmt"
	"os"

	"github.com/gobuffalo/buffalo"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/instagram"
	"github.com/markbates/goth/providers/twitter"
	"github.com/slashk/mtbcal/models"
)

func init() {
	gothic.Store = App().SessionStore

	goth.UseProviders(
		twitter.New(os.Getenv("TWITTER_KEY"), os.Getenv("TWITTER_SECRET"), fmt.Sprintf("%s%s", App().Host, "/auth/twitter/callback")),
		facebook.New(os.Getenv("FACEBOOK_KEY"), os.Getenv("FACEBOOK_SECRET"), fmt.Sprintf("%s%s", App().Host, "/auth/facebook/callback")),
		github.New(os.Getenv("GITHUB_KEY"), os.Getenv("GITHUB_SECRET"), fmt.Sprintf("%s%s", App().Host, "/auth/github/callback")),
		// fitbit.New(os.Getenv("FITBIT_KEY"), os.Getenv("FITBIT_SECRET"), fmt.Sprintf("%s%s", App().Host, "/auth/fitbit/callback")),
		instagram.New(os.Getenv("INSTAGRAM_KEY"), os.Getenv("INSTAGRAM_SECRET"), fmt.Sprintf("%s%s", App().Host, "/auth/instagram/callback")),
	)

	app := App().Group("/auth")
	app.GET("/{provider}", buffalo.WrapHandlerFunc(gothic.BeginAuthHandler))
	app.GET("/{provider}/callback", AuthCallback)
}

// AuthCallback provides provider callback handler for identity providers
func AuthCallback(c buffalo.Context) error {
	user, err := gothic.CompleteUserAuth(c.Response(), c.Request())
	if err != nil {
		return c.Error(401, err)
	}

	// TODO find user and create session
	// TODO if not found, register the user and create session

	c.LogField("User response", user)

	// register
	u := models.User{
		Login:    user.Name,
		Hometown: user.Location,
		Avatar:   user.AvatarURL,
		Email:    user.Email,
		// Provider:      user.Provider,
		ProviderID:    user.UserID,
		Active:        true,
		Admin:         false,
		PublicProfile: false,
	}
	// TODO register in DB only if not registered ?
	// err = models.DB.Save(&u)
	// if err != nil {
	// 	return c.Error(500, err)
	// }
	c.Set("user", u)

	// // if OK, set username in session
	// c.Session().Set("username", c.Request().FormValue("username"))
	// c.Session().Save()

	// c.Session().Set("AccessToken", user.AccessToken)
	// c.Session().Set("AccessTokenSecret", user.AccessTokenSecret)
	// c.Session().Set("RefreshToken", user.RefreshToken)
	// c.Session().Save()

	return c.Render(200, r.HTML("users/show.html"))
}

func authMiddleware() buffalo.MiddlewareFunc {
	return func(h buffalo.Handler) buffalo.Handler {
		return func(c buffalo.Context) error {
			name := c.Session().Get("username")
			fmt.Println("Name", name)
			if name == "" || name == nil {
				return c.Redirect(302, "/login")
			}
			c.LogField("User", name)
			return h(c)
		}
	}
}
