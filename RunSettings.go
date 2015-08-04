package main

import (
	"errors"
	"flag"
	"fmt"
	"strings"
	"path/filepath"
)

/*
  Represents arguments given by the user to this program.
*/
type RunSettings struct {

	parameters map[string]string

	skeletonPath string
	targetPath   string

	inspectionRun bool
	addRegistry bool
	showRegistry bool
}

/*
  Parses runtime flags and positional arguments, returning the result.
*/
func FindRunSettings() (RunSettings, error) {

	var ret RunSettings
	var parameterGroup string
	var err error

	flag.StringVar(&parameterGroup, "p", "", "Semicolon-separated list of parameters in k=v form.")
	flag.BoolVar(&ret.inspectionRun, "i", false, "Whether or not to show a list of available parameters for the skeleton")
	flag.BoolVar(&ret.addRegistry, "a", false, "Whether or not to register the template at the given path")
	flag.BoolVar(&ret.showRegistry, "l", false, "Whether or not to show a list of all available registered templates")
	flag.Parse()

	ret.skeletonPath = flag.Arg(0)
	ret.targetPath = flag.Arg(1)

	// if we're not just showing the registry, and not skeleton path is specified...
	if !ret.showRegistry && ret.skeletonPath == "" {
		errorMsg := fmt.Sprintf("Skeleton project path not specified")
		return ret, errors.New(errorMsg)
	}

	// if we're actually generating a project, and no target path is specified...
	if !ret.showRegistry &&
		!ret.inspectionRun &&
		!ret.addRegistry &&
		ret.targetPath == "" {

		errorMsg := fmt.Sprintf("Target output path not specified")
		return ret, errors.New(errorMsg)
	}

	// make parameters, set default project.name
	ret.parameters = make(map[string]string)
	ret.parameters["project.name"] = filepath.Base(ret.targetPath)

	err = parseParametersTo(parameterGroup, ret.parameters)
	if err != nil {
		return ret, err
	}

	return ret, nil
}

/*
  Given a sequence of k=v strings, this parses them out into a map.
*/
func parseParametersTo(parameterGroup string, destination map[string]string) (error) {

	var groups, pair []string

	parameterGroup = strings.Trim(parameterGroup, " ")

	if len(parameterGroup) == 0 {
		return nil
	}

	groups = strings.Split(parameterGroup, ";")
	if len(groups) == 0 {
		groups = []string{
			parameterGroup,
		}
	}

	for _, group := range groups {

		pair = strings.Split(group, "=")

		if len(pair) != 2 {
			errorMsg := fmt.Sprintf("Unable to parse parameters, expected exactly one '=' per semicolon-separated set")
			return errors.New(errorMsg)
		}

		destination[pair[0]] = pair[1]
	}

	return nil
}
