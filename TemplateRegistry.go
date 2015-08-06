package main

import (
  "io/ioutil"
  "os"
  "fmt"
  "errors"
  "strings"
  "path/filepath"
  "github.com/mitchellh/go-homedir"
  "github.com/jhoonb/archivex"
)

/*
  Manages templates that live in a persistent template directory on the local disk.
*/
type TemplateRegistry struct {

  templatePaths StringSet
  path string
}

/*
  Creates a new registry located in the user's home directory.
  If this directory does not exist, it is created and left empty, returning and empty TemplateRegistry struct.
  If it exists, this will populate a new TemplateRegistry struct and return it.
  If the path is not readable, this returns an empty struct.
*/
func NewTemplateRegistry() *TemplateRegistry {

  var ret *TemplateRegistry
  var files []os.FileInfo
  var name string
  var err error

  ret = new(TemplateRegistry)

  ret.path, err = getRegistryPath()
  if(err != nil) {
    return ret
  }

  err = os.MkdirAll(ret.path, 0700)
  if(err != nil) {
    return ret
  }

  files, err = ioutil.ReadDir(ret.path)
  if(err != nil) {
    return ret
  }

  for _, file := range files {

    name = file.Name()

    // strip extension
    extensionIndex := len(filepath.Ext(name))
    if(extensionIndex > 0) {

      nameLen := len(name)
      name = name[0:nameLen-extensionIndex]
    }

    ret.templatePaths.Add(name)
  }

  return ret
}

/*
  Returns true if this registry contains a template of the given [name],
  false otherwise.
*/
func (this TemplateRegistry) Contains(name string) bool {
  return this.templatePaths.Contains(name)
}

/*
  Returns true if the given path is something that should be interpreted
  as a registry template.
*/
func (this TemplateRegistry) IsPathRegistry(path string) bool {
  return !strings.Contains(path, string(os.PathSeparator)) && !strings.Contains(path, ".")
}

/*
  Returns a path for the expanded template identified by [name].
  The path returned is temporary, and is not the "actual" location of the template,
  merely where it can be immediately read.
  The returned path should not be persisted.
*/
func (this TemplateRegistry) GetTemplatePath(name string) (string, error) {

  var path string

  if(!this.Contains(name)) {
    errorMsg := fmt.Sprintf("Cannot find any template by the name '%s'\n", name)
    return "", errors.New(errorMsg)
  }

  path = (this.path + name + ".zip")
  return filepath.Abs(path)
}

/*
  Registers the given [path] in the registry, by copying the archive file to it.
  If the file is not an archive, or it cannot be read, or the registry cannot be written,
  an error is returned.
*/
func (this *TemplateRegistry) RegisterTemplate(path string) (string, error) {

  var targetPath, name string
  var err error

  // if the given path is a directory (not a zip file),
  // archive it and prepare for registration.
  if(!strings.HasSuffix(path, archiveMarker)) {

    path, err = archivePath(path)
    if(err != nil) {
      return "", err
    }
  }

  name = filepath.Base(path)
  targetPath = fmt.Sprintf("%s%s%s", this.path, string(os.PathSeparator), name)

  _, err = CopyFile(path, targetPath)
  return name, err
}

func archivePath(path string) (string, error) {

  var zip archivex.ZipFile
  var name string
  var tempPath string
  var err error

  tempPath, err = ioutil.TempDir("", "")
  if(err != nil) {
    return "", err
  }

  name = filepath.Base(path)
  tempPath = fmt.Sprintf("%s%s%s.zip", tempPath, string(os.PathSeparator), name)

  zip.Create(tempPath)
  zip.AddAll(path, false)
  zip.Close()

  return tempPath, nil
}

/*
  Returns a slice representing every template in this registry.
*/
func (this TemplateRegistry) GetTemplateList() []string {
  return this.templatePaths.GetSlice()
}

func getRegistryPath() (string, error) {

  var ret string
  var err error

  ret, err = homedir.Dir()
  if(err != nil) {
    errorMsg := fmt.Sprintf("Unable to determine user home directory: %s\n", err.Error())
    return "", errors.New(errorMsg)
  }

  ret, err = homedir.Expand(ret)
  if(err != nil) {
    errorMsg := fmt.Sprintf("Unable to expand home directory: %s\n", err.Error())
    return "", errors.New(errorMsg)
  }

  return fmt.Sprintf("%s%s.wisk/", ret, string(os.PathSeparator)), nil
}
