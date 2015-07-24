package main

import (
	"errors"
	"flag"
	"fmt"
	"strings"
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
	flag.Parse()

	ret.parameters, err = parseParameters(parameterGroup)
	if err != nil {
		return ret, err
	}

	ret.skeletonPath = flag.Arg(0)
	ret.targetPath = flag.Arg(1)

	if ret.skeletonPath == "" {
		errorMsg := fmt.Sprintf("Skeleton project path not specified")
		return ret, errors.New(errorMsg)
	}

	if !ret.inspectionRun && !ret.addRegistry && ret.targetPath == "" {
		errorMsg := fmt.Sprintf("Target output path not specified")
		return ret, errors.New(errorMsg)
	}

	return ret, nil
}

/*
  Given a sequence of k=v strings, this parses them out into a map.
*/
func parseParameters(parameterGroup string) (map[string]string, error) {

	var groups, pair []string
	var ret map[string]string

	parameterGroup = strings.Trim(parameterGroup, " ")
	ret = make(map[string]string, len(groups))

	if len(parameterGroup) == 0 {
		return ret, nil
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
			return ret, errors.New(errorMsg)
		}

		ret[pair[0]] = pair[1]
	}

	return ret, nil
}
