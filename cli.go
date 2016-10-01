// Package main contains the Toggl⥃CSV commandline utility for importing
// and exporting Toggl time records.
package main

import (
	"fmt"
	"io"
	"time"

	"github.com/jinzhu/now"

	"gopkg.in/alecthomas/kingpin.v2"
)

// exportDateFormat defines the date format for start and end dates of the export command.
const exportDateFormat = "2006-01-02"

func init() {
	now.FirstDayMonday = true
}

type togglCli struct {
	importerFactory func(apiToken string) CSVImporter
	exporterFactory func(apiToken string) CSVExporter
}

// Execute parses the given arguments and performs the selected action.
func (cli *togglCli) Execute(input io.Reader, output, errorOutput io.Writer, args []string) (success bool) {
	app := kingpin.New(applicationName, "Toggl⥃CSV is an csv-based import/export utility for Toggl time tracking data (see: https://github.com/andreaskoch/togglcsv)")
	app.Version(applicationVersion)
	app.Writer(errorOutput)
	app.Terminate(func(int) {
		return
	})

	// export
	exportCommand := app.Command("export", "Export your Toggl time tracking records as CSV")
	exportAPIToken := exportCommand.Arg("token", "The Toggl API token of the source account").Required().String()
	exportStartDate := exportCommand.Arg("startdate", "The start date (e.g. \"2006-01-26\")").Required().String()
	exportEndDate := exportCommand.Arg("enddate", "The start date (e.g. \"2006-01-26\")").String()

	// import
	importCommand := app.Command("import", "Import CSV-based time tracking records into Toggl from stdin")
	importAPIToken := importCommand.Arg("token", "The Toggl API token of the target account").Required().String()

	command, err := app.Parse(args)
	if err != nil {
		app.Fatalf("%s", err.Error())
	}

	switch command {

	// export
	case exportCommand.FullCommand():

		// start date (required)
		startDate, startDateError := time.Parse(exportDateFormat, *exportStartDate)
		if startDateError != nil {
			app.Fatalf("Failed to parse the given start date %q. %s", *exportStartDate, startDateError.Error())
		}

		// end date (optional)
		now := time.Now()
		endDate := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.UTC) // use current date as the default

		if len(*exportEndDate) > 0 {
			endDateParsed, endDateError := time.Parse(exportDateFormat, *exportEndDate)
			if endDateError != nil {
				app.Fatalf("Failed to parse the given end date %q. %s", *exportEndDate, endDateError.Error())
			}

			endDate = endDateParsed
		}

		exporter := cli.exporterFactory(*exportAPIToken)
		if exportError := exporter.Export(startDate, endDate, output); exportError != nil {
			fmt.Fprintf(errorOutput, "Error: %s\n", exportError.Error())
			return false
		}

		return true

	// import
	case importCommand.FullCommand():
		importer := cli.importerFactory(*importAPIToken)
		if importError := importer.Import(input); importError != nil {
			fmt.Fprintf(errorOutput, "Error: %s\n", importError.Error())
			return false
		}

		return true

	}

	return false
}
