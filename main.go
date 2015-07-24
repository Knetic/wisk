package main

import (
  "fmt"
  "os"
)

func main() {

  var project *TemplatedProject
  var settings RunSettings
  var err error

  settings, err = FindRunSettings()
  if(err != nil) {
    exitWith("Unable to parse run arguments: %s\n", err, 1)
    return
  }

  project, err = NewTemplatedProject(settings.skeletonPath)
  if(err != nil) {
    exitWith("Unable to read templated project: %s\n", err, 1)
    return
  }

  err = project.GenerateAt(settings.targetPath, settings.parameters)
  if(err != nil) {
    exitWith("Unable to generate project: %s\n", err, 1)
    return
  }
}

func exitWith(message string, err error, code int) {

  errorMsg := fmt.Sprintf(message, err)
  fmt.Fprintf(os.Stderr, errorMsg)
  os.Exit(code)
}
