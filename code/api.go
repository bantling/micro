package code

// SPDX-License-Identifier: Apache-2.0

import (
  "fmt"
  "io/fs"
  "os"
  "path"

  "github.com/bantling/micro/funcs"
  "github.com/bantling/micro/io/writer"
  "github.com/bantling/micro/union"
)

const (
  errBasePathCreateMsg       = "The base path %s could not be created: %s"
	errBasePathReadMsg         = "The base path %s could not be read: %s"
	errBasePathDeleteMsg       = "The base path %s could not be deleted: %s"
	errDirExistsInBasePathMsg  = "The dir %s already exists in base path %s"
	errCreateDirInBasePathMsg  = "The dir %s could not be created in base path %s: it already exists"
	errCreateSrcFileInDirMsg   = "The file %s could not be created in dir %s: %s"
)

// Generator is the interface describing an easy to use fluent API for generating packages of code.
// If anything goes wrong, the generator panics.
type Generator interface {
  // Return the base path to generate dirs under
	GetBasePath() string
  // Set the base path once to generate dirs under. Setting again panics.
	SetBasePath(basePath string) Generator
  // Creates a directory, or wipes out the contents (assuming they are from a prior run of the same generator).
  // The caller must generate parent directories before child directories, or a panic occurs.
  // This will be the current directory until another call is made
	Dir(name string) Generator
  // Creates a src file in the current directory
  // This will be the current source file under another call is made
	Src(name string) Generator
  // Create global constants in the current source file
	GlobalConsts(constants ...VarDef) Generator
  // Create global vars in the current source file
	GlobalVars(globals ...VarDef) Generator
  // Create types in the current source file
	Types(objects ...ObjectDef) Generator
  // Create funcs in the current source file
	Funcs(funcs ...FuncDef) Generator
}

// BaseGenerator contains the base parts of the Generator interface
type BaseGenerator struct {
  // Base path
  basePath union.Maybe[string]
  // All dirs created
  dirs map[string]bool
  // Current dir
  currentDir string
  // All source files created, mapped under containing dirs
  srcs map[string]map[string]bool
  // Current source
  currentSrc writer.Writer[string]
}

// GetBasePath panics if SetBasePath has not been called
func (bg BaseGenerator) GetBasePath() string {
  return bg.basePath.Get()
}

// SetBasePath panics if the base path has already been set
func (bg *BaseGenerator) SetBasePath(bp string) {
  // Die if base path has already been set
  bg.basePath.SetOrError(bp)

  // Use clean path
  cleanPath := path.Clean(bp)

  // Does the base path exist and contain stuff already from a previous run?
  _, err := os.ReadDir(cleanPath)
  if err != nil {
    if os.IsNotExist(err) {
      // Doesn't exist is ok, we'll just create all the parts that are missing for it
      if err = os.MkdirAll(cleanPath, fs.ModeDir); err != nil {
        // Could not create some part
        panic(fmt.Errorf(errBasePathCreateMsg, cleanPath, err))
      }
    } else {
      // Exists but can't be read
      panic(fmt.Errorf(errBasePathReadMsg, cleanPath, err))
    }
  } else {
    // Yes, we have stuff from previous run, delete last path part and recreate it
    if err = os.RemoveAll(cleanPath); err != nil {
      // Could not delete path
      panic(fmt.Errorf(errBasePathDeleteMsg, cleanPath, err))
    }

    if err = os.Mkdir(cleanPath, fs.ModeDir); err != nil {
      // Could not recreate path
      panic(fmt.Errorf(errBasePathCreateMsg, cleanPath, err))
    }
  }
}

// Dir panics if dir already exists, or creating the dir has an error
func (bg *BaseGenerator) Dir(name string) {
	// Get
	dirPath := path.Clean(bg.basePath.Get() + "/" + name)

	// Die if dir already exists
	if _, haveIt := bg.dirs[dirPath]; haveIt {
		panic(fmt.Errorf(errDirExistsInBasePathMsg, dirPath, bg.basePath))
	}

	if err := os.MkdirAll(dirPath, fs.ModeDir); err != nil {
		// Could not create dir under base path
		panic(fmt.Errorf(errCreateDirInBasePathMsg, dirPath, err))
	}

  bg.dirs[dirPath] = true
}

// Src panics if the file already exists
func (bg *BaseGenerator) Src(name string) {
  if funcs.Map2Test(bg.srcs, bg.currentDir, name) {
    // Die if src already exists
		panic(fmt.Errorf(errCreateSrcFileInDirMsg, name, bg.Dir))
  }

	path := bg.currentDir + "/" + name
	f, err := os.Create(path)
	if err != nil {
		// Could not create src file under dir
		panic(fmt.Errorf(errCreateSrcFileInDirMsg, name, bg.Dir, err))
	}

  funcs.Map2Set(&bg.srcs, bg.currentDir, name, true)
  bg.currentSrc = writer.OfIOWriterAsStrings(f)
}

// CurrentSrc returns the writer for the current src
func (bg BaseGenerator) CurrentSrc() writer.Writer[string] {
  return bg.currentSrc
}
