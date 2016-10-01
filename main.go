// Package main contains the Toggl⥃CSV commandline utility for importing
// and exporting Toggl time records.
package main

import (
	"io"
	"os"

	"github.com/andreaskoch/togglapi"
	"github.com/andreaskoch/togglapi/date"
	"github.com/andreaskoch/togglcsv/toggl"
	"github.com/jinzhu/now"
)

const applicationName = "togglcsv"
const applicationVersion = "1.0.0-dev"
const togglAPIBaseURL = "https://www.toggl.com/api/v8"

var out io.Writer
var err io.Writer
var in io.Reader
var args []string

func init() {
	now.FirstDayMonday = true

	out = os.Stdout
	err = os.Stderr
	in = os.Stdin
	args = os.Args[1:]
}

func main() {

	cli := togglCli{
		importerFactory: getCSVImporter,
		exporterFactory: getCSVExporter,
	}

	cli.Execute(in, out, err, args)
}

// getCSVExporter creates a new CSVExporter instance for the given API token.
func getCSVExporter(apiToken string) CSVExporter {
	dateFormatter := date.NewISO8601Formatter()
	csvTimeRecordMapper := NewCSVTimeRecordMapper(dateFormatter)

	togglAPI := togglapi.NewAPI(togglAPIBaseURL, apiToken)
	workspaces := toggl.NewWorkspaceRepository(togglAPI)
	clients := toggl.NewClientRepository(togglAPI, workspaces)
	projects := toggl.NewProjectRepository(togglAPI, workspaces, clients)

	timeRecords := toggl.NewTimeRecordRepository(togglAPI, workspaces, projects, clients)

	return &TogglCSVExporter{
		csvMapper:            csvTimeRecordMapper,
		timeRecordRepository: timeRecords,
	}
}

// getCSVImporter creates a new CSVImporter instance for the given API token.
func getCSVImporter(apiToken string) CSVImporter {
	dateFormatter := date.NewISO8601Formatter()
	csvTimeRecordMapper := NewCSVTimeRecordMapper(dateFormatter)

	togglAPI := togglapi.NewAPI(togglAPIBaseURL, apiToken)
	workspaces := toggl.NewWorkspaceRepository(togglAPI)
	clients := toggl.NewClientRepository(togglAPI, workspaces)
	projects := toggl.NewProjectRepository(togglAPI, workspaces, clients)

	timeRecords := toggl.NewTimeRecordRepository(togglAPI, workspaces, projects, clients)

	return &TogglCSVImporter{
		csvMapper:            csvTimeRecordMapper,
		timeRecordRepository: timeRecords,
		output:               os.Stdout,
	}
}
