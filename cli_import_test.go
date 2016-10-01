package main

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func Test_togglCli_Execute_ImportActionIsGiven_NoTokenArgumemtGiven_ErrorIsPrinted(t *testing.T) {
	// arrange
	inputString := ``
	inputReader := strings.NewReader(inputString)

	var outputBuffer bytes.Buffer
	outputWriter := bufio.NewWriter(&outputBuffer)

	var errorBuffer bytes.Buffer
	errorWriter := bufio.NewWriter(&errorBuffer)

	arguments := []string{
		"import",
	}

	cli := togglCli{
		importerFactory: func(string) CSVImporter {
			return getMockCSVImporter(nil)
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

func Test_togglCli_Execute_ImportActionIsGiven_TokenArgumemtGiven_InvalidCSVInputGiven_ErrorIsPrinted(t *testing.T) {
	// arrange
	inputString := `ddas; -das ... -> dasdas;,,:`
	inputReader := strings.NewReader(inputString)

	var outputBuffer bytes.Buffer
	outputWriter := bufio.NewWriter(&outputBuffer)

	var errorBuffer bytes.Buffer
	errorWriter := bufio.NewWriter(&errorBuffer)

	arguments := []string{
		"import",
		"1971800d4d82861d8f2c1651fea4d212",
	}

	cli := togglCli{
		importerFactory: func(string) CSVImporter {
			return getMockCSVImporter(fmt.Errorf("invalid csv"))
		},
	}

	// act
	cli.Execute(inputReader, outputWriter, errorWriter, arguments)

	// assert
	outputWriter.Flush()
	errorWriter.Flush()

	if !strings.Contains(errorBuffer.String(), "invalid csv") {
		t.Fail()
		t.Logf("togglCli_Execute should print an error text if the given input is not valid CSV but wrote this instead: %s", errorBuffer.String())
	}
}

func Test_togglCli_Execute_ImportActionIsGiven_TokenArgumemtGiven_NoInputGiven_NoErrorIsPrinted(t *testing.T) {
	// arrange
	inputString := ``
	inputReader := strings.NewReader(inputString)

	var outputBuffer bytes.Buffer
	outputWriter := bufio.NewWriter(&outputBuffer)

	var errorBuffer bytes.Buffer
	errorWriter := bufio.NewWriter(&errorBuffer)

	arguments := []string{
		"import",
		"1971800d4d82861d8f2c1651fea4d212",
	}

	cli := togglCli{
		importerFactory: func(string) CSVImporter {
			return getMockCSVImporter(nil)
		},
	}

	// act
	cli.Execute(inputReader, outputWriter, errorWriter, arguments)

	// assert
	outputWriter.Flush()
	errorWriter.Flush()

	if len(errorBuffer.String()) != 0 {
		t.Fail()
		t.Logf("togglCli_Execute should not print any output if no input was given but returned this instead: %s", errorBuffer.String())
	}
}

func Test_togglCli_Execute_ImportActionIsGiven_TokenArgumemtGiven_ValidCSV_ImportSucceeds_NoErrorIsReturned(t *testing.T) {
	// arrange
	inputString := `Start,Stop,Workspace Name,Project Name,Tag(s),Description
	2016-08-12T07:54:47+01:00,2016-08-12T08:19:02+01:00,My Workspace,Project A,"Meetings, Sprint",Retrospective
` // doesn't really matter if this is valid or not
	inputReader := strings.NewReader(inputString)

	var outputBuffer bytes.Buffer
	outputWriter := bufio.NewWriter(&outputBuffer)

	var errorBuffer bytes.Buffer
	errorWriter := bufio.NewWriter(&errorBuffer)

	arguments := []string{
		"import",
		"1971800d4d82861d8f2c1651fea4d212",
	}

	cli := togglCli{
		importerFactory: func(string) CSVImporter {
			return getMockCSVImporter(nil)
		},
	}

	// act
	cli.Execute(inputReader, outputWriter, errorWriter, arguments)

	// assert
	outputWriter.Flush()
	errorWriter.Flush()

	if len(errorBuffer.String()) != 0 {
		t.Fail()
		t.Logf("togglCli_Execute should not print an error if the import succeeded: %s", errorBuffer.String())
	}
}

func Test_togglCli_Execute_ImportActionIsGiven_TokenArgumemtGiven_ImportFails_ErrorIsReturned(t *testing.T) {
	// arrange
	inputString := ``
	inputReader := strings.NewReader(inputString)

	var outputBuffer bytes.Buffer
	outputWriter := bufio.NewWriter(&outputBuffer)

	var errorBuffer bytes.Buffer
	errorWriter := bufio.NewWriter(&errorBuffer)

	arguments := []string{
		"import",
		"1971800d4d82861d8f2c1651fea4d212",
	}

	cli := togglCli{
		importerFactory: func(string) CSVImporter {
			return getMockCSVImporter(fmt.Errorf("Import failed"))
		},
	}

	// act
	cli.Execute(inputReader, outputWriter, errorWriter, arguments)

	// assert
	outputWriter.Flush()
	errorWriter.Flush()

	if !strings.Contains(errorBuffer.String(), "Import failed") {
		t.Fail()
		t.Logf("togglCli_Execute should print an error if the CSV importer returns one: %s", errorBuffer.String())
	}
}
