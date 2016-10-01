// +build ignore

// Copyright 2016 Andreas Koch. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The make.go script builds the Toggl⥃CSV command-line utility.
//
// $ go run make.go -install
//
// View the README.md for further details.
//
// The output binaries go into the ./bin/ directory
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ProjectName contains the name of the project
const ProjectName = "togglcsv"

// PackagePath contains the package name of the project
const ProjectNamespace = "github.com/andreaskoch/togglcsv"

// GOPATH environment variable name.
const GOPATH_ENVIRONMENT_VARIABLE = "GOPATH"

// GOBIN environment variable name
const GOBIN_ENVIRONMENT_VARIBALE = "GOBIN"

// GOOS environment variable name
const GOOS_ENVIRONMENT_VARIBALE = "GOOS"

// GOARCH environment variable name
const GOARCH_ENVIRONMENT_VARIBALE = "GOARCH"

// GO15VENDOREXPERIMENT environment variable name
const GO15VENDOREXPERIMENT_ENVIRONMENT_VARIBALE = "GO15VENDOREXPERIMENT"

var (

	// packages
	packages = []string{
		"",
		"toggl",
	}

	// command line flags
	verboseFlagIsSet      = flag.Bool("v", false, "Verbose mode")
	installFlagIsSet      = flag.Bool("install", false, fmt.Sprintf("Compiles the %s binaries", ProjectName))
	crossCompileFlagIsSet = flag.Bool("crosscompile", false, fmt.Sprintf("Compile %s binaries for all platforms and architectures", ProjectName))
	testFlagIsSet         = flag.Bool("test", false, "Run the tests")
	coverageFlagIsSet     = flag.Bool("coverage", false, "Run the test and create a code coverage report")

	// The GOPATH for the current project
	goPath = getWorkingDirectory()

	// The GOBIN for the current project
	goBin = filepath.Join(goPath, "bin")

	// A list of all supported compilation targets (e.g. "windows/amd64")
	compilationTargets = []compilationTarget{
		compilationTarget{"darwin", "amd64", []string{}},
		compilationTarget{"linux", "amd64", []string{}},
		compilationTarget{"linux", "arm", []string{}},
		compilationTarget{"linux", "arm", []string{"GOARM=5"}},
		compilationTarget{"linux", "arm", []string{"GOARM=6"}},
		compilationTarget{"linux", "arm", []string{"GOARM=7"}},
		compilationTarget{"windows", "amd64", []string{}},
	}
)

// Compilation Target Definition
type compilationTarget struct {
	OperatingSystem string
	Architecture    string
	OtherVariables  []string
}

func (target *compilationTarget) String() string {
	return fmt.Sprintf("%s/%s", target.OperatingSystem, target.Architecture)
}

func init() {

	executableName := "make.go"

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s provides functions for compiling %s.\n", executableName, ProjectName)
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  go run %s [options]\n", executableName)
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "\n")
		flag.PrintDefaults()
	}
}

func main() {
	log.SetFlags(0)
	flag.Parse()

	if len(os.Args) < 2 {
		flag.Usage()
		return
	}

	if *verboseFlagIsSet {
		fmt.Printf("%s: %s\n", GOPATH_ENVIRONMENT_VARIABLE, goPath)
		fmt.Printf("%s: %s\n", GOBIN_ENVIRONMENT_VARIBALE, goBin)
	}

	if *installFlagIsSet {
		if err := install(); err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}

		return
	}

	if *crossCompileFlagIsSet {
		if err := crossCompile(); err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}

		return
	}

	if *testFlagIsSet {
		if err := test(); err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}

		return
	}

	if *coverageFlagIsSet {
		if err := codecoverage(); err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}

		return
	}

	flag.Usage()
}

// Install all parts of allmark using go install.
func install() error {

	// prepare the environment variables
	environmentVariables := cleanGoEnv()
	environmentVariables = setEnv(environmentVariables, GOBIN_ENVIRONMENT_VARIBALE, goBin)
	environmentVariables = setEnv(environmentVariables, GO15VENDOREXPERIMENT_ENVIRONMENT_VARIBALE, "1")

	return runCommand(os.Stdout, os.Stderr, goPath, environmentVariables, "go", "install", fmt.Sprintf("%s/...", ProjectNamespace))

}

// test runs all tests
func test() error {

	// prepare the environment variables
	environmentVariables := cleanGoEnv()
	environmentVariables = setEnv(environmentVariables, GOBIN_ENVIRONMENT_VARIBALE, goBin)
	environmentVariables = setEnv(environmentVariables, GO15VENDOREXPERIMENT_ENVIRONMENT_VARIBALE, "1")

	testError := runCommand(os.Stdout, os.Stderr, goPath, environmentVariables, "go", "test", fmt.Sprintf("%s/...", ProjectNamespace))
	if testError != nil {
		return testError
	}

	return nil
}

// Cross-compile all parts of allmark for all supported platforms.
func crossCompile() error {

	// prepare the environment variables
	environmentVariables := cleanGoEnv()
	environmentVariables = setEnv(environmentVariables, GOBIN_ENVIRONMENT_VARIBALE, goBin)
	environmentVariables = setEnv(environmentVariables, GO15VENDOREXPERIMENT_ENVIRONMENT_VARIBALE, "1")

	// iterate over all supported compilation targets
	for _, target := range compilationTargets {

		// assemble the target path
		targetFile := filepath.Join(goBin, fmt.Sprintf("%s_%s_%s", ProjectName, target.OperatingSystem, target.Architecture))

		// prepare environment variables for cross-compilation
		crossCompileEnvironemntVariables := environmentVariables
		crossCompileEnvironemntVariables = setEnv(crossCompileEnvironemntVariables, GOOS_ENVIRONMENT_VARIBALE, target.OperatingSystem)
		crossCompileEnvironemntVariables = setEnv(crossCompileEnvironemntVariables, GOARCH_ENVIRONMENT_VARIBALE, target.Architecture)

		// add additional environment variables
		for _, additionalEnvVariable := range target.OtherVariables {
			components := strings.Split(additionalEnvVariable, "=")
			name := components[0]
			value := components[1]

			// append additional environment variables to target file name
			targetFile += fmt.Sprintf("_%s_%s", strings.ToLower(name), strings.ToLower(value))

			crossCompileEnvironemntVariables = setEnv(crossCompileEnvironemntVariables, name, value)
		}

		// build the package for the specified os and arch
		if *verboseFlagIsSet {
			fmt.Printf("Compiling %s for %s\n", ProjectName, target.String())
		}

		// add .exe extension for windows
		if target.OperatingSystem == "windows" {
			targetFile += ".exe"
		}

		err := runCommand(os.Stdout,
			os.Stderr,
			goPath,
			crossCompileEnvironemntVariables,
			"go",
			"build",
			"-o",
			targetFile,
			fmt.Sprintf("%s", ProjectNamespace))

		if err != nil {
			return err
		}

	}

	return nil
}

// codecoverage runs all tests and creates a code coverage report
func codecoverage() error {

	// prepare the environment variables
	environmentVariables := cleanGoEnv()
	environmentVariables = setEnv(environmentVariables, GOBIN_ENVIRONMENT_VARIBALE, goBin)
	environmentVariables = setEnv(environmentVariables, GO15VENDOREXPERIMENT_ENVIRONMENT_VARIBALE, "1")

	for _, packageName := range packages {

		packageTitle := packageName
		if packageTitle == "" {
			packageTitle = "cli"
		}

		packagePath := "./" + packageName

		coverageError := runCommand(os.Stdout, os.Stderr, goPath, environmentVariables, "go", "test", fmt.Sprintf("-coverprofile=coverage-%s.out", packageTitle), packagePath)
		if coverageError != nil {
			return coverageError
		}

		reportError := runCommand(os.Stdout, os.Stderr, goPath, environmentVariables, "go", "tool", "cover", fmt.Sprintf("-html=coverage-%s.out", packageTitle), "-o", fmt.Sprintf("coverage-%s.html", packageTitle))
		if reportError != nil {
			return reportError
		}

	}

	return nil
}

// getWorkingDirectory returns the current working directory path or fails.
func getWorkingDirectory() string {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current directory: %v", err)
	}

	return wd
}

// Execute go in the specified go path with the supplied command arguments.
func runCommand(stdout, stderr io.Writer, workingDirectory string, environmentVariables []string, command string, args ...string) error {

	// Create the command
	cmdName := fmt.Sprintf("%s %s", command, strings.Join(args, " "))
	cmd := exec.Command(command, args...)

	// Set the working directory
	cmd.Dir = workingDirectory

	// set environment variables
	cmd.Env = environmentVariables

	// Capture the output
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if *verboseFlagIsSet {
		log.Printf("Running %s", cmdName)
	}

	// execute the command
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Error running %s: %v", cmdName, err.Error())
	}

	return nil
}

// cleanGoEnv returns a copy of the current environment with GOPATH_ENVIRONMENT_VARIABLE and GOBIN_ENVIRONMENT_VARIBALE removed.
func cleanGoEnv() (clean []string) {
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, GOBIN_ENVIRONMENT_VARIBALE+"=") {
			continue
		}

		clean = append(clean, env)
	}

	return
}

// setEnv sets the given key & value in the provided environment.
// Each value in the env list should be of the form key=value.
func setEnv(env []string, key, value string) []string {
	for i, s := range env {
		if strings.HasPrefix(s, fmt.Sprintf("%s=", key)) {
			env[i] = envPair(key, value)
			return env
		}
	}
	env = append(env, envPair(key, value))
	return env
}

// Create an environment variable of the form key=value.
func envPair(key, value string) string {
	return fmt.Sprintf("%s=%s", key, value)
}
