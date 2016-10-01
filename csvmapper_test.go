package main

import (
	"strings"
	"testing"

	"github.com/andreaskoch/togglapi/date"
	"github.com/andreaskoch/togglcsv/toggl"
)

func Test_GetTimeRecord_InvalidNumberOfColumns_ErrorIsReturned(t *testing.T) {
	// arrange
	dateFormatter := date.NewISO8601Formatter()
	csvMapper := &CSVTimeRecordMapper{
		dateFormatter: dateFormatter,
		columnNames:   []string{"Start", "Stop", "Workspace Name", "Project Name", "Client Name", "Tag(s)", "Description"},
		tagsSeparator: ",",
	}

	inputRows := [][]string{

		// Too many
		[]string{"2015-03-26T08:53:56Z", "2015-03-26T12:55:10+01:00", "Workspace", "Project XY", "Client X", "", "Some stuff"},

		// Too few
		[]string{"Project XY", "Client X", "", "Some stuff"},
	}

	for _, row := range inputRows {
		// act
		_, err := csvMapper.GetTimeRecord(row)

		// assert
		if err == nil {
			t.Fail()
			t.Logf("GetTimeRecord should return an error the number of columns is wrong")
		}
	}
}

func Test_GetTimeRecord_InvalidStartDateFormat_ErrorIsReturned(t *testing.T) {
	// arrange
	dateFormatter := date.NewISO8601Formatter()
	csvMapper := &CSVTimeRecordMapper{
		dateFormatter: dateFormatter,
		columnNames:   []string{"Start", "Stop", "Workspace Name", "Project Name", "Client Name", "Tag(s)", "Description"},
		tagsSeparator: ",",
	}

	row := []string{"2015-03-26", "2015-03-26T12:55:10+01:00", "Workspace", "Project XY", "Client X", "", "Some stuff"}

	// act
	_, err := csvMapper.GetTimeRecord(row)

	// assert
	if err == nil {
		t.Fail()
		t.Logf("GetTimeRecord should return an error if the given start date format is invalid")
	}
}

func Test_GetTimeRecord_InvalidStopDateFormat_ErrorIsReturned(t *testing.T) {
	// arrange
	dateFormatter := date.NewISO8601Formatter()
	csvMapper := &CSVTimeRecordMapper{
		dateFormatter: dateFormatter,
		columnNames:   []string{"Start", "Stop", "Workspace Name", "Project Name", "Client Name", "Tag(s)", "Description"},
		tagsSeparator: ",",
	}

	row := []string{"2015-03-26T08:00:00+01:00", "2015-03-26", "Workspace", "Project XY", "Client X", "", "Some stuff"}

	// act
	_, err := csvMapper.GetTimeRecord(row)

	// assert
	if err == nil {
		t.Fail()
		t.Logf("GetTimeRecord should return an error if the given stop date format is invalid")
	}
}

func Test_GetTimeRecord_TooLongDescription_ErrorIsReturned(t *testing.T) {
	// arrange
	dateFormatter := date.NewISO8601Formatter()
	csvMapper := &CSVTimeRecordMapper{
		dateFormatter: dateFormatter,
		columnNames:   []string{"Start", "Stop", "Workspace Name", "Project Name", "Client Name", "Tag(s)", "Description"},
		tagsSeparator: ",",
	}

	oneHundredCharacters := `Lorem ipsum dolor sit amet, consectetuer adipiscing elit. Aenean commodo ligula eget dolor. Aenean m`
	veryLongDescription := ""
	for i := 0; i < 40; i++ {
		veryLongDescription = veryLongDescription + oneHundredCharacters
	}

	row := []string{"2015-03-26T08:00:00+01:00", "2015-03-26T11:00:00+01:00", "Workspace", "Project XY", "Client X", "", veryLongDescription}

	// act
	_, err := csvMapper.GetTimeRecord(row)

	// assert
	if err == nil {
		t.Fail()
		t.Logf("GetTimeRecord should return an error if the given description is too long")
	}
}

func Test_GetTimeRecord_ValidInput_NoErrorIsReturned(t *testing.T) {
	// arrange
	dateFormatter := date.NewISO8601Formatter()
	csvMapper := &CSVTimeRecordMapper{
		dateFormatter: dateFormatter,
		columnNames:   []string{"Start", "Stop", "Workspace Name", "Project Name", "Client Name", "Tag(s)", "Description"},
		tagsSeparator: ",",
	}

	row := []string{"2015-03-26T08:00:00+01:00", "2015-03-26T11:30:00+01:00", "Workspace", "Project XY", "Client X", "", "Some stuff"}

	// act
	_, err := csvMapper.GetTimeRecord(row)

	// assert
	if err != nil {
		t.Fail()
		t.Logf("GetTimeRecord have returned a time entry but returned an error instead: %s", err)
	}
}

func Test_GetTimeRecord_ValidInput_TimeRecordIsReturned(t *testing.T) {
	// arrange
	dateFormatter := date.NewISO8601Formatter()
	csvMapper := &CSVTimeRecordMapper{
		dateFormatter: dateFormatter,
		columnNames:   []string{"Start", "Stop", "Workspace Name", "Project Name", "Client Name", "Tag(s)", "Description"},
		tagsSeparator: ",",
	}

	row := []string{"2015-03-26T08:00:00+01:00", "2015-03-26T11:30:00+01:00", "Workspace", "Project XY", "Client X", "", "Some stuff"}

	// act
	timeRecord, _ := csvMapper.GetTimeRecord(row)

	// assert
	if timeRecord.WorkspaceName != "Workspace" {
		t.Fail()
		t.Logf("The workspace name was not returned correctly: %#v", timeRecord)
	}

	if timeRecord.ProjectName != "Project XY" {
		t.Fail()
		t.Logf("The project name was not returned correctly: %#v", timeRecord)
	}

	if len(timeRecord.Tags) == 0 {
		t.Fail()
		t.Logf("The tags were not returned correctly: %#v", timeRecord)
	}

	if timeRecord.Description != "Some stuff" {
		t.Fail()
		t.Logf("The description was not returned correctly: %#v", timeRecord)
	}
}

func Test_GetTimeRecords_NoRowsGiven_EmptyListReturned(t *testing.T) {
	// arrange
	dateFormatter := date.NewISO8601Formatter()
	csvMapper := &CSVTimeRecordMapper{
		dateFormatter: dateFormatter,
		columnNames:   []string{"Start", "Stop", "Workspace Name", "Project Name", "Client Name", "Tag(s)", "Description"},
		tagsSeparator: ",",
	}

	var rows [][]string

	// act
	records, _ := csvMapper.GetTimeRecords(rows)

	// assert
	if len(records) > 0 {
		t.Fail()
		t.Logf("GetTimeRecords should not return time records if the given rows were empty")
	}

}

func Test_GetTimeRecords_OnlyHeadlineGiven_EmptyListReturned(t *testing.T) {
	// arrange
	dateFormatter := date.NewISO8601Formatter()
	csvMapper := &CSVTimeRecordMapper{
		dateFormatter: dateFormatter,
		columnNames:   []string{"Start", "Stop", "Workspace Name", "Project Name", "Client Name", "Tag(s)", "Description"},
		tagsSeparator: ",",
	}

	rows := [][]string{
		[]string{"Start", "Stop", "Workspace Name", "Project Name", "Client Name", "Tag(s)", "Description"},
	}

	// act
	records, _ := csvMapper.GetTimeRecords(rows)

	// assert
	if len(records) > 0 {
		t.Fail()
		t.Logf("GetTimeRecords should not return time records if the given rows only contained the headline")
	}

}

func Test_GetTimeRecords_HeadlineGiven_RowCannotBeParsed_ErrorIsReturned(t *testing.T) {
	// arrange
	dateFormatter := date.NewISO8601Formatter()
	csvMapper := &CSVTimeRecordMapper{
		dateFormatter: dateFormatter,
		columnNames:   []string{"Start", "Stop", "Workspace Name", "Project Name", "Client Name", "Tag(s)", "Description"},
		tagsSeparator: ",",
	}

	rows := [][]string{
		[]string{"Start", "Stop", "Workspace Name", "Project Name", "Client Name", "Tag(s)", "Description"},
		[]string{"Invalid Date", "2015-03-26T11:30:00+01:00", "Workspace", "Project XY", "Client X", "", "Some stuff"},
	}

	// act
	_, err := csvMapper.GetTimeRecords(rows)

	// assert
	if err == nil {
		t.Fail()
		t.Logf("GetTimeRecords should have returned an error")
	}

}

func Test_GetTimeRecords_HeadlineGiven_ValidRow_OneRecordIsReturned(t *testing.T) {
	// arrange
	dateFormatter := date.NewISO8601Formatter()
	csvMapper := &CSVTimeRecordMapper{
		dateFormatter: dateFormatter,
		columnNames:   []string{"Start", "Stop", "Workspace Name", "Project Name", "Client Name", "Tag(s)", "Description"},
		tagsSeparator: ",",
	}

	rows := [][]string{
		[]string{"Start", "Stop", "Workspace Name", "Project Name", "Client Name", "Tag(s)", "Description"},
		[]string{"2015-03-26T08:00:00+01:00", "2015-03-26T11:30:00+01:00", "Workspace", "Project XY", "Client X", "", "Some stuff"},
	}

	// act
	records, _ := csvMapper.GetTimeRecords(rows)

	// assert
	if len(records) != 1 {
		t.Fail()
		t.Logf("GetTimeRecords should have returned one time record")
	}

}

func Test_GetTimeRecords_HeadlineGiven_ThreeRows_ThreeRecordsAreReturned(t *testing.T) {
	// arrange
	dateFormatter := date.NewISO8601Formatter()
	csvMapper := &CSVTimeRecordMapper{
		dateFormatter: dateFormatter,
		columnNames:   []string{"Start", "Stop", "Workspace Name", "Project Name", "Client Name", "Tag(s)", "Description"},
		tagsSeparator: ",",
	}

	rows := [][]string{
		[]string{"Start", "Stop", "Workspace Name", "Project Name", "Client Name", "Tag(s)", "Description"},
		[]string{"2015-03-26T08:00:00+01:00", "2015-03-26T11:30:00+01:00", "Workspace", "Project XY", "Client X", "", "Some stuff"},
		[]string{"2015-03-27T08:00:00+01:00", "2015-03-27T11:30:00+01:00", "Workspace", "Project XY", "Client X", "", "Some more stuff"},
		[]string{"2015-03-28T08:00:00+01:00", "2015-03-28T11:30:00+01:00", "Workspace", "Project XY", "Client X", "", "Some more stuff"},
	}

	// act
	records, _ := csvMapper.GetTimeRecords(rows)

	// assert
	if len(records) != 3 {
		t.Fail()
		t.Logf("GetTimeRecords should have returned three time records")
	}

}

func Test_GetColumnNames_GivenColumnNamesAreReturnedWithoutModification(t *testing.T) {
	// arrange
	dateFormatter := date.NewISO8601Formatter()
	csvMapper := &CSVTimeRecordMapper{
		dateFormatter: dateFormatter,
		columnNames:   []string{"Start", "Stop", "Workspace Name", "Project Name", "Client Name", "Tag(s)", "Description"},
		tagsSeparator: ",",
	}

	// act
	columnNames := csvMapper.GetColumnNames()

	// assert
	expected := "Start,Stop,Workspace Name,Project Name,Client Name,Tag(s),Description"
	if strings.Join(columnNames, ",") != expected {
		t.Fail()
		t.Logf("GetColumnNames should have returned %q", expected)
	}

}

func Test_GetRow(t *testing.T) {
	// arrange
	dateFormatter := date.NewISO8601Formatter()
	csvMapper := &CSVTimeRecordMapper{
		dateFormatter: dateFormatter,
		columnNames:   []string{"Start", "Stop", "Workspace Name", "Project Name", "Client Name", "Tag(s)", "Description"},
		tagsSeparator: ",",
	}

	timeRecord := toggl.TimeRecord{
		WorkspaceName: "Workspace",
		ProjectName:   "Project XY",
		ClientName:    "Client X",
		Tags:          []string{"Tag 1", "Tag 2", "XYZ"},
		Description:   "Some stuff",
	}

	// act
	row := csvMapper.GetRow(timeRecord)

	// assert
	expected := "0001-01-01T00:00:00+00:00|0001-01-01T00:00:00+00:00|Workspace|Project XY|Client X|Tag 1,Tag 2,XYZ|Some stuff"
	if strings.Join(row, "|") != expected {
		t.Fail()
		t.Logf("GetRow returned an invalid value. Expected: %q, Actual: %q", expected, strings.Join(row, "|"))
	}
}
