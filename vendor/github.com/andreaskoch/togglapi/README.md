# Toggl API for go

A go wrapper for the Toggl API

`github.com/andreaskoch/togglapi` is a simple wrapper for the Toggl API (https://github.com/toggl/toggl_api_docs).

[![Build Status](https://travis-ci.org/andreaskoch/togglapi.svg?branch=master)](https://travis-ci.org/andreaskoch/togglapi)

## Installation

Download `github.com/andreaskoch/togglapi`:

```bash
go get github.com/andreaskoch/togglapi
```

## Supported API methods

`github.com/andreaskoch/togglapi` currently only supports the following methods of the Toggl API:

- Clients
	- `CreateClient(client Client) (Client, error)`
	- `GetClients() ([]Client, error)`
- Workspaces
	- `GetWorkspaces() ([]Workspace, error)`
- Projects
	- `CreateProject(project Project) (Project, error)`
	- `GetProjects(workspaceID int) ([]Project, error)`
- Time Entries
	- `CreateTimeEntry(timeEntry TimeEntry) (TimeEntry, error)`
	- `GetTimeEntries(start, end time.Time) ([]TimeEntry, error)`

I might add the missing methods in the future, but if you need them now please add them and send me a pull-request.

## Usage

```go
package main

import (
	"time"

	"github.com/andreaskoch/togglapi"
)

func main() {
	apiToken := "Your-API-Token"
	baseURL := "https://www.toggl.com/api/v8"
	api := togglapi.NewAPI(baseURL, apiToken)

  // workspaces
	workspaces, workspacesError := api.GetWorkspaces()
	...

  // clients
	clients, clientsError := api.GetClients()
  ...

  // projects by workspace
  for _, workspace := range workspaces {
		projects, projectsError := api.GetProjects(workspace.ID)
		...
	}

	// time entries
	stop := time.Now()
	start := stop.AddDate(0, -1, 0)
	timeEntries, timeEntriesError := api.GetTimeEntries(start, stop)
  ...
}
```

You can also have a look at the **example command line utility**: [example/main.go](example/main.go)

```bash
cd $GOPATH/src/github.com/andreaskoch/togglapi/example
go run main.go Your-Toggl-API-Token
```

## Development

Run the unit tests:

```bash
cd $GOPATH/src/github.com/andreaskoch/togglapi
make test
```

Create code coverage reports:

```bash
cd $GOPATH/src/github.com/andreaskoch/togglapi
make coverage
```

## Licensing

TogglCSV is licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE) for the full license text.
