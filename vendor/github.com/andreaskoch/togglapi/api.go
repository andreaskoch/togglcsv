// Package togglapi provides access to Toggls' time tracking API.
// The togglapi package provides functions for creating and retrieving
// workspaces, clients, projects and time entries.
//
// To learn more about the Toggl API visit:
// https://github.com/toggl/toggl_api_docs
package togglapi

import (
	"time"

	"github.com/andreaskoch/togglapi/date"
	"github.com/andreaskoch/togglapi/model"
)

const clientName = "github.com/andreaskoch/togglapi"

// NewAPI create a new instance of the Toggl API.
func NewAPI(baseURL, token string) model.TogglAPI {
	restAPI := &togglRESTAPIClient{
		baseURL: baseURL,
		token:   token,

		// The Toggl API only allows roughly 1 request per second
		// see: https://github.com/toggl/toggl_api_docs
		pauseBetweenRequests: time.Millisecond * 1000,
	}

	dateFormatter := date.NewISO8601Formatter()

	return &API{
		&WorkspaceAPI{restAPI},
		&ProjectAPI{restAPI},
		&TimeEntryAPI{restAPI, dateFormatter},
		&ClientAPI{restAPI},
	}
}

// API provides functions for interacting with the Toggl API.
type API struct {
	model.WorkspaceAPI
	model.ProjectAPI
	model.TimeEntryAPI
	model.ClientAPI
}
