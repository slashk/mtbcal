package actions

import (
	"errors"
	"fmt"
	"github.com/gobuffalo/buffalo"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/instagram"
	"github.com/markbates/goth/providers/twitter"
	"github.com/slashk/mtbcal/models"
	"net/http"
	"os"
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
	c.LogField("response", c.Response().Header())
	user, err := gothic.CompleteUserAuth(c.Response(), c.Request())
	if err != nil {
		return c.Error(401, err)
	}

	// TODO find user and create session
	// TODO if not found, register the user and create session

	c.LogField("User", user)

	// register
	u := models.User{
		Login:         user.Name,
		Hometown:      user.Location,
		Email:         user.Email,
		Avatar:        user.AvatarURL,
		Provider:      user.Provider,
		ProviderID:    user.UserID,
		Active:        true,
		Admin:         false,
		PublicProfile: false,
	}

	// TODO lookup user via provider and providerID
	// TODO register in DB only if not found
	// err = models.DB.Save(&u)
	// if err != nil {
	// 	return c.Error(500, err)
	// }
	c.Set("user", u)

	c.Session().Set("username", u.Login)
	c.Session().Set("AccessToken", user.AccessToken)
	c.Session().Set("AccessTokenSecret", user.AccessTokenSecret)
	c.Session().Set("RefreshToken", user.RefreshToken)
	c.Session().Set("Provider", user.Provider)
	c.Session().Save()

	c.Set("username", u.Login)

	return c.Render(200, r.HTML("users/show.html"))
}

// LoginHandler renders the login form for /login
func LoginHandler(c buffalo.Context) error {
	c.Set("page", pageDefault)
	return c.Render(200, r.HTML("login.html"))
}

// LogoutHandler logs the user out
func LogoutHandler(c buffalo.Context) error {
	// gothic.GetProviderName = retrieveMTBcalProviderName
	// err := gothic.Logout(c.Response(), c.Request())
	err := deleteGothSession(c)
	if err != nil {
		c.LogField("logout_failure", err)
	}
	c.Session().Clear()
	c.Session().Session.Options.MaxAge = -1
	c.Session().Save()
	// NOTE do not redirect 301 here ... cached by browser
	return c.Render(200, r.HTML("login.html"))
}

// AuthMiddleware grabs current session
//
// Shamelessly stolen from https://godoc.org/github.com/gobuffalo/buffalo#MiddlewareStack.Skip
func AuthMiddleware(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		loggedin := false
		admin := false
		name := c.Session().Get("username")
		if name != nil {
			loggedin = true
		}
		c.LogField("session_user", name)
		if loggedin && isAdmin(name.(string)) {
			admin = true
		}
		c.Set("loggedin", loggedin)
		c.Set("admin", admin)
		err := next(c)
		return err
	}
}

func isAdmin(u string) bool {
	// TODO lookup user in DB or stuff in session
	if u == "Ken Pepple" {
		return true
	}
	return false
}

func retrieveMTBcalProviderName(req *http.Request) (string, error) {
	return "twitter", nil
}

func deleteGothSession(c buffalo.Context) error {
	p := c.Session().Get("Provider")
	if p == nil {
		return errors.New("Could not delete user session ")
	}
	session, err := gothic.Store.Get(c.Request(), p.(string)+gothic.SessionName)
	if err != nil {
		return errors.New("Could not delete user session ")
	}
	session.Options.MaxAge = -1
	session.Values = make(map[interface{}]interface{})
	err = session.Save(c.Request(), c.Response())
	if err != nil {
		return errors.New("Could not delete user session ")
	}
	return nil
}
