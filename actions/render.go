package actions

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	rice "github.com/GeertJohan/go.rice"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/buffalo/render/resolvers"
	"github.com/slashk/mtbcal/models"
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
			"formatDate":   FormatDate,
			"formatLatLng": formatLatLng,
			"formatBool":   formatBool,
			"raceName":     formatRaceName,
			"dateRange":    DateRange,
		},
	})
}

func assetsPath() http.FileSystem {
	box := rice.MustFindBox("../public/assets")
	return box.HTTPBox()
}

func formatRaceName(i int) string {
	return models.FormatbyID[i].Name
}

// FormatDate mimics the format_date rails helper
func FormatDate(t time.Time, f string) string {
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
	case "year":
		return t.Format("2006")
	case "datepicker":
		return t.Format("2006-01-02")
	}
	return t.Format("Jan 2 2006")
}

// formatLatLng returns a formatted string to 3 decimals from a float32
func formatLatLng(f float32) string {
	return strconv.FormatFloat(float64(f), 'f', 3, 64)
}

// formatBool return "true" or "false" strings
func formatBool(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// DateRange formats a range of dates for humans
func DateRange(start time.Time, end time.Time) string {
	f := "common"
	if start == end {
		return FormatDate(start, f)
	}
	return fmt.Sprintf("%s - %s", FormatDate(start, f), FormatDate(end, f))
}
