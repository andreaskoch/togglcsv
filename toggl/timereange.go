package toggl

import (
	"fmt"
	"time"

	"github.com/jinzhu/now"
	"github.com/pkg/errors"
)

// A timeRangeProvider interface calculates time ranges between a given start and end date.
type timeRangeProvider interface {
	// GetTimeRanges returns all time ranges between the given start and end date.
	// Returns an error if no ranges can be calculated for the given start and end date.
	GetTimeRanges(startDate, endDate time.Time) ([]timeRange, error)
}

// fullMonthTimeRangeProvider returns time ranges for a given start and end date.
type fullMonthTimeRangeProvider struct {
}

// GetTimeRanges returns a set of timeRange objects starting
// the given year and month up until the end of the current month.
// Returns an error if no ranges can be calculated for the given start and end date.
func (fullMonthTimeRangeProvider) GetTimeRanges(startDate, endDate time.Time) ([]timeRange, error) {

	// normalize the time of the start and end date
	start := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 1, 0, time.UTC)
	end := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 0, time.UTC)

	// validate start and end date
	if startDate.After(endDate) {
		return nil, fmt.Errorf("The start date cannot be before the end date")
	}

	var ranges []timeRange
	month := start
	for month.Before(end) {

		// determine the start
		var monthStart time.Time
		isFirstMonthInRange := month.Year() == start.Year() && month.Month() == start.Month()
		if isFirstMonthInRange {
			// respect the day if it is the first month range
			monthStart = start
		} else {
			// use the beginning of the month for all following ranges
			monthStart = now.New(month).BeginningOfMonth()
		}

		// determine the end
		var monthEnd time.Time
		isLastMonthInRange := month.Year() == end.Year() && month.Month() == end.Month()
		if isLastMonthInRange {
			// use the given end date
			monthEnd = end
		} else {
			// use the end of the month
			monthEnd = now.New(month).EndOfMonth()
		}

		monthRange, newRangeError := newTimeRange(monthStart, monthEnd)
		if newRangeError != nil {
			return nil, errors.Wrap(newRangeError, fmt.Sprintf("Failed to calculate time ranges between %q and %q", startDate, endDate))
		}

		ranges = append(ranges, monthRange)

		// make sure the next month starts at the beginning
		if month.Day() != 1 {
			month = now.New(month).BeginningOfMonth()
		}

		// next month
		month = month.AddDate(0, 1, 0)
	}

	return ranges, nil
}

// newTimeRange creates a new timeRange instance.
func newTimeRange(start, end time.Time) (timeRange, error) {
	if start.After(end) {
		return timeRange{}, fmt.Errorf("The start date cannot be before the end date")
	}

	if end.Equal(start) {
		return timeRange{}, fmt.Errorf("The start and end date cannot be the same")
	}

	return timeRange{start, end}, nil
}

// A timeRange defines a time range based on a given start and stop date.
type timeRange struct {
	start time.Time
	stop  time.Time
}

// Start returns the start date and time of the range.
func (timeRange timeRange) Start() time.Time {
	return timeRange.start
}

// Stop returns the stop date and time of the range.
func (timeRange timeRange) Stop() time.Time {
	return timeRange.stop
}
