// Package date provides functions for parsing and formatting
// dates according the requirements of Toggl (ISO 8601).
package date

import "time"

// A Formatter interface provides functions for parsing and formatting dates.
type Formatter interface {
	// GetDateString returns a formatted date string for the given date.
	GetDateString(date time.Time) string

	// GetDate returns a time.Time model for a given date string.
	// Returns and error if the date could not be parsed.
	GetDate(date string) (time.Time, error)
}

// NewISO8601Formatter creates a new ISO 8601 date formatter.
func NewISO8601Formatter() Formatter {
	return &iso80601Formatter{}
}

// iso8601DateFormat contains the date format for ISO 8601
const iso8601DateFormat = "2006-01-02T15:04:05-07:00"

// iso80601Formatter parses and formats ISO 8061 dates.
type iso80601Formatter struct {
}

// GetDateString returns an ISO 8601 formatted date from the given time.Time object.
func (iso80601Formatter) GetDateString(date time.Time) string {
	return date.Format(iso8601DateFormat)
}

// GetDate returns a time.Time model for the given date ISO 8601 date string.
// Returns an error of the date could not be parsed.
func (iso80601Formatter) GetDate(date string) (time.Time, error) {
	return time.Parse(iso8601DateFormat, date)
}
