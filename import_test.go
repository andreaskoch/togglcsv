package main

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/andreaskoch/togglcsv/toggl"
)

func Test_Import_NoInput_NoErrorIsReturned(t *testing.T) {
	// arrange
	timeRecords := []toggl.TimeRecord{}

	timeRecordMapper := &mockCSVTimeRecordMapper{
		columnNames: []string{"Start", "Stop", "..."},
		getTimeRecords: func(rows [][]string) ([]toggl.TimeRecord, error) {
			return timeRecords, nil
		},
	}

	timeRecordRepository := &mockTimeRecordRepository{}

	importer := TogglCSVImporter{
		csvMapper:            timeRecordMapper,
		timeRecordRepository: timeRecordRepository,
	}

	input := ""
	inputReader := strings.NewReader(input)

	// act
	err := importer.Import(inputReader)

	// assert
	if err != nil {
		t.Fail()
		t.Logf("Import should not return an error if no input was given but this was returned: %s", err)
	}
}

func Test_Import_InvalidCSV_ErrorIsReturned(t *testing.T) {
	// arrange
	timeRecords := []toggl.TimeRecord{}

	timeRecordMapper := &mockCSVTimeRecordMapper{
		columnNames: []string{"Start", "Stop", "..."},
		getTimeRecords: func(rows [][]string) ([]toggl.TimeRecord, error) {
			return timeRecords, nil
		},
	}

	timeRecordRepository := &mockTimeRecordRepository{}

	importer := TogglCSVImporter{
		csvMapper:            timeRecordMapper,
		timeRecordRepository: timeRecordRepository,
	}

	input := `"Start","Stop,`
	inputReader := strings.NewReader(input)

	// act
	err := importer.Import(inputReader)

	// assert
	if err == nil {
		t.Fail()
		t.Logf("Import should return an error if the given CSV is invalid but didn't")
	}
}

func Test_Import_TimeRecordMapperReturnsError_ErrorIsReturned(t *testing.T) {
	// arrange
	timeRecordMapper := &mockCSVTimeRecordMapper{
		columnNames: []string{"Start", "Stop", "..."},
		getTimeRecords: func(rows [][]string) ([]toggl.TimeRecord, error) {
			return nil, fmt.Errorf("Some error")
		},
	}

	timeRecordRepository := &mockTimeRecordRepository{}

	importer := TogglCSVImporter{
		csvMapper:            timeRecordMapper,
		timeRecordRepository: timeRecordRepository,
	}

	input := ``
	inputReader := strings.NewReader(input)

	// act
	err := importer.Import(inputReader)

	// assert
	if err == nil {
		t.Fail()
		t.Logf("Import should return an error if the time record mapper returns an error")
	}
}

func Test_Import_CreateTimeRecordReturnsError_ErrorIsReturned(t *testing.T) {
	// arrange
	timeRecords := []toggl.TimeRecord{
		toggl.TimeRecord{},
	}

	timeRecordMapper := &mockCSVTimeRecordMapper{
		columnNames: []string{"Start", "Stop", "..."},
		getTimeRecords: func(rows [][]string) ([]toggl.TimeRecord, error) {
			return timeRecords, nil
		},
	}

	timeRecordRepository := &mockTimeRecordRepository{
		createTimeRecord: func(timeRecord toggl.TimeRecord) error {
			return fmt.Errorf("Some error")
		},
	}

	importer := TogglCSVImporter{
		csvMapper:            timeRecordMapper,
		timeRecordRepository: timeRecordRepository,
	}

	input := ``
	inputReader := strings.NewReader(input)

	// act
	err := importer.Import(inputReader)

	// assert
	if err == nil {
		t.Fail()
		t.Logf("Import should return an error if the create function returns an error")
	}
}

func Test_Import_CreateTimeRecordNoError_NoErrorIsReturned(t *testing.T) {
	// arrange
	timeRecords := []toggl.TimeRecord{
		toggl.TimeRecord{},
	}

	timeRecordMapper := &mockCSVTimeRecordMapper{
		columnNames: []string{"Start", "Stop", "..."},
		getTimeRecords: func(rows [][]string) ([]toggl.TimeRecord, error) {
			return timeRecords, nil
		},
	}

	timeRecordRepository := &mockTimeRecordRepository{
		createTimeRecord: func(timeRecord toggl.TimeRecord) error {
			return nil
		},
	}

	importer := TogglCSVImporter{
		csvMapper:            timeRecordMapper,
		timeRecordRepository: timeRecordRepository,
	}

	input := ``
	inputReader := strings.NewReader(input)

	// act
	err := importer.Import(inputReader)

	// assert
	if err != nil {
		t.Fail()
		t.Logf("Import should no return an error if everthing went well")
	}
}

func Test_Import_CreateTimeRecordNoError_ProgressBarIsWritten(t *testing.T) {
	// arrange
	timeRecords := []toggl.TimeRecord{
		toggl.TimeRecord{},
	}

	timeRecordMapper := &mockCSVTimeRecordMapper{
		columnNames: []string{"Start", "Stop", "..."},
		getTimeRecords: func(rows [][]string) ([]toggl.TimeRecord, error) {
			return timeRecords, nil
		},
	}

	var outputBuffer bytes.Buffer
	output := bufio.NewWriter(&outputBuffer)

	timeRecordRepository := &mockTimeRecordRepository{
		createTimeRecord: func(timeRecord toggl.TimeRecord) error {
			return nil
		},
	}

	importer := TogglCSVImporter{
		csvMapper:            timeRecordMapper,
		timeRecordRepository: timeRecordRepository,
		output:               output,
	}

	input := ``
	inputReader := strings.NewReader(input)

	// act
	importer.Import(inputReader)

	// assert
	output.Flush()
	if len(outputBuffer.String()) == 0 {
		t.Fail()
		t.Logf("Import should write a progress bar to the output but didn't")
	}
}
