package main

import (

  "strings"
  "os"
)

type TemplatedFile struct {

  path string
  mode os.FileMode
}

func NewTemplatedFile(path string, root string, info os.FileInfo) TemplatedFile {

  var ret TemplatedFile

  ret.path = strings.Replace(path, root, "", -1)
  ret.mode = info.Mode()

  return ret
}
