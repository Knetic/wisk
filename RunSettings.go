package main

import (
	"errors"
	"flag"
	"fmt"
	"path/filepath"
	"time"
	"strings"
)

/*
  Represents arguments given by the user to this program.
*/
type RunSettings struct {
	parameters map[string][]string

	skeletonPath string
	targetPath   string

	basicAuthUser string
	basicAuthPass string

	inspectionRun bool
	addRegistry   bool
	showRegistry  bool
	forceGenerate GenerateMode

	flagList	  bool
}

type GenerateMode uint8

const (
	GENERATE_NONE GenerateMode = iota
	GENERATE_OVERWRITE
	GENERATE_DELETE
)

var FLAGS = []string{
	"-p",
	"-i",
	"-a",
	"-l",
	"-f",
	"-w",
	"-bu",
	"-bp",
	"-flags",
}

/*
  Parses runtime flags and positional arguments, returning the result.
*/
func FindRunSettings() (RunSettings, error) {

	var ret RunSettings
	var parameterGroup string
	var forceGenerate, forceDelete bool
	var err error

	flag.StringVar(&parameterGroup, "p", "", "Semicolon-separated list of parameters in k=v form.")
	flag.StringVar(&ret.basicAuthUser, "bu", "", "The 'user' to use when a remote archive requests Basic Authentication")
	flag.StringVar(&ret.basicAuthPass, "bp", "", "The 'password' to use when a remote archive requests Basic Authentication")
	flag.BoolVar(&ret.inspectionRun, "i", false, "Whether or not to show a list of available parameters for the skeleton")
	flag.BoolVar(&ret.addRegistry, "a", false, "Whether or not to register the template at the given path (can be http/https URLs)")
	flag.BoolVar(&ret.showRegistry, "l", false, "Whether or not to show a list of all available registered templates")
	flag.BoolVar(&forceGenerate, "f", false, "Whether or not to overwrite existing files during generation")
	flag.BoolVar(&forceDelete, "d", false, "Whether or not to delete every existing file in the output path first (only valid when -f is specified)")
	flag.BoolVar(&ret.flagList, "flags", false, "Whether or not to list the flags")
	flag.Parse()

	ret.skeletonPath = flag.Arg(0)
	ret.targetPath = flag.Arg(1)

	if(ret.flagList || ret.showRegistry) {
		return ret, nil
	}

	// if we're not just showing the registry, and not skeleton path is specified...
	if ret.skeletonPath == "" {
		errorMsg := fmt.Sprintf("Skeleton project path not specified")
		return ret, errors.New(errorMsg)
	}

	// if we're actually generating a project, and no target path is specified...
	if !ret.inspectionRun &&
		!ret.addRegistry &&
		ret.targetPath == "" {

		errorMsg := fmt.Sprintf("Target output path not specified")
		return ret, errors.New(errorMsg)
	}

	// set up overwrite strategy
	if(forceGenerate) {
		if(forceDelete) {
			ret.forceGenerate = GENERATE_DELETE
		} else {
			ret.forceGenerate = GENERATE_OVERWRITE
		}
	} else {
		ret.forceGenerate = GENERATE_NONE
	}

	// make parameters, set default project.name
	ret.parameters = make(map[string][]string)
	ret.parameters["project.name"] = []string{filepath.Base(ret.targetPath)}
	ret.parameters["project.createdDate"] = []string{time.Now().Format("2006-01-02")}

	err = parseParametersTo(parameterGroup, ret.parameters)
	if err != nil {
		return ret, err
	}

	return ret, nil
}

/*
  Given a sequence of k=v strings, this parses them out into a map.
*/
func parseParametersTo(parameterGroup string, destination map[string][]string) error {

	var groups, pair, values []string
	var key, value string

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

		key = strings.Trim(pair[0], " ")
		value = strings.Trim(pair[1], " ")
		values = strings.Split(value, ",")

		if len(values) <= 0 {
			errorMsg := fmt.Sprintf("Unable to parse parameters, parameter '%s', value was empty\n", key)
			return errors.New(errorMsg)
		}

		destination[key] = values
	}

	return nil
}
