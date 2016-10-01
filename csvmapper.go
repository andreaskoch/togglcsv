package main

import (
	"fmt"
	"strings"

	"github.com/andreaskoch/togglapi/date"
	"github.com/andreaskoch/togglcsv/toggl"
)

// The TimeRecordMapper interface provides functions for mapping CSV records to time records and vice versa.
type TimeRecordMapper interface {
	// GetTimeRecords returns a list of time records for the given CSV table rows.
	GetTimeRecords(rows [][]string) ([]toggl.TimeRecord, error)

	// GetTimeRecord returns a TimeRecord model from an CSV row.
	GetTimeRecord(row []string) (toggl.TimeRecord, error)

	// GetColumnNames returns the names of the CSV columns.
	GetColumnNames() []string

	// GetRow returns an CSV row for the given TimeRecord model.
	GetRow(timeRecord toggl.TimeRecord) []string
}

// NewCSVTimeRecordMapper converts CSV rows to TimeRecord models and vice versa.
func NewCSVTimeRecordMapper(dateFormatter date.Formatter) TimeRecordMapper {
	return &CSVTimeRecordMapper{
		dateFormatter: dateFormatter,
		columnNames:   []string{"Start", "Stop", "Workspace Name", "Project Name", "Client Name", "Tag(s)", "Description"},
		tagsSeparator: ",",
	}
}

// CSVTimeRecordMapper converts CSV time records into TimeRecord models.
type CSVTimeRecordMapper struct {
	dateFormatter date.Formatter

	// columnNames contains the list of all CSV column names for CSV-based time reports used for import or export.
	columnNames []string

	// tagsSeparator contains the separator sign/string that is used to split and concatenate tags
	tagsSeparator string
}

// GetTimeRecords returns a list of time records for the given CSV table rows.
func (mapper *CSVTimeRecordMapper) GetTimeRecords(rows [][]string) ([]toggl.TimeRecord, error) {

	// cut the headline
	if len(rows) > 0 {
		firstLine := rows[0]
		firstColumnName := mapper.GetColumnNames()[0]
		if firstLine[0] == firstColumnName {
			rows = rows[1:]
		}
	}

	// create time record models from each row
	var timeRecords []toggl.TimeRecord
	for _, row := range rows {

		timeRecord, timeRecordError := mapper.GetTimeRecord(row)
		if timeRecordError != nil {
			return nil, fmt.Errorf("Failed to create time entry from (%v): %s", row, timeRecordError.Error())
		}

		timeRecords = append(timeRecords, timeRecord)

	}

	return timeRecords, nil
}

// GetTimeRecord returns a TimeRecord model from an CSV row.
func (mapper *CSVTimeRecordMapper) GetTimeRecord(row []string) (toggl.TimeRecord, error) {

	// check the number of columns
	if len(row) != len(mapper.GetColumnNames()) {
		return toggl.TimeRecord{}, fmt.Errorf("Wrong number of values in the given row. The required: %d. Given: %d", len(mapper.GetColumnNames()), len(row))
	}

	// Start date
	startDateVal := row[0]
	startDate, startDateError := mapper.dateFormatter.GetDate(startDateVal)
	if startDateError != nil {
		return toggl.TimeRecord{}, fmt.Errorf("Cannot parse the start date: %s", startDateError)
	}

	// Stop Date
	stopDateVal := row[1]
	stopDate, stopDateError := mapper.dateFormatter.GetDate(stopDateVal)
	if stopDateError != nil {
		return toggl.TimeRecord{}, fmt.Errorf("Cannot parse the stop date: %s", stopDateError)
	}

	// Workspace Name
	workspaceVal := row[2]
	workspaceVal = strings.TrimSpace(workspaceVal)

	// Project Name
	projectVal := row[3]
	projectVal = strings.TrimSpace(projectVal)

	// Client Name
	clientVal := row[4]
	clientVal = strings.TrimSpace(clientVal)

	// Tags
	tagsVal := row[5]
	tags := strings.Split(tagsVal, mapper.tagsSeparator)
	for index, tag := range tags {
		tags[index] = strings.TrimSpace(tag)
	}

	// Description
	descriptionVal := row[6]
	description := strings.TrimSpace(descriptionVal)
	if len(description) >= 3000 {
		return toggl.TimeRecord{}, fmt.Errorf("The description text of the time entry %q is too long", startDate)
	}

	entry := toggl.TimeRecord{
		Start:         startDate,
		Stop:          stopDate,
		WorkspaceName: workspaceVal,
		ProjectName:   projectVal,
		ClientName:    clientVal,
		Description:   description,
		Tags:          tags,
	}

	return entry, nil

}

// GetColumnNames returns the names of the CSV columns.
func (mapper *CSVTimeRecordMapper) GetColumnNames() []string {
	return mapper.columnNames
}

// GetRow returns an CSV row for the given TimeRecord model.
func (mapper *CSVTimeRecordMapper) GetRow(timeRecord toggl.TimeRecord) []string {
	return []string{
		mapper.dateFormatter.GetDateString(timeRecord.Start),
		mapper.dateFormatter.GetDateString(timeRecord.Stop),
		timeRecord.WorkspaceName,
		timeRecord.ProjectName,
		timeRecord.ClientName,
		strings.Join(timeRecord.Tags, mapper.tagsSeparator),
		timeRecord.Description,
	}
}
