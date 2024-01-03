package code

// SPDX-License-Identifier: Apache-2.0

// PackageGenerator is the top level interface describing an easy to use fluent API for generating packages of code.
// Creates a directory, or wipes out the contents (assuming they are from a prior run of the same generator).
//
// The SrcGenerator returns this same instance, allowing additional directories to be created, or stop the whole process of
// generating code, returning to the caller.
//
// The caller must generate parent directories before child directories, or a panic occurs.
type PackageGenerator interface {
  GetBasePath() string // Return the base path to generate dirs under
  SetBasePath(basePath string) // Set the base path once to generate dirs under. Setting again panics.
  Dir(name string) SrcGenerator // Create a dir under base path, it can be a name or relative path. Creating existing dir panics.
  EndPackage() // Doesn't necessarily do anything, mostly for completeness
}

// SrcGenerator creates zero or more source files in the current directory.
// The SrcPartsGenerator returns this same instance, allowing additional source files to be created in the same
// directory, or stop generation in this directory. A given source file can only be defined once, or a panic occurs.
// It is allowable to create an empty directory, or a directory that only contains other directories.
type SrcGenerator interface {
  Src(name string) SrcPartsGenerator
  EndDir() PackageGenerator
}

// SrcPartsGenerator populates a source file, with the following parts in any order:
// - Global constants
// - Global variables
// - Types
// - Functions
type SrcPartsGenerator interface {
  GlobalConsts(constants ...VarDef) SrcPartsGenerator
  GlobalVars(globals ...VarDef) SrcPartsGenerator
  Types(objects ...ObjectDef) SrcPartsGenerator
  Funcs(funcs ...FuncDef) SrcPartsGenerator
  EndSrc() SrcGenerator
}
