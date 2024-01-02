package code

// SPDX-License-Identifier: Apache-2.0

import (
  "fmt"
  "io/fs"
  "path"
  "strings"
  "os"

  "github.com/bantling/micro/io/writer"
)

var (
  errLanguageExistsMsg = "A language named %s has already been registered, the generator is %s"
  errLanguageDoesNotExistMsg = "No language named %s has been registered, choices are: %s"
  errBasePathCreateMsg = "The base path %s could not be created: %s"
  errBasePathReadMsg = "The base path %s could not be read: %s"
  errBasePathDeleteMsg = "The base path %s could not be deleted: %s"
  errCreateDirInBasePathMsg = "The dir %s could not be created in base path %s: %s"
  errCreateSrcFileInDirMsg = "The file %s could not be created in dir %s: %s"

  // generators is a map of languages to generate code in.
  // Each implementation populates this map with an init function.
  generators = map[string]PackageGenerator{}
)

// AddLanguage must be called by the init function of each supported language
func AddLanguage(language string, generator PackageGenerator) {
  if lang, haveIt := generators[language]; haveIt {
    panic(fmt.Errorf(errLanguageExistsMsg, reflect.TypeOf(lang)))
  }

  generators[language] = generator
}

// BasePackageGenerator contains base implementation of PackageGenerator
type BasePackageGenerator struct {
  basePath string // BasePath is the path prefix of zero or more dirs that contain all generated artifacts
  dirs map[string]bool // Dirs is the set of dirs created under BasePath
}

// Construct a generator for a specific language
func Of(language string, basePath string) PackageGenerator {
  generator, haveIt := generators[language]
  if !haveIt {
    panic(
      fmt.Errorf(
        errLanguageDoesNotExistMsg,
        language,
        iter.Maybe(
          stream.ReduceToSlice(
            stream.Map(
              func(t tuple.Two[string, PackageGenerator]) string { return t.T },
            )(iter.OfMap(generators)),
          ),
        ).Get(),
      ),
    )
  }

  // Use clean path
  cleanPath := path.Clean(basePath)

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

  generator.SetBasePath(cleanPath)

  return generator
}

// GetBasePath from PackageGenerator
func (bpg *BasePackageGenerator) GetBasePath() string {
  return bpg.basePath
}

// SetBasePath from PackageGenerator
func (bpg *BasePackageGenerator) SetBasePath(basePath string) {
  if bpg.basePath != "" {
    panic(fmt.Errorf(errBasePathAlreadySetMsg, basePath, bpg.basePath)
  }

  bpg.basePath = basePath
}

// Dir from PackageGenerator
func (bpg BasePackageGenerator) Dir(name string) SrcGenerator {
  // Get
  dirPath := path.Clean(bpg.basePath + "/" + name)

  // Die if dir already exists
  if bpg.dirs[dirpath]

  if err := os.MkdirAll(dirPath, fs.ModeDir); err != nil {
    // Could not create dir under base path
    panic(fmt.Errorf(name, dirPath, err))
  }

  return BaseSrcGenerator{Dir: dirPath}
}

// PackageGenerator
func (bpg BasePackageGenerator) EndProgram() {
  fmt.Println("Program generation ended")
}

// BaseSrcGenerator contains base implementation of SrcGenerator
type BaseSrcGenerator struct {
  Dir string // The path of the dir containing generated source files
}

// Src from SrcGenerator
func (bsg *BaseSrcGenerator) Src(name string) SrcPartsGenerator {
  path := bsg.Dir + "/" + name

  f, err := os.Create(path)
  if err != nil {
    // Could not create src file under dir
      panic(fmt.Errorf(errCreateSrcFileInDirMsg, name, bsg.Dir, err))
  }

  return BaseSrcPartsGenerator{File: path, Writer: writer.OfIOWriterAsStrings(f)}
}

// EndDir from SrcGenerator
func (bsg *BaseSrcGenerator) EndDir() ProgramGenerator

type BaseSrcPartsGenerator struct {
  File string // The path to the file, for error messages
  Writer writer.Writer[string] // The file to write to with unicode strings
}
