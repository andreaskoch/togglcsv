package model

import "time"

// The ProjectAPI interface provides functions for creating and fetching projects.
type ProjectAPI interface {
	// CreateProject creates a new project.
	CreateProject(project Project) (Project, error)

	// GetProjects returns all projects for the given workspace.
	GetProjects(workspaceID int) ([]Project, error)
}

// The ClientAPI interface provides functions for creating and fetching clients.
type ClientAPI interface {
	// CreateClient creates a new client.
	CreateClient(client Client) (Client, error)

	// GetClients returns all clients.
	GetClients() ([]Client, error)
}

// The WorkspaceAPI interface provides functions for fetching workspacs.
type WorkspaceAPI interface {
	// GetWorkspaces returns all workspaces for the current user.
	GetWorkspaces() ([]Workspace, error)
}

// The TimeEntryAPI interface provides functions for fetching and creating time entries.
type TimeEntryAPI interface {
	// CreateTimeEntry creates a new time entry.
	CreateTimeEntry(timeEntry TimeEntry) (TimeEntry, error)

	// GetTimeEntries returns all time entries created between the given start and end date.
	// Returns nil and an error if the time entries could not be retrieved.
	GetTimeEntries(start, end time.Time) ([]TimeEntry, error)
}

// A TogglAPI interface implements some of the Toggl API methods.
type TogglAPI interface {
	WorkspaceAPI
	ProjectAPI
	TimeEntryAPI
	ClientAPI
}
