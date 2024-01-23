package golang

// SPDX-License-Identifier: Apache-2.0

import (
  "github.com/bantling/micro/code"
)

// GoGenerator generates Go code
type GoGenerator struct {
  bg *code.BaseGenerator
}

// Construct a GoGenerator
func Of() GoGenerator {
  return GoGenerator{bg: &code.BaseGenerator{}}
}

// GetBasePath implements code.Generator using code.BaseGenerator
func (gg GoGenerator) GetBasePath() string {
  return gg.bg.GetBasePath()
}

// SetBasePath implements code.Generator using code.BaseGenerator
func (gg *GoGenerator) SetBasePath(basePath string) code.Generator {
  gg.bg.SetBasePath(basePath)
  return gg
}

// Dir implements code.Generator using code.BaseGenerator
func (gg *GoGenerator) Dir(name string) code.Generator {
  gg.bg.Dir(name)
  return gg
}

// Src implements code.Generator using code.BaseGenerator
func (gg *GoGenerator) Src(name string) code.Generator {
  gg.bg.Src(name)
  return gg
}

// GlobalConsts implements code.Generator
func (gg *GoGenerator) GlobalConsts(constants ...code.VarDef) code.Generator {
  src := gg.bg.CurrentSrc()

  // switch len(constants) {
  // case 1:
  //   c = constants[0]
  //   src.Write("const ", typeMap[c.Type], " = ",

  // // Start const block
  // src.Write("const (\n")

  return gg
}

// GlobalVars implements code.Generator
func (gg *GoGenerator) GlobalVars(globals ...code.VarDef) code.Generator {
  return gg
}

// Types implements code.Generator
func (gg *GoGenerator) Types(objects ...code.ObjectDef) code.Generator {
  return gg
}

// Funcs implements code.Generator
func (gg *GoGenerator) Funcs(funcs ...code.FuncDef) code.Generator {
  return gg
}
