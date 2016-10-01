package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"
)

func Test_togglCli_Execute_ExportActionIsGiven_NoTokenArgumemtGiven_ErrorIsPrinted(t *testing.T) {
	// arrange
	inputString := ``
	inputReader := strings.NewReader(inputString)

	var outputBuffer bytes.Buffer
	outputWriter := bufio.NewWriter(&outputBuffer)

	var errorBuffer bytes.Buffer
	errorWriter := bufio.NewWriter(&errorBuffer)

	arguments := []string{
		"export",
	}

	cli := togglCli{
		exporterFactory: func(string) CSVExporter {
			return getMockCSVExporter(nil)
		},
	}

	// act
	cli.Execute(inputReader, outputWriter, errorWriter, arguments)

	// assert
	outputWriter.Flush()
	errorWriter.Flush()

	if !strings.Contains(errorBuffer.String(), `required argument 'token' not provided`) {
		t.Fail()
		t.Logf("togglCli_Execute should print an error text if the required token parameter is not given but wrote this instead: %s", errorBuffer.String())
	}
}

func Test_togglCli_Execute_ExportActionIsGiven_TokenArgumemtGiven_NoStartDateProvided_ErrorIsPrinted(t *testing.T) {
	// arrange
	inputString := ``
	inputReader := strings.NewReader(inputString)

	var outputBuffer bytes.Buffer
	outputWriter := bufio.NewWriter(&outputBuffer)

	var errorBuffer bytes.Buffer
	errorWriter := bufio.NewWriter(&errorBuffer)

	arguments := []string{
		"export",
		"1971800d4d82861d8f2c1651fea4d212",
	}

	cli := togglCli{
		exporterFactory: func(string) CSVExporter {
			return getMockCSVExporter(nil)
		},
	}

	// act
	cli.Execute(inputReader, outputWriter, errorWriter, arguments)

	// assert
	outputWriter.Flush()
	errorWriter.Flush()

	if !strings.Contains(errorBuffer.String(), `required argument 'startdate' not provided`) {
		t.Fail()
		t.Logf("togglCli_Execute should print an error text if the required token parameter is not given but wrote this instead: %s", errorBuffer.String())
	}
}

func Test_togglCli_Execute_ExportActionIsGiven_TokenArgumemtGiven_StartDateArgumentInvalid_ErrorIsReturned(t *testing.T) {
	// arrange
	inputString := ``
	inputReader := strings.NewReader(inputString)

	var outputBuffer bytes.Buffer
	outputWriter := bufio.NewWriter(&outputBuffer)

	var errorBuffer bytes.Buffer
	errorWriter := bufio.NewWriter(&errorBuffer)

	arguments := []string{
		"export",
		"1971800d4d82861d8f2c1651fea4d212",
		"2013  03.   10", // invalid start date
	}

	cli := togglCli{
		exporterFactory: func(string) CSVExporter {
			return getMockCSVExporter(nil)
		},
	}

	// act
	cli.Execute(inputReader, outputWriter, errorWriter, arguments)

	// assert
	outputWriter.Flush()
	errorWriter.Flush()

	if !strings.Contains(errorBuffer.String(), "Failed to parse the given start date") {
		t.Fail()
		t.Logf("togglCli_Execute should print an error if the given start date is invalid: %s", errorBuffer.String())
	}
}

func Test_togglCli_Execute_ExportActionIsGiven_TokenArgumemtGiven_StartDateValid_InvalidEndDate_ErrorIsReturned(t *testing.T) {
	// arrange
	inputString := ``
	inputReader := strings.NewReader(inputString)

	var outputBuffer bytes.Buffer
	outputWriter := bufio.NewWriter(&outputBuffer)

	var errorBuffer bytes.Buffer
	errorWriter := bufio.NewWriter(&errorBuffer)

	arguments := []string{
		"export",
		"1971800d4d82861d8f2c1651fea4d212",
		"2013-03-10",
		"2013 09.  -30", // invalid end date
	}

	cli := togglCli{
		exporterFactory: func(string) CSVExporter {
			return getMockCSVExporter(nil)
		},
	}

	// act
	cli.Execute(inputReader, outputWriter, errorWriter, arguments)

	// assert
	outputWriter.Flush()
	errorWriter.Flush()

	if !strings.Contains(errorBuffer.String(), "Failed to parse the given end date") {
		t.Fail()
		t.Logf("togglCli_Execute should print an error if the given end date is invalid: %s", errorBuffer.String())
	}
}

func Test_togglCli_Execute_ExportActionIsGiven_TokenArgumemtGiven_StartDateValid_EndDateValid_ExportSucceeds_NoErrorIsReturned(t *testing.T) {
	// arrange
	inputString := ``
	inputReader := strings.NewReader(inputString)

	var outputBuffer bytes.Buffer
	outputWriter := bufio.NewWriter(&outputBuffer)

	var errorBuffer bytes.Buffer
	errorWriter := bufio.NewWriter(&errorBuffer)

	arguments := []string{
		"export",
		"1971800d4d82861d8f2c1651fea4d212",
		"2013-03-10",
		"2013-09-30",
	}

	cli := togglCli{
		exporterFactory: func(string) CSVExporter {
			return getMockCSVExporter(nil)
		},
	}

	// act
	cli.Execute(inputReader, outputWriter, errorWriter, arguments)

	// assert
	outputWriter.Flush()
	errorWriter.Flush()

	if len(errorBuffer.String()) != 0 {
		t.Fail()
		t.Logf("togglCli_Execute should not print an error if the export succeeded: %s", errorBuffer.String())
	}
}

func Test_togglCli_Execute_ExportActionIsGiven_TokenArgumemtGiven_StartDateValid_NoEndDateGiven_CurrentDateIsUsedAsEndDate(t *testing.T) {
	// arrange
	inputString := ``
	inputReader := strings.NewReader(inputString)

	var outputBuffer bytes.Buffer
	outputWriter := bufio.NewWriter(&outputBuffer)

	var errorBuffer bytes.Buffer
	errorWriter := bufio.NewWriter(&errorBuffer)

	arguments := []string{
		"export",
		"1971800d4d82861d8f2c1651fea4d212",
		"2013-03-10",
	}

	mockCSVExporter := &MockCSVExporter{
		exportFunc: func(startDate, endDate time.Time, writer io.Writer) error {

			// assert
			now := time.Now()
			expectedEndDate := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.UTC)
			if !endDate.Equal(expectedEndDate) {
				t.Fail()
				t.Logf("togglCli_Execute should have used %q as the end date but used %q instead", expectedEndDate, endDate)
			}

			return nil
		},
	}

	cli := togglCli{
		exporterFactory: func(string) CSVExporter {
			return mockCSVExporter
		},
	}

	// act
	cli.Execute(inputReader, outputWriter, errorWriter, arguments)
}

func Test_togglCli_Execute_ExportActionIsGiven_TokenArgumemtGiven_StartDateValid_ExportSucceeds_NoErrorIsReturned(t *testing.T) {
	// arrange
	inputString := ``
	inputReader := strings.NewReader(inputString)

	var outputBuffer bytes.Buffer
	outputWriter := bufio.NewWriter(&outputBuffer)

	var errorBuffer bytes.Buffer
	errorWriter := bufio.NewWriter(&errorBuffer)

	arguments := []string{
		"export",
		"1971800d4d82861d8f2c1651fea4d212",
		"2013-03-10",
	}

	cli := togglCli{
		exporterFactory: func(string) CSVExporter {
			return getMockCSVExporter(nil)
		},
	}

	// act
	cli.Execute(inputReader, outputWriter, errorWriter, arguments)

	// assert
	outputWriter.Flush()
	errorWriter.Flush()

	if len(errorBuffer.String()) != 0 {
		t.Fail()
		t.Logf("togglCli_Execute should not print an error if the export succeeded: %s", errorBuffer.String())
	}
}

func Test_togglCli_Execute_ExportActionIsGiven_TokenArgumemtGiven_StartDateArgumentGiven_ExportFails_ErrorIsReturned(t *testing.T) {
	// arrange
	inputString := ``
	inputReader := strings.NewReader(inputString)

	var outputBuffer bytes.Buffer
	outputWriter := bufio.NewWriter(&outputBuffer)

	var errorBuffer bytes.Buffer
	errorWriter := bufio.NewWriter(&errorBuffer)

	arguments := []string{
		"export",
		"1971800d4d82861d8f2c1651fea4d212",
		"2013-03-10",
	}

	cli := togglCli{
		exporterFactory: func(string) CSVExporter {
			return getMockCSVExporter(fmt.Errorf("Export failed"))
		},
	}

	// act
	cli.Execute(inputReader, outputWriter, errorWriter, arguments)

	// assert
	outputWriter.Flush()
	errorWriter.Flush()

	if !strings.Contains(errorBuffer.String(), "Export failed") {
		t.Fail()
		t.Logf("togglCli_Execute should print an error if the CSV exporter returns one: %s", errorBuffer.String())
	}
}
