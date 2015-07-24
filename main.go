package main

import (
	"fmt"
	"os"
)

func main() {

	var project *TemplatedProject
	var settings RunSettings
	var parameters []string
	var err error

	settings, err = FindRunSettings()
	if err != nil {
		exitWith("Unable to parse run arguments: %s\n", err, 1)
		return
	}

	project, err = NewTemplatedProject(settings.skeletonPath)
	if err != nil {
		exitWith("Unable to read templated project: %s\n", err, 1)
		return
	}

	// inspect only?
	if settings.inspectionRun {

		parameters, err = project.FindParameters()

		if err != nil {
			exitWith("Unable to inspect skeleton: %s\n", err, 1)
			return
		}

		for _, parameter := range parameters {
			fmt.Println(parameter)
		}
		return
	}

	// generate a project
	err = project.GenerateAt(settings.targetPath, settings.parameters)
	if err != nil {
		exitWith("Unable to generate project: %s\n", err, 1)
		return
	}
}

func exitWith(message string, err error, code int) {

	errorMsg := fmt.Sprintf(message, err)
	fmt.Fprintf(os.Stderr, errorMsg)
	os.Exit(code)
}
