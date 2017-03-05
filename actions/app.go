package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/middleware"
	"github.com/gobuffalo/envy"
	"github.com/slashk/mtbcal/models"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
func App() *buffalo.App {
	if app == nil {
		app = buffalo.Automatic(buffalo.Options{
			Env:         ENV,
			SessionName: "_mtbcal_session",
		})

		app.Use(middleware.PopTransaction(models.DB))
		app.Use(AuthMiddleware)
		// skip the AuthorizeUser middleware for the login and Logout handlers
		// app.Middleware.Skip(AuthMiddleware, LogoutHandler, LoginHandler)

		app.GET("/login", LoginHandler)
		app.GET("/logout", LogoutHandler)
		app.GET("/", HomeHandler)

		app.ServeFiles("/assets", assetsPath())
	}

	return app
}
