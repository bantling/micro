package code

// SPDX-License-Identifier: Apache-2.0

import (
  "github.com/bantling/micro/union"
)

// Dir is a directory of source files
struct Dir {
  // Sources is one or more source files in a directory
  Sources []Source

  // Init is an optional initialization function for the directory.
  // The set of all Init functions execute in some arbitrary order at runtime.
  // Depending on the target language, they may all execute before main starts, or they may execute some time later, such
  // as when files that need them are loaded.
  Init union.Maybe[Func]
}
