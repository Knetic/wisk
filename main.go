package main

import (
  "fmt"
  "flag"
  "os"
)

func main() {

  var root *TemplatedProject
  var err error

  flag.Parse()

  root, err = NewTemplatedProject(flag.Arg(0))
  if(err != nil) {
    exitWith("Unable to read templated project: %s\n", err, 1)
  }

  fmt.Printf("whisk: %v\n", root)
}

func exitWith(message string, err error, code int) {

  errorMsg := fmt.Sprintf(message, err)
  fmt.Fprintf(os.Stderr, errorMsg)
  os.Exit(code)
}
