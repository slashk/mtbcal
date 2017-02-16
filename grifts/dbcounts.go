package grifts

import (
	"os"
	"strconv"

	. "github.com/markbates/grift/grift"
	"github.com/olekukonko/tablewriter"
	// "github.com/slashk/mtbcal/actions"
	"github.com/slashk/mtbcal/models"
)

var _ = Add("counts", func(c *Context) error {
	cnt := map[string]int{}

	cnt["Events"], _ = models.DB.Count(&models.Event{})
	cnt["Races"], _ = models.DB.Count(&models.Race{})
	cnt["Users"], _ = models.DB.Count(&models.User{})

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Table", "Records"})
	for k, v := range cnt {
		table.Append([]string{k, strconv.Itoa(v)})
	}
	table.SetCenterSeparator("|")
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.Render()
	return nil
})
