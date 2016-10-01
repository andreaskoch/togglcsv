package togglapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/andreaskoch/togglapi/date"
	"github.com/andreaskoch/togglapi/model"
	"github.com/pkg/errors"
)

// NewTimeEntryAPI create a new client for the Toggl time entry API.
func NewTimeEntryAPI(baseURL, token string) model.TimeEntryAPI {
	return &TimeEntryAPI{
		restClient: &togglRESTAPIClient{
			baseURL: baseURL,
			token:   token,
		},
		dateFormatter: date.NewISO8601Formatter(),
	}
}

// TimeEntryAPI provides functions for interacting with Toggls' time entry API.
type TimeEntryAPI struct {
	restClient    RESTRequester
	dateFormatter date.Formatter
}

// CreateTimeEntry creates a new time entry.
func (repository *TimeEntryAPI) CreateTimeEntry(timeEntry model.TimeEntry) (model.TimeEntry, error) {

	duration := int(timeEntry.Stop.Sub(timeEntry.Start).Seconds())

	timeEntryModel := struct {
		Wid         int       `json:"wid"`
		Pid         int       `json:"pid"`
		Start       time.Time `json:"start"`
		Duration    int       `json:"duration"`
		Billable    bool      `json:"billable"`
		Description string    `json:"description"`
		Tags        []string  `json:"tags"`
		CreatedWith string    `json:"created_with"`
	}{
		Wid:         timeEntry.Wid,
		Pid:         timeEntry.Pid,
		Start:       timeEntry.Start,
		Duration:    duration,
		Billable:    timeEntry.Billable,
		Description: timeEntry.Description,
		Tags:        timeEntry.Tags,
		CreatedWith: clientName,
	}

	// create the request object
	timeEntryRequest := struct {
		TimeEntry interface{} `json:"time_entry"`
	}{
		TimeEntry: timeEntryModel,
	}

	jsonBody, marshalError := json.Marshal(timeEntryRequest)
	if marshalError != nil {
		return model.TimeEntry{}, errors.Wrap(marshalError, "Failed to serialize the time entry")
	}

	content, err := repository.restClient.Request(http.MethodPost, "time_entries", bytes.NewBuffer(jsonBody))
	if err != nil {
		return model.TimeEntry{}, errors.Wrap(err, "Failed to create time entry")
	}

	var timeEntryResponse struct {
		TimeEntry model.TimeEntry `json:"data"`
	}

	if unmarshalError := json.Unmarshal(content, &timeEntryResponse); unmarshalError != nil {
		return model.TimeEntry{}, errors.Wrap(unmarshalError, "Failed to deserialize the time entry")
	}

	return timeEntryResponse.TimeEntry, nil
}

// GetTimeEntries returns all time entries created between the given start and end date.
// Returns nil and an error if the time entries could not be retrieved.
func (repository *TimeEntryAPI) GetTimeEntries(start, end time.Time) ([]model.TimeEntry, error) {
	route := fmt.Sprintf(
		"time_entries?start_date=%s&end_date=%s",
		url.QueryEscape(repository.dateFormatter.GetDateString(start)),
		url.QueryEscape(repository.dateFormatter.GetDateString(end)),
	)

	content, err := repository.restClient.Request(http.MethodGet, route, nil)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Failed to retrieve time entries (Start: %q, Stop: %q)", start, end))
	}

	var timeEntries []model.TimeEntry
	if unmarshalError := json.Unmarshal(content, &timeEntries); unmarshalError != nil {
		return nil, errors.Wrap(unmarshalError, "Failed to deserialize time entries")
	}

	return timeEntries, nil
}
