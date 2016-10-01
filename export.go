package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"time"

	"github.com/andreaskoch/togglcsv/toggl"
)

// The CSVExporter interface exports time records from a Toggl account as CSV.
type CSVExporter interface {
	// Export prints all time records from the given start date as CSV
	Export(startDate, endDate time.Time, writer io.Writer) error
}

// TogglCSVExporter provides import and export functionality Toggl accounts.
type TogglCSVExporter struct {
	csvMapper            TimeRecordMapper
	timeRecordRepository toggl.TimeRecorder
}

// Export prints all time records from the given start date as CSV.
func (exporter *TogglCSVExporter) Export(startDate, endDate time.Time, writer io.Writer) error {

	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()

	// write the header
	csvWriter.Write(exporter.csvMapper.GetColumnNames())
	csvWriter.Flush()

	records, timeRecordsError := exporter.timeRecordRepository.GetTimeRecords(startDate, endDate)
	if timeRecordsError != nil {
		return fmt.Errorf("Failed to retrieve time records between %q and %q: %s", startDate, endDate, timeRecordsError.Error())
	}

	// write the records one-by-one
	for _, record := range records {
		row := exporter.csvMapper.GetRow(record)
		csvWriter.Write(row)
		csvWriter.Flush()
	}

	return nil
}
