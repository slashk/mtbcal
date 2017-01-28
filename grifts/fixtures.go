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
	}
	events := models.Events{}
	events = append(events, one)
	events = append(events, two)
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
