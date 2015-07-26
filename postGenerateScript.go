package main

import (
  "path/filepath"
  "os"
  "fmt"
  "os/exec"
)

/*
	Checks for the existence of a post-generation script. If it exists (and is executable)
	this will execute it, with the working directory set to the generated project.
*/
func executePostGenerate(sourcePath string, generatedPath string) error {

	var command *exec.Cmd
  var output []byte
	var scriptPath string
	var workingDirectory string
	var err error

	scriptPath = fmt.Sprintf("%s%s_postGenerate.sh", sourcePath, string(os.PathSeparator))
	scriptPath, err = filepath.Abs(scriptPath)
	if(err != nil) {
		return err
	}

  generatedPath, err = filepath.Abs(generatedPath)
  if(err != nil) {
    return err
  }

  // file doesn't exist, exit quietly.
  if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
    return nil
  }

	workingDirectory, err = os.Getwd()
	if(err != nil) {
		return err
	}

	err = os.Chdir(generatedPath)
	if(err != nil) {
		return err
	}
	defer os.Chdir(workingDirectory)

  command = exec.Command(scriptPath, "")

  output, err = command.CombinedOutput()
  if(err != nil) {
    return err
  }

  fmt.Printf(string(output))
  return nil
}
