// Package main contains the Toggl⥃CSV commandline utility for importing
// and exporting Toggl time records.
package main

import (
	"encoding/csv"
	"fmt"
	"io"

	"github.com/andreaskoch/togglcsv/toggl"
	"github.com/pkg/errors"
	"gopkg.in/cheggaaa/pb.v1"
)

// The CSVImporter interface imports CSV records into a Toggl account.
type CSVImporter interface {
	// Import reads time records supplied via Stdin and imports them into a Toggl account.
	Import(input io.Reader) error
}

// TogglCSVImporter provides import and export functionality Toggl accounts.
type TogglCSVImporter struct {
	csvMapper            TimeRecordMapper
	timeRecordRepository toggl.TimeRecorder
	output               io.Writer
}

// Import reads time records supplied via Stdin and imports them into a Toggl account.
func (togglCSVImporter *TogglCSVImporter) Import(input io.Reader) error {

	// read the CSV data
	csvReader := csv.NewReader(input)
	rows, csvError := csvReader.ReadAll()
	if csvError != nil {
		return fmt.Errorf("Failed to read time records from CSV: %s", csvError.Error())
	}

	timeRecords, timeRecordsError := togglCSVImporter.csvMapper.GetTimeRecords(rows)
	if timeRecordsError != nil {
		return timeRecordsError
	}

	// abort if no time records were returned
	if len(timeRecords) == 0 {
		return nil
	}

	// upload the time entries to toggl
	progressbar := pb.New(len(timeRecords))
	progressbar.ShowTimeLeft = true
	if togglCSVImporter.output != nil {
		progressbar.Output = togglCSVImporter.output
		progressbar.Start()
	}

	// create the records
	for recordIndex, record := range timeRecords {

		if err := togglCSVImporter.timeRecordRepository.CreateTimeRecord(record); err != nil {
			return errors.Wrap(err, fmt.Sprintf("Failed to create time record %d of %d", recordIndex+1, len(timeRecords)))
		}

		if togglCSVImporter.output != nil {
			progressbar.Increment()
		}

	}

	if togglCSVImporter.output != nil {
		// FinishPrint writes to os.Stdout
		// see:
		// https://github.com/cheggaaa/pb/issues/87
		// https://github.com/cheggaaa/pb/commit/7f4253899ba18226b3c52aca004d298182360edc#commitcomment-18923803
		// progressbar.FinishPrint("Import complete.")
		progressbar.Finish()
	}

	return nil
}
