package code

// SPDX-License-Identifier: Apache-2.0

import (
  "fmt"
  "io/fs"
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
  BasePath string // BasePath is the path prefix of zero or more dirs that contain all generated artifacts
  Dirs map[string]bool // Dirs is the set of dirs created under BasePath
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

  // Does the base path exist and contain stuff already from a previous run?
  _, err := os.ReadDir(basePath)
  if err != nil {
    if os.IsNotExist(err) {
      // Doesn't exist is ok, we'll just create all the parts that are missing for it
      if err = os.MkdirAll(basePath, fs.ModeDir); err != nil {
        // Could not create some part
        panic(fmt.Errorf(errBasePathCreateMsg, basePath, err))
      }
    } else {
      // Exists but can't be read
      panic(fmt.Errorf(errBasePathReadMsg, basePath, err))
    }
  } else {
    // Yes, we have stuff from previous run, delete last path part and recreate it
    if err = os.RemoveAll(basePath); err != nil {
      // Could not delete path
      panic(fmt.Errorf(errBasePathDeleteMsg, basePath, err))
    }

    if err = os.Mkdir(basePath, fs.ModeDir); err != nil {
      // Could not recreate path
      panic(fmt.Errorf(errBasePathCreateMsg, basePath, err))
    }
  }

  generator.SetBasePath(basePath)

  return generator
}

// SetBasePath from PackageGenerator
func (bpg *BasePackageGenerator) SetBasePath(basePath string) {
  bpg.BasePath = basePath
}

// Dir from PackageGenerator
func (bpg BasePackageGenerator) Dir(name string) SrcGenerator {
  path := bpg.BasePath + "/" + name

  if err := os.MkdirAll(path, fs.ModeDir); err != nil {
    // Could not create dir under base path
    panic(fmt.Errorf(name, bpg.BasePath, err))
  }

  return BaseSrcGenerator{Dir: path}
}

// PackageGenerator
func (bpg BasePackageGenerator) EndProgram() {
  fmt.Println("Program generation ended")
}

// BaseSrcGenerator contains base implementation of SrcGenerator
type BaseSrcGenerator struct {
  Dir string // The path of the dir containing generated source files
  File writer.Writer[string] // The file to write to with unicode strings
}

// Src from SrcGenerator
func (bsg *BaseSrcGenerator) Src(name string) SrcPartsGenerator {
  path := bsg.Dir + "/" + name

  f, err := os.Create(path)
  if err != nil {
    // Could not create src file under dir
      panic(fmt.Errorf(errCreateSrcFileInDirMsg, name, bsg.Dir, err))
  }

  bsg.File = writer.OfIOWriterAsStrings(f)

  return nil
}
