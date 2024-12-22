package fake

import (
	"fmt"
	"math/rand"
	"time"
)

func FakeTimeRange(start time.Time, end time.Time) time.Time {
	startUnix := start.Unix()
	endUnix := end.Unix()
	unix := startUnix + rand.Int63n(endUnix-startUnix)
	return time.Unix(unix, 0)
}

func FakeYearRange(start, end int) time.Time {
	startTime, _ := time.Parse("2006-01-02", fmt.Sprintf("%d-01-01", start))
	endTime, _ := time.Parse("2006-01-02", fmt.Sprintf("%d-12-31", end))
	return FakeTimeRange(startTime, endTime)
}

func FakeYearsBefore(years int) time.Time {
	endTime := time.Now()
	startTime := endTime.AddDate(-years, 0, 0)
	return FakeTimeRange(startTime, endTime)
}
