package actions

import (
	"net/http"
	"time"

	rice "github.com/GeertJohan/go.rice"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/buffalo/render/resolvers"
)

var r *render.Engine

func init() {
	r = render.New(render.Options{
		HTMLLayout:     "application.html",
		CacheTemplates: ENV == "production",
		FileResolverFunc: func() resolvers.FileResolver {
			return &resolvers.RiceBox{
				Box: rice.MustFindBox("../templates"),
			}
		},
		Helpers: map[string]interface{}{
			"formatDate": formatDate,
		},
	})
}

func assetsPath() http.FileSystem {
	box := rice.MustFindBox("../public/assets")
	return box.HTTPBox()
}

// formatDate mimics the format_date rails helper
func formatDate(t time.Time, f string) string {
	// time = Time.now                    # => Thu Jan 18 06:10:17 CST 2007
	// time.to_formatted_s(:time)         # => "06:10"
	// time.to_s(:time)                   # => "06:10"
	// time.to_formatted_s(:db)           # => "2007-01-18 06:10:17"
	// time.to_formatted_s(:number)       # => "20070118061017"
	// time.to_formatted_s(:short)        # => "18 Jan 06:10"
	// time.to_formatted_s(:long)         # => "January 18, 2007 06:10"
	// time.to_formatted_s(:long_ordinal) # => "January 18th, 2007 06:10"
	// time.to_formatted_s(:rfc822)       # => "Thu, 18 Jan 2007 06:10:17 -0600"
	// time.to_formatted_s(:iso8601)      # => "2007-01-18T06:10:17-06:00"
	// golang reference time is Mon Jan 2 15:04:05 MST 2006
	switch f {
	case "time":
		return t.Format("15:04")
	case "db":
		return t.Format("2006-01-02 15:04:05")
	case "number":
		return t.Format("20060102150405")
	case "short":
		return t.Format("Jan 2 15:04")
	case "long":
		return t.Format("January 2, 2006 15:04")
	case "long_ordinal":
		// humanize 02 then regex it into :long
		// return t.Format("")
		return "unimplemented -- choose another format"
	case "rfc822":
		return t.Format("Mon, 2 Jan 2006 15:04:05 -0500")
	case "iso8601":
		return t.Format("2006-01-02T15:04:05-05:00")
	case "common":
		return t.Format("Jan 2 2006")
	}
	return ""
}
