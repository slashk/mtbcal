package grifts

import (
	"log"
	"os"
	"time"

	. "github.com/markbates/grift/grift"
	"github.com/olekukonko/tablewriter"
	"github.com/slashk/mtbcal/models"
)

var _ = Add("fixtures", func(c *Context) error {
	one := models.Event{
		Name:         "BigBoy",
		Location:     "Folsom, CA",
		State:        "CA",
		WebReg:       true,
		Active:       true,
		PublishedAt:  time.Now(),
		StartDate:    time.Now(),
		EndDate:      time.Now(),
		RegOpenDate:  time.Now(),
		RegCloseDate: time.Now(),
		Description:  "This is the bigboy event",
		URL:          "http://www.mtbcalendar.com",
		UserID:       "admin",
		Country:      "US",
		Lng:          95.2353,
		Lat:          38.9717,
	}
	two := models.Event{
		Name:         "LittleBoy",
		Location:     "Boulder, CO",
		State:        "CO",
		WebReg:       false,
		Active:       false,
		PublishedAt:  time.Now(),
		StartDate:    time.Now().Add(300 * time.Hour),
		EndDate:      time.Now(),
		RegOpenDate:  time.Now(),
		RegCloseDate: time.Now(),
		Description:  "This is the bigboy event",
		URL:          "http://www.mtbcalendar.com",
		UserID:       "boss",
		Country:      "US",
		Lng:          95.2353,
		Lat:          38.9717,
	}
	longone := models.Event{
		Name:         "The Otway Odyssey and Great Otway Gravel Grind (GOGG)",
		Location:     "Forrest VIC 3236, Australia",
		State:        "VIC",
		WebReg:       false,
		Active:       false,
		PublishedAt:  time.Now(),
		StartDate:    time.Now().Add(600 * time.Hour),
		EndDate:      time.Now(),
		RegOpenDate:  time.Now(),
		RegCloseDate: time.Now(),
		Description: `The Otway Odyssey presented by Focus is one of Australia’s most respected off-road cycling events thanks to the quality of the courses with lung busting hills, blitzing descents and an electrifying race atmosphere that welcomes any rider – mountain bikers, roadies and gravel grinders.
For 10 years, riders have travelled to the famous Otway rainforest of Forrest in south-western Victoria to take on the region’s iconic single track, but in 2017 even more riders can be part of the action with the addition of the Great Otway Gravel Grind (The GOGG) 49km and 97km races to the schedule of events.
The action gets underway on Saturday 25th February with the principle mountain bike races; the 100km Odyssey, the 50km Shorty and the 30km Rookie, all starting and finishing amidst a huge bike expo and food festival at the Forrest Football Ground. The GOGG races and the kids 10km Pioneer MTB will be held on Sunday 26th February.

A Gravel Grind is an event held on unsealed 2WD roads and the Otway Ranges around Forrest are infused with these, making for some superb riding on traffic free roads through the wilderness. The GOGG is achievable for all styles of rider on either road bikes, cyclocross, mountain or hybrid bikes.

A *unique* feature of The GOGG races will be time-out sections where riders can stop for a coffee or visit a café during their race. Many will also stop and wait for their friends, adding to the social aspects of the event especially for the long course riders where the time-out zone will be in the coastal townships of Kennett River and Wye River on the spectacular Great Ocean Road.
The Otway Odyssey always attracts the best riders from across Australia who come not only to chase the titles but also to get their hands on the prizemoney which includes over $7,000 in cash across the MTB races and $3,700 up for grabs in The GOGG.


The Otway Odyssey MTB Marathons have grown to become the pre-eminent mountain bike races in Australia and on the bucket list for many cyclists. The courses are well known for their tough climbs and technical trails, but also for the friendliness on course and the race atmosphere around the Forrest Football Ground expo area.
An extensive event expo with sponsors stands, food, drinks, entertainment and activities at the Start / Finish makes this a weekend when all riders come together to celebrate off-road riding and gravel grinding in the best of company. So mark it in your diary, get on your bike and we’ll see you at Forrest over 25th and 26th February, 2017.`,
		URL:     "http://www.mtbcalendar.com",
		UserID:  "boss",
		Country: "US",
		Lng:     95.2353,
		Lat:     38.9717,
	}
	events := models.Events{}
	events = append(events, one)
	events = append(events, two)
	events = append(events, longone)
	for x := range events {
		err := models.DB.Create(&events[x])
		if err != nil {
			log.Panic(err.Error())
		}
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Date"})
	for _, e := range events {
		table.Append([]string{e.Name, e.StartDate.Format("1 Feb 2016")})
	}
	table.SetCenterSeparator("|")
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.Render()
	return nil
})
