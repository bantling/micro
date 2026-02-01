package code

//
// // SPDX-License-Identifier: Apache-2.0
//
// import (
// 	"fmt"
// 	"io/fs"
// 	"os"
// 	"path"
// 	goreflect "reflect"
//
// 	"github.com/bantling/micro/io/writer"
// 	"github.com/bantling/micro/iter"
// 	"github.com/bantling/micro/stream"
// 	"github.com/bantling/micro/tuple"
// )
//
// var (
// 	errLanguageExistsMsg       = "A language named %s has already been registered, the generator is %s"
// 	errLanguageDoesNotExistMsg = "No language named %s has been registered, choices are: %s"
// 	errBasePathAlreadySetMsg   = "The base path is already set to %s"
// 	errBasePathDeleteMsg       = "The base path %s could not be deleted: %s"
// 	errDirExistsInBasePathMsg  = "The dir %s already exists in base path %s"
// 	errCreateDirInBasePathMsg  = "The dir %s could not be created in base path %s: it already exists"
// 	errCreateSrcFileExistsMsg   = "The file %s has already been created in dir %s"
// 	errCreateSrcFileInDirMsg   = "The file %s could not be created in dir %s: %s"
// 	errMethodNotImplementedMsg = "The method %s is not implemented"
//
// 	// generators is a map of languages to generate code in.
// 	// Each implementation populates this map with an init function.
// 	generators = map[string]PackageGenerator{}
// )
//
// // AddLanguage must be called by the init function of each supported language
// func AddLanguage(language string, generator PackageGenerator) {
// 	if lang, haveIt := generators[language]; haveIt {
// 		panic(fmt.Errorf(errLanguageExistsMsg, goreflect.TypeOf(lang)))
// 	}
//
// 	generators[language] = generator
// }
//
// // BasePackageGenerator contains base implementation of PackageGenerator
// type BasePackageGenerator struct {
// 	basePath string          // BasePath is the path prefix of zero or more dirs that contain all generated artifacts
// 	dirs     map[string]bool // Dirs is the set of dirs created under BasePath
// }
//
// // Construct a generator for a specific language
// func Of(language string, basePath string) PackageGenerator {
// 	generator, haveIt := generators[language]
// 	if !haveIt {
// 		panic(
// 			fmt.Errorf(
// 				errLanguageDoesNotExistMsg,
// 				language,
// 				iter.Maybe(
// 					stream.ReduceToSlice(
// 						stream.Map(
// 							func(t tuple.Two[string, PackageGenerator]) string { return t.T },
// 						)(iter.OfMap(generators)),
// 					),
// 				).Get(),
// 			),
// 		)
// 	}
//
// 	// Use clean path
// 	cleanPath := path.Clean(basePath)
//
// 	// Does the base path exist and contain stuff already from a previous run?
// 	_, err := os.ReadDir(cleanPath)
// 	if err != nil {
// 		if os.IsNotExist(err) {
// 			// Doesn't exist is ok, we'll just create all the parts that are missing for it
// 			if err = os.MkdirAll(cleanPath, fs.ModeDir); err != nil {
// 				// Could not create some part
// 				panic(fmt.Errorf(errBasePathCreateMsg, cleanPath, err))
// 			}
// 		} else {
// 			// Exists but can't be read
// 			panic(fmt.Errorf(errBasePathReadMsg, cleanPath, err))
// 		}
// 	} else {
// 		// Yes, we have stuff from previous run, delete last path part and recreate it
// 		if err = os.RemoveAll(cleanPath); err != nil {
// 			// Could not delete path
// 			panic(fmt.Errorf(errBasePathDeleteMsg, cleanPath, err))
// 		}
//
// 		if err = os.Mkdir(cleanPath, fs.ModeDir); err != nil {
// 			// Could not recreate path
// 			panic(fmt.Errorf(errBasePathCreateMsg, cleanPath, err))
// 		}
// 	}
//
// 	generator.SetBasePath(cleanPath)
//
// 	return generator
// }
//
// // GetBasePath from PackageGenerator
// func (bpg BasePackageGenerator) GetBasePath() string {
// 	return bpg.basePath
// }
//
// // SetBasePath from PackageGenerator
// func (bpg *BasePackageGenerator) SetBasePath(basePath string) {
// 	if bpg.basePath != "" {
// 		panic(fmt.Errorf(errBasePathAlreadySetMsg, bpg.basePath))
// 	}
//
// 	bpg.basePath = basePath
// }
//
// // Dir from PackageGenerator
// func (bpg BasePackageGenerator) Dir(name string) SrcGenerator {
// 	// Get
// 	dirPath := path.Clean(bpg.basePath + "/" + name)
//
// 	// Die if dir already exists
// 	if _, haveIt := bpg.dirs[dirPath]; haveIt {
// 		panic(fmt.Errorf(errDirExistsInBasePathMsg, dirPath, bpg.basePath))
// 	}
//
// 	if err := os.MkdirAll(dirPath, fs.ModeDir); err != nil {
// 		// Could not create dir under base path
// 		panic(fmt.Errorf(errCreateDirInBasePathMsg, dirPath, err))
// 	}
//
// 	return BaseSrcGenerator{parent: bpg, Dir: dirPath}
// }
//
// // PackageGenerator
// func (bsg BasePackageGenerator) EndPackage() {
// 	fmt.Println("Program generation ended")
// }
//
// // BaseSrcGenerator contains base implementation of SrcGenerator
// type BaseSrcGenerator struct {
// 	parent BasePackageGenerator
// 	Dir    string // The path of the dir containing generated source files
// }
//
// // Src from SrcGenerator
// func (bsg BaseSrcGenerator) Src(name string) SrcPartsGenerator {
// 	path := bsg.Dir + "/" + name
//
// 	f, err := os.Create(path)
// 	if err != nil {
// 		// Could not create src file under dir
// 		panic(fmt.Errorf(errCreateSrcFileInDirMsg, name, bsg.Dir, err))
// 	}
//
// 	return BaseSrcPartsGenerator{parent: bsg, File: path, Writer: writer.OfIOWriterAsStrings(f)}
// }
//
// // EndDir from SrcGenerator
// func (bsg BaseSrcGenerator) EndDir() PackageGenerator {
// 	return &bsg.parent
// }
//
// // BaseSrcPartsGenerator contains base implementation of SrcPartsGenerator.
// // Only the EndSrc method is implemented, since everything else is language specific.
// // The remaining methods satisfy the interface by panicking
// type BaseSrcPartsGenerator struct {
// 	parent BaseSrcGenerator
// 	File   string                // The path to the file, for error messages
// 	Writer writer.Writer[string] // The file to write to with unicode strings
// }
//
// // GlobalConsts panics
// func (bspg BaseSrcPartsGenerator) GlobalConsts(constants ...VarDef) SrcPartsGenerator {
// 	panic(fmt.Errorf(errMethodNotImplementedMsg, "GlobalConsts"))
// }
//
// // GlobalVars panics
// func (bspg BaseSrcPartsGenerator) GlobalVars(vars ...VarDef) SrcPartsGenerator {
// 	panic(fmt.Errorf(errMethodNotImplementedMsg, "GlobalVars"))
// }
//
// // Types panics
// func (bspg BaseSrcPartsGenerator) Types(types ...ObjectDef) SrcPartsGenerator {
// 	panic(fmt.Errorf(errMethodNotImplementedMsg, "Types"))
// }
//
// // Funcs panics
// func (bspg BaseSrcPartsGenerator) Funcs(types ...FuncDef) SrcPartsGenerator {
// 	panic(fmt.Errorf(errMethodNotImplementedMsg, "Funcs"))
// }
//
// // EndSrc from SrcPartsGenerator
// func (bspg BaseSrcPartsGenerator) EndSrc() SrcGenerator {
// 	return bspg.parent
// }
