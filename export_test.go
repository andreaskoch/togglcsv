package main

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/andreaskoch/togglcsv/toggl"
)

type mockCSVTimeRecordMapper struct {
	getTimeRecords func(rows [][]string) ([]toggl.TimeRecord, error)
	getTimeRecord  func(row []string) (toggl.TimeRecord, error)
	columnNames    []string
	getRow         func(timeRecord toggl.TimeRecord) []string
}

func (mapper *mockCSVTimeRecordMapper) GetTimeRecords(rows [][]string) ([]toggl.TimeRecord, error) {
	return mapper.getTimeRecords(rows)
}

func (mapper *mockCSVTimeRecordMapper) GetTimeRecord(row []string) (toggl.TimeRecord, error) {
	return mapper.GetTimeRecord(row)
}

func (mapper *mockCSVTimeRecordMapper) GetColumnNames() []string {
	return mapper.columnNames
}

func (mapper *mockCSVTimeRecordMapper) GetRow(timeRecord toggl.TimeRecord) []string {
	return mapper.getRow(timeRecord)
}

type mockTimeRecordRepository struct {
	createTimeRecord func(timeRecord toggl.TimeRecord) error
	getTimeRecords   func(start, stop time.Time) ([]toggl.TimeRecord, error)
}

func (repository *mockTimeRecordRepository) CreateTimeRecord(timeRecord toggl.TimeRecord) error {
	return repository.createTimeRecord(timeRecord)
}

func (repository *mockTimeRecordRepository) GetTimeRecords(start, stop time.Time) ([]toggl.TimeRecord, error) {
	return repository.getTimeRecords(start, stop)
}

func Test_Export_NoTimeRecordsReturned_OnlyTheCSVHeaderIsWritten(t *testing.T) {
	// arrange
	timeRecordMapper := &mockCSVTimeRecordMapper{
		columnNames: []string{"Col 1", "Col 2", "Col 3"},
	}

	timeRecords := []toggl.TimeRecord{}
	timeRecordRepository := &mockTimeRecordRepository{
		getTimeRecords: func(start, stop time.Time) ([]toggl.TimeRecord, error) {
			return timeRecords, nil
		},
	}

	exporter := TogglCSVExporter{
		csvMapper:            timeRecordMapper,
		timeRecordRepository: timeRecordRepository,
	}

	startDate := time.Date(2016, 5, 3, 0, 0, 1, 0, time.UTC)
	endDate := time.Date(2016, 8, 3, 0, 0, 1, 0, time.UTC)
	var outputBuffer bytes.Buffer
	writer := bufio.NewWriter(&outputBuffer)

	// act
	exporter.Export(startDate, endDate, writer)
	result := outputBuffer.String()

	// assert
	expected := "Col 1,Col 2,Col 3\n"
	if result != expected {
		t.Fail()
		t.Logf("The output of the Export function should have been '%s' but was '%s' instead", expected, result)
	}
}

func Test_Export_TimeRecordRepositoryReturnsError_ErrorIsReturned(t *testing.T) {
	// arrange
	timeRecordMapper := &mockCSVTimeRecordMapper{
		columnNames: []string{"Col 1", "Col 2", "Col 3"},
		getRow: func(timeRecord toggl.TimeRecord) []string {
			return []string{"row"}
		},
	}

	timeRecordRepository := &mockTimeRecordRepository{
		getTimeRecords: func(start, stop time.Time) ([]toggl.TimeRecord, error) {
			return nil, fmt.Errorf("Some error")
		},
	}

	exporter := TogglCSVExporter{
		csvMapper:            timeRecordMapper,
		timeRecordRepository: timeRecordRepository,
	}

	startDate := time.Date(2016, 5, 3, 0, 0, 1, 0, time.UTC)
	endDate := time.Date(2016, 8, 3, 0, 0, 1, 0, time.UTC)
	var outputBuffer bytes.Buffer
	writer := bufio.NewWriter(&outputBuffer)

	// act
	err := exporter.Export(startDate, endDate, writer)

	// assert
	if err == nil {
		t.Fail()
		t.Logf("Export should have returned an error but didn't")
	}
}

func Test_Export_OneTimeRecordReturned_HeaderAndRowIsWritten(t *testing.T) {
	// arrange
	timeRecordMapper := &mockCSVTimeRecordMapper{
		columnNames: []string{"Col 1", "Col 2", "Col 3"},
		getRow: func(timeRecord toggl.TimeRecord) []string {
			return []string{"row"}
		},
	}

	timeRecords := []toggl.TimeRecord{
		toggl.TimeRecord{},
	}

	timeRecordRepository := &mockTimeRecordRepository{
		getTimeRecords: func(start, stop time.Time) ([]toggl.TimeRecord, error) {
			return timeRecords, nil
		},
	}

	exporter := TogglCSVExporter{
		csvMapper:            timeRecordMapper,
		timeRecordRepository: timeRecordRepository,
	}

	startDate := time.Date(2016, 5, 3, 0, 0, 1, 0, time.UTC)
	endDate := time.Date(2016, 8, 3, 0, 0, 1, 0, time.UTC)
	var outputBuffer bytes.Buffer
	writer := bufio.NewWriter(&outputBuffer)

	// act
	exporter.Export(startDate, endDate, writer)
	result := outputBuffer.String()

	// assert
	expected := "Col 1,Col 2,Col 3\n" + "row\n"
	if result != expected {
		t.Fail()
		t.Logf("The output of the Export function should have been '%s' but was '%s' instead", expected, result)
	}
}

func Test_Export_MultipleTimeRecordReturned_MultipleRowsAreWritten(t *testing.T) {
	// arrange
	rowIndex := 0
	timeRecordMapper := &mockCSVTimeRecordMapper{
		columnNames: []string{"Col 1", "Col 2", "Col 3"},
		getRow: func(timeRecord toggl.TimeRecord) []string {
			rowIndex++
			return []string{fmt.Sprintf("row %d", rowIndex)}
		},
	}

	timeRecords := []toggl.TimeRecord{
		toggl.TimeRecord{},
		toggl.TimeRecord{},
		toggl.TimeRecord{},
	}

	timeRecordRepository := &mockTimeRecordRepository{
		getTimeRecords: func(start, stop time.Time) ([]toggl.TimeRecord, error) {
			return timeRecords, nil
		},
	}

	exporter := TogglCSVExporter{
		csvMapper:            timeRecordMapper,
		timeRecordRepository: timeRecordRepository,
	}

	startDate := time.Date(2016, 5, 3, 0, 0, 1, 0, time.UTC)
	endDate := time.Date(2016, 8, 3, 0, 0, 1, 0, time.UTC)
	var outputBuffer bytes.Buffer
	writer := bufio.NewWriter(&outputBuffer)

	// act
	exporter.Export(startDate, endDate, writer)
	result := outputBuffer.String()

	// assert
	expected := "Col 1,Col 2,Col 3\n" + "row 1\n" + "row 2\n" + "row 3\n"
	if result != expected {
		t.Fail()
		t.Logf("The output of the Export function should have been '%s' but was '%s' instead", expected, result)
	}
}
