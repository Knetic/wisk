package main

import (

  "errors"
  "os"
  "bytes"
  "io/ioutil"
  "strings"
  "path/filepath"
)

type TemplatedProject struct {

  files []TemplatedFile
  rootDirectory string
}

func NewTemplatedProject(path string) (*TemplatedProject, error) {

  var ret *TemplatedProject
  var stat os.FileInfo
  var err error

  path, err = filepath.Abs(path)
  if(err != nil) {
    return nil, err
  }

  stat, err = os.Stat(path)
  if(err != nil) {
    return nil, err
  }

  if(!stat.IsDir()) {
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
func (this TemplatedProject) GenerateAt(targetPath string, parameters map[string]string) error {

  return nil
}

func (this TemplatedProject) replaceFileContents(input TemplatedFile, outPath string, parameters map[string]string) error {

  var contentBytes []byte
  var contents string
  var err error

  contentBytes, err = ioutil.ReadFile(input.path)
  if(err != nil) {
    return err
  }

  contents = string(contentBytes)
  contents = this.replaceStringParameters(contents, parameters)

  ioutil.WriteFile(outPath, []byte(contents), input.mode)
  return nil
}

func (this TemplatedProject) replaceStringParameters(input string, parameters map[string]string) string {

  var resultBuffer bytes.Buffer
  var characters chan rune
  var sequence, parameter string
  var exists bool

  characters = make(chan rune)
  go readRunes(input, characters)

  for {

    sequence, exists = readUntil("${{=", characters)
    resultBuffer.WriteString(sequence)

    if(!exists) {
      break
    }

    // read a parameter, then replace it.
    sequence, exists = readUntil("=}}", characters)

    if(!exists) {
      resultBuffer.WriteString("${{=")
      resultBuffer.WriteString(sequence)
      break
    }

    parameter, exists = parameters[sequence]

    if(exists) {
      resultBuffer.WriteString(parameter)
    }
  }

  return resultBuffer.String()
}

func (this *TemplatedProject) getFolderWalker() (func(string, os.FileInfo, error) error) {

  return func(path string, fileStat os.FileInfo, err error) error {

    var file TemplatedFile

    if(fileStat.IsDir()) {
      return nil
    }

    file = NewTemplatedFile(path, this.rootDirectory, fileStat)
    this.files = append(this.files, file)

    return nil
  }
}

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

  var buffer bytes.Buffer
  var sequence string
  var character rune
  var done bool

  for {

      character, done =<- characters

      if(done) {
        return buffer.String(), false
      }

      buffer.WriteString(string(character))

      sequence = buffer.String()
      if(strings.LastIndex(sequence, pattern) >= 0) {

        return sequence, true
      }
  }
}
