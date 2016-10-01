package main

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
)

func Test_getCSVExporter_IntegrationTest_ResultIsNotNull(t *testing.T) {
	// arrange
	apiToken := "dkasjdlkjsadkljas3123j12kl"

	// act
	exporter := getCSVExporter(apiToken)

	// assert
	if exporter == nil {
		t.Fail()
		t.Logf("getCSVExporter should not have returned nil")
	}
}

func Test_getCSVImporter_IntegrationTest_ResultIsNotNull(t *testing.T) {
	// arrange
	apiToken := "dkasjdlkjsadkljas3123j12kl"

	// act
	exporter := getCSVImporter(apiToken)

	// assert
	if exporter == nil {
		t.Fail()
		t.Logf("getCSVImporter should not have returned nil")
	}
}

func Test_IntegrationTest_main_HelpOrInvalidArguments_HelpTextIsPrintedToStderr(t *testing.T) {
	// arrange
	argumentInputs := [][]string{
		[]string{},
		[]string{"help"},
		[]string{"--help"},
		[]string{"help export"},
		[]string{"help import"},
		[]string{"version"},
		[]string{"some-invalid-action-name"},
	}

	for _, arguments := range argumentInputs {

		inputString := ``
		inputReader := strings.NewReader(inputString)

		var outputBuffer bytes.Buffer
		outputWriter := bufio.NewWriter(&outputBuffer)

		var errorBuffer bytes.Buffer
		errorWriter := bufio.NewWriter(&errorBuffer)

		// override globals
		out = outputWriter
		err = errorWriter
		in = inputReader
		args = arguments

		// act
		main()

		outputWriter.Flush()
		errorWriter.Flush()

		// assert
		if len(errorBuffer.String()) == 0 {
			t.Fail()
			t.Logf("main(%q) should have printed help information to the stderr but didn't", strings.Join(arguments, ", "))
		}
	}
}
