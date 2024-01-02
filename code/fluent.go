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
  SetBasePath(basePath string) // Set the base path to generate dirs under
  Dir(name string) SrcGenerator // Create a dir under base path. Dir can be a name or relative path
  EndProgram()
}

// SrcGenerator creates zero or more source files in the current directory.
// The SrcPartsGenerator eventually returns this same instance, allowing additional source files to be created in the same
// directory, or stop generation in this directory. A given source file can only be defined once, or a panic occurs.
// It is allowable to create an empty directory, or a directory that only contains other directories.
type SrcGenerator interface {
  Src(name string) SrcPartsGenerator
  EndDir() ProgramGenerator
}

// SrcPartsGenerator populates a source file, in the following order:
// - Global constants
// - Global variables
// - Objects
// - Functions
type SrcPartsGenerator interface {
  GlobalConst(constants ...Var) SrcGlobalVarsGenerator
  GlobalVars(globals ...Var) SrcObjectsGenerator
  Objects(objects ...Object) SrcFuncsGenerator
  Funcs(funcs ...Func) SrcEnder
}

// SrcGlobalVarsGenerator populates a source file, in the following order:
// - Global variables
// - Objects
// - Functions
type SrcGlobalVarsGenerator interface {
  GlobalVars(globals ...Var) SrcObjectsGenerator
  Objects(objects ...Object) SrcFuncsGenerator
  Funcs(funcs ...Func) SrcEnder
}

// SrcObjectsGenerator populates a source file, in the following order:
// - Objects
// - Functions
type SrcObjectsGenerator interface {
  Objects(objects ...Object) SrcFuncsGenerator
  Funcs(funcs ...Func) SrcEnder
}

// SrcFuncsGenerator populates a source file, in the following order:
// - Functions
type SrcFuncsGenerator interface {
  Funcs(funcs ...Func) SrcEnder
}

// SrcEnder ends the current source file, returning to the parent SrcGenerator instance
type SrcEnder interface {
  EndSrc() SrcGenerator
}
