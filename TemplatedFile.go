package main

import (

  "strings"
  "os"
)

/*
  Represents a single file that is part of a project skeleton.
*/
type TemplatedFile struct {

  path string
  mode os.FileMode
}

/*
  Returns a new templated file, with the given [path], [root] path of the project skeleton,
  and the given filemode [info].
*/
func NewTemplatedFile(path string, root string, info os.FileInfo) TemplatedFile {

  var ret TemplatedFile

  ret.path = strings.Replace(path, root, "", -1)
  ret.mode = info.Mode()

  return ret
}
