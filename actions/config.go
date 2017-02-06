package actions

// PageDefaults are the defaults for every html page
type PageDefaults struct {
	Title       string
	Description string
	Keywords    string
}

// D is the default instantiation of PageDefaults
var pageDefault = PageDefaults{
	Title:       "MTBCal",
	Description: "The best mtb calendar",
	Keywords:    "mtb, mountain bike, xc",
}
