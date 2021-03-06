package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	PLACEHOLDER_OPEN  = "${{="
	PLACEHOLDER_CLOSE = "=}}"

	PARAMETER_JOIN_OPEN  = "["
	PARAMETER_JOIN_CLOSE = "]"

	PLACEHOLDER_ITERATIVE_VALUE      = "${{value}}"
	PLACEHOLDER_ITERATIVE_RECURSE    = "${{recurse}}"
	PLACEHOLDER_ITERATIVE_RECURSE_LN = "${{recurse}}\n"

	archiveMarker = ".zip"
)

/*
  Represents an entire skeleton project, capable of generating new projects.
*/
type TemplatedProject struct {
	files         []TemplatedFile
	rootDirectory string

	missingParameters StringSet
	incorrectParameters StringSet
}

/*
  Creates a new skeleton project rooted at the given [path].
  Every file below that path (of any size or location) is included.
*/
func NewTemplatedProject(path string) (*TemplatedProject, error) {

	var tempDir string
	var err error

	path, err = filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	// extract archive to temporary directory, then read it.
	if strings.HasSuffix(path, archiveMarker) {

		tempDir, err = ioutil.TempDir("", "")
		if err != nil {
			return nil, err
		}

		Unzip(path, tempDir)
		return createTemplatedProjectFromFile(tempDir)
	}
	return createTemplatedProjectFromFile(path)
}

func createTemplatedProjectFromFile(path string) (*TemplatedProject, error) {

	var ret *TemplatedProject
	var stat os.FileInfo
	var err error

	stat, err = os.Stat(path)
	if err != nil {
		return nil, err
	}

	if !stat.IsDir() {
		return nil, errors.New("Path is not a directory")
	}

	ret = new(TemplatedProject)
	ret.rootDirectory = path

	err = filepath.Walk(path, ret.getFolderWalker())
	return ret, err
}

/*
  Creates a copy of this project's template at the given [targetPath]
  using the given [parameters].
*/
func (this *TemplatedProject) GenerateAt(targetPath string, parameters map[string][]string) error {

	var file TemplatedFile
	var inputPath, outputPath string
	var err error

	for _, file = range this.files {

		outputPath = (targetPath + file.path)
		outputPath, err = filepath.Abs(outputPath)
		if err != nil {
			return err
		}

		inputPath = (this.rootDirectory + file.path)
		outputPath = this.replaceStringParameters(outputPath, parameters)

		if strings.HasSuffix(outputPath, "/") {
			fmt.Printf("Could not create file at '%s<empty file name>', because the file name is empty.\n", outputPath)
			continue
		}

		err = this.replaceFileContents(inputPath, outputPath, file.mode, parameters)

		if err != nil {
			return err
		}
	}

	return nil
}

/*
  Returns a deduplicated list of all parameters used by this skeleton.
*/
func (this TemplatedProject) FindParameters() ([]string, error) {

	var file TemplatedFile
	var parameters StringSet
	var contentBytes []byte
	var inputPath string
	var err error

	for _, file = range this.files {

		inputPath = (this.rootDirectory + file.path)
		parameters.Add(this.findParametersInString(file.path)...)

		contentBytes, err = ioutil.ReadFile(inputPath)
		if err != nil {
			return nil, err
		}

		parameters.Add(this.findParametersInString(string(contentBytes))...)
	}

	return parameters.GetSlice(), nil
}

func (this *TemplatedProject) findParametersInString(target string) []string {

	var parameters StringSet
	var characters chan rune
	var sequence string
	var exists bool

	characters = make(chan rune)
	go readRunes(target, characters)

	for {

		sequence, exists = readUntil(PLACEHOLDER_OPEN, characters)
		if !exists {
			break
		}

		// read a parameter, then replace it.
		sequence, exists = readUntil(PLACEHOLDER_CLOSE, characters)
		if !exists {
			break
		}

		// content placeholder?
		if sequence[0:1] == ":" || sequence[0:1] == ";" {
			sequence = sequence[1:len(sequence)]
		}

		// joined list?
		index := strings.LastIndex(sequence, PARAMETER_JOIN_OPEN)
		if index > 0 {
			sequence = sequence[0:index]
		}

		parameters.Add(sequence)
	}

	return parameters.GetSlice()
}

/*
  Reads the contents of the file at [inPath], replaces placeholders with the given [parameters],
  then writes the results to the given [outPath] (with the given [mode]).
  Any directories that do not exist in the [outPath] tree will be created.
*/
func (this *TemplatedProject) replaceFileContents(inPath, outPath string, mode os.FileMode, parameters map[string][]string) error {

	var contentBytes []byte
	var contents, path, base string
	var err error

	path, err = filepath.Abs(inPath)
	if err != nil {
		return err
	}

	contentBytes, err = ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	// ensure base path exists
	base = fmt.Sprintf("%s%s", string(os.PathSeparator), filepath.Base(outPath))
	index := strings.LastIndex(outPath, base)
	base = (outPath[0:index])

	err = os.MkdirAll(base, 0755)
	if err != nil {
		return err
	}

	// write replaced contents
	contents = string(contentBytes)
	contents = this.replaceStringParameters(contents, parameters)

	err = ioutil.WriteFile(outPath, []byte(contents), mode)
	return err
}

/*
  Replaces all placeholders in the given [input] with their equivalent values in [parameters],
  returning the resultant string.
*/
func (this *TemplatedProject) replaceStringParameters(input string, parameters map[string][]string) string {

	var characters chan rune
	var sequence, separator, parameterName string
	var parameterValues []string
	var exists bool

	characters = make(chan rune)
	go readRunes(input, characters)

	resultBuffer := bytes.NewBuffer([]byte{})

	for {

		sequence, exists = readUntil(PLACEHOLDER_OPEN, characters)
		resultBuffer.WriteString(sequence)

		if !exists {
			break
		}

		// read a parameter, then replace it.
		sequence, exists = readUntil(PLACEHOLDER_CLOSE, characters)

		if !exists {
			resultBuffer.WriteString(PLACEHOLDER_OPEN)
			resultBuffer.WriteString(sequence)
			break
		}

		// write parameter. If the parameter is unspecified, add it to the list of missing parameters.

		// check if the parameter has a separator
		exists, parameterName, separator = determineParameterSeparator(sequence)
		if exists {

			parameterValues, exists = parameters[parameterName]

			if !exists {
				this.missingParameters.Add(parameterName)
			} else {

				if(len(parameterValues) == 1) {
					this.incorrectParameters.Add(parameterName)
				}

				sequence = strings.Join(parameterValues, separator)
				resultBuffer.WriteString(sequence)
			}
			continue
		}

		// is this a content placeholder?
		if sequence[0:1] == ":" {

			parameterName = sequence[1:len(sequence)]
			sequence, exists = readUntil(fmt.Sprintf("${{=;%s=}}", parameterName), characters)

			// no closing content identifier?
			if !exists {
				resultBuffer.WriteString(sequence)
				continue
			}

			parameterValues, exists = parameters[parameterName]
			if !exists {
				this.missingParameters.Add(parameterName)
				continue
			}

			sequence = this.fillContentPlaceholder(parameterValues, sequence, parameters)
			resultBuffer.WriteString(sequence)
			continue
		}

		// this must be a normal parameter.
		parameterValues, exists = parameters[sequence]
		if !exists {
			this.missingParameters.Add(sequence)
		} else {
			resultBuffer.WriteString(parameterValues[0])
		}
	}

	return resultBuffer.String()
}

/*
	Given the contents of a content placeholder and a specific parameter name / parameters, this
	returns a string which represents the exploded/replaced content.
*/
func (this *TemplatedProject) fillContentPlaceholder(parameterValues []string, contents string, parameters map[string][]string) string {

	var ret bytes.Buffer
	var contentTemplate, replaced, current string
	var recurseTemplate string
	var recursive bool

	// explode any inner content placeholders or regular placeholders.
	contentTemplate = this.replaceStringParameters(contents, parameters)
	contentTemplate = strings.TrimLeft(contentTemplate, " \t\n")

	current = PLACEHOLDER_ITERATIVE_RECURSE_LN
	recursive = strings.Index(contentTemplate, PLACEHOLDER_ITERATIVE_RECURSE) > 0

	if recursive {

		if strings.Index(contentTemplate, PLACEHOLDER_ITERATIVE_RECURSE_LN) > 0 {
			recurseTemplate = PLACEHOLDER_ITERATIVE_RECURSE_LN
		} else {
			recurseTemplate = PLACEHOLDER_ITERATIVE_RECURSE
		}
	}

	// iterate over every item in the list of parameter values, replacing "${{value}}" with the current item.
	for _, value := range parameterValues {

		replaced = strings.Replace(contentTemplate, PLACEHOLDER_ITERATIVE_VALUE, value, 1)

		// if this is a recursive template, replace the recursion indicator
		// with the replaced value.
		if recursive {

			current = strings.Replace(current, recurseTemplate, replaced, 1)
		} else {
			// otherwise, just append.
			ret.WriteString(replaced)
		}
	}

	if recursive {

		current = strings.Replace(current, PLACEHOLDER_ITERATIVE_RECURSE, "", 1)
		return current
	}
	return ret.String()
}

func determineParameterSeparator(parameter string) (exists bool, name string, separator string) {

	var start, end int

	start = strings.LastIndex(parameter, PARAMETER_JOIN_OPEN)
	end = strings.LastIndex(parameter, PARAMETER_JOIN_CLOSE)

	if start > 0 && end > 0 {

		separator = parameter[start+1 : end]
		if len(separator) <= 0 {
			separator = string(os.PathSeparator)
		}

		return true, parameter[0:start], separator
	}

	return false, "", ""
}

/*
  Creates a directory walker that discovers files and appends then into this templatedProject's
  list of files.
*/
func (this *TemplatedProject) getFolderWalker() func(string, os.FileInfo, error) error {

	var scmPaths []string
	var scmPath bytes.Buffer

	scmPath.WriteString(".git")
	scmPath.WriteRune(os.PathSeparator)

	scmPaths = []string {
		scmPath.String(),
	}

	return func(path string, fileStat os.FileInfo, err error) error {

		var file TemplatedFile

		if fileStat.IsDir() {
			return nil
		}

		// Check to see if the path is underneath an SCM root (like ".git")
		for _, forbiddenPath := range scmPaths {
			if(strings.Contains(path, forbiddenPath)) {
				return nil
			}
		}

		file = NewTemplatedFile(path, this.rootDirectory, fileStat)
		this.files = append(this.files, file)

		return nil
	}
}

/*
  Reads all runes individually from the given [input],
  writing each of them into the given [results] channel.
  Closes the channel after all runes have been read.
*/
func readRunes(input string, results chan rune) {

	for _, character := range input {

		results <- character
	}

	close(results)
}

/*
  Reads from the given channel until the given [pattern] is found.
  Returns a string representing all characters not part of the pattern,
  and a bool representing whether or not the end of the channel was reached
  before a pattern was found.
*/
func readUntil(pattern string, characters chan rune) (string, bool) {

	var sequence string
	var character rune
	var index int
	var done bool

	buffer := bytes.NewBuffer([]byte{})

	for {

		character, done = <-characters

		if !done {
			return buffer.String(), false
		}

		buffer.WriteString(string(character))
		sequence = buffer.String()
		index = strings.LastIndex(sequence, pattern)

		if index >= 0 {

			// remove pattern from sequence
			return sequence[0:index], true
		}
	}
}
