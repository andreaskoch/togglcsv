// Package model contains the models and interface for the Toggl API
package model

import "time"

// Project defines the key properties of a Toggl project
type Project struct {
	ID          int    `json:"id"`
	WorkspaceID int    `json:"wid"`
	ClientID    int    `json:"cid"`
	Name        string `json:"name"`
}

// Workspace defines the key properties of a Toggl workspace
type Workspace struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// TimeEntry represents a single Toggle time tracking record
type TimeEntry struct {

	// ID contains the id of the time entry
	ID int `json:"id"`

	// Wid contains the workspace ID
	Wid int `json:"wid"`

	// Pid contains the project id
	Pid int `json:"pid"`

	// Start contains the start time of the entry.
	Start time.Time `json:"start"`

	// Stop contains the end time of the entry.
	Stop time.Time `json:"stop"`

	// Billable contains a flag indicating whether this time entry is billable or not.
	Billable bool `json:"billable"`

	// The Description of the time entry
	Description string `json:"description"`

	// Tags contains a list of names for the time entry
	Tags []string `json:"tags"`

	CreatedWith string `json:"created_with"`
}

// Client defines the key properties of a Toggl client
type Client struct {
	ID          int    `json:"id"`
	WorkspaceID int    `json:"wid"`
	Name        string `json:"name"`
	Notes       string `json:"notes"`
}
