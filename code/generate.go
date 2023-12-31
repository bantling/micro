package code

// SPDX-License-Identifier: Apache-2.0

import (
  "fmt"
)

var (
  errLanguageExistsFmt = "A language named %s has already been registered, the generator is %s"

  // languages is a map of languages to generate code in.
  // Each implementation populates this map with an init function.
  languages = map[string]PackageGenerator{}
)

// AddLanguage must be called by the init function of each supported language
func AddLanguage(language string, pkgGen PackageGenerator) {
  if lang, haveIt := languages[language]; haveIt {
    panic(fmt.Errorf(errLanguageExistsFmt, reflect.TypeOf(lang)))
  }
}

// BasePackageGenerator contains generic code for generating packages
type BasePackageGenerator struct {
  BasePath string // BasePath is the path prefix of zero or more dirs that contain all generated artifacts
}

func Of(language string) {

}

// Dir
func (bpg BasePackageGenerator) Dir(name string) BaseSrcGenerator {

}

// PackageGenerator
func (bpg BasePackageGenerator) EndProgram() {
  fmt.Println("Program generation ended")
}
