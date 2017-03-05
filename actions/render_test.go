package actions_test

import (
	"testing"
	"time"

	"github.com/slashk/mtbcal/actions"
	"github.com/stretchr/testify/require"
)

func Test_dateRange(t *testing.T) {
	r := require.New(t)
	s := time.Now()
	e := s
	r.Equal(actions.FormatDate(s, "common"), actions.DateRange(s, e))
	e = time.Now().Add(3 * time.Hour)
	r.NotEqual(actions.FormatDate(s, "common"), actions.DateRange(s, e))
	s = time.Date(1968, 11, 1, 1, 0, 0, 0, time.UTC)
	e = time.Date(1968, 11, 5, 13, 0, 0, 0, time.UTC)
	r.Equal("Nov 1 1968 - Nov 5 1968", actions.DateRange(s, e))
}

func Test_FormatDate(t *testing.T) {
	r := require.New(t)
	d := time.Date(1968, 11, 1, 1, 0, 0, 0, time.UTC)
	r.Equal("Nov 1 1968", actions.FormatDate(d, "common"), "Old common date not formatted correctly")
	r.Equal("1968-11-01", actions.FormatDate(d, "datepicker"), "Datepicker date not formatted correctly")
	d = time.Date(1968, 11, 1, 13, 10, 10, 0, time.UTC)
	r.Equal("Nov 1 1968", actions.FormatDate(d, "common"), "Dates (common) with time not formatted correctly")
	r.Equal("1968-11-01", actions.FormatDate(d, "datepicker"), "Dates (datepicker) with time not formatted correctly")
}
