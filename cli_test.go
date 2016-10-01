package main

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"testing"
	"time"
)

type MockCSVImporter struct {
	importFunc func(input io.Reader) error
}

func (importer *MockCSVImporter) Import(input io.Reader) error {
	return importer.importFunc(input)
}

type MockCSVExporter struct {
	exportFunc func(startDate, endDate time.Time, writer io.Writer) error
}

func (exporter *MockCSVExporter) Export(startDate, endDate time.Time, writer io.Writer) error {
	return exporter.exportFunc(startDate, endDate, writer)
}

func getMockCSVImporter(returnVal error) CSVImporter {
	return &MockCSVImporter{
		importFunc: func(io.Reader) error {
			return returnVal
		},
	}
}

func getMockCSVExporter(returnVal error) CSVExporter {
	return &MockCSVExporter{
		exportFunc: func(startDate, endDate time.Time, writer io.Writer) error {
			return returnVal
		},
	}
}

func Test_togglCli_Execute_NoArguments_HelpIsPrinted(t *testing.T) {
	// arrange
	inputString := ``
	inputReader := strings.NewReader(inputString)

	var outputBuffer bytes.Buffer
	outputWriter := bufio.NewWriter(&outputBuffer)

	var errorBuffer bytes.Buffer
	errorWriter := bufio.NewWriter(&errorBuffer)

	cli := togglCli{
		importerFactory: func(string) CSVImporter {
			return getMockCSVImporter(nil)
		},
		exporterFactory: func(string) CSVExporter {
			return getMockCSVExporter(nil)
		},
	}

	arguments := []string{}

	// act
	cli.Execute(inputReader, outputWriter, errorWriter, arguments)

	// assert
	outputWriter.Flush()
	errorWriter.Flush()

	if !strings.Contains(errorBuffer.String(), "Toggl⥃CSV is an csv-based import/export utility for Toggl time tracking data") {
		t.Fail()
		t.Logf("togglCli_Execute should print help text if no arguments were given but wrote this instead: %s", errorBuffer.String())
	}
}

func Test_togglCli_Execute_HelpActionIsGiven_HelpIsPrinted(t *testing.T) {
	// arrange
	inputString := ``
	inputReader := strings.NewReader(inputString)

	var outputBuffer bytes.Buffer
	outputWriter := bufio.NewWriter(&outputBuffer)

	var errorBuffer bytes.Buffer
	errorWriter := bufio.NewWriter(&errorBuffer)

	arguments := []string{
		"help",
	}

	cli := togglCli{
		importerFactory: func(string) CSVImporter {
			return getMockCSVImporter(nil)
		},
		exporterFactory: func(string) CSVExporter {
			return getMockCSVExporter(nil)
		},
	}

	// act
	cli.Execute(inputReader, outputWriter, errorWriter, arguments)

	// assert
	outputWriter.Flush()
	errorWriter.Flush()

	if !strings.Contains(errorBuffer.String(), "Toggl⥃CSV is an csv-based import/export utility for Toggl time tracking data") {
		t.Fail()
		t.Logf("togglCli_Execute should print help text if no arguments were given but wrote this instead: %s", errorBuffer.String())
	}
}

func Test_togglCli_Execute_IvalidActionIsGiven_ErrorIsPrinted(t *testing.T) {
	// arrange
	inputString := ``
	inputReader := strings.NewReader(inputString)

	var outputBuffer bytes.Buffer
	outputWriter := bufio.NewWriter(&outputBuffer)

	var errorBuffer bytes.Buffer
	errorWriter := bufio.NewWriter(&errorBuffer)

	arguments := []string{
		"some-invalid-action",
	}

	cli := togglCli{
		importerFactory: func(string) CSVImporter {
			return getMockCSVImporter(nil)
		},
		exporterFactory: func(string) CSVExporter {
			return getMockCSVExporter(nil)
		},
	}

	// act
	cli.Execute(inputReader, outputWriter, errorWriter, arguments)

	// assert
	outputWriter.Flush()
	errorWriter.Flush()

	if !strings.Contains(errorBuffer.String(), `expected command but got "some-invalid-action"`) {
		t.Fail()
		t.Logf("togglCli_Execute should print an error text if the given action is not known but wrote this instead: %s", errorBuffer.String())
	}
}
