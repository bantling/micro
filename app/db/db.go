package db

// SPDX-License-Identifier: Apache-2.0

import (
  "database/sql"

  "github.com/bantling/micro/app"
)

// DB defines common operations DBs need to support
type DB Interface {
  // LoadSchema loads the actual schema from the database, to compare against the expected schema. Panics on any errors.
  LoadSchema(conn *sql.DB) app.Configuration
}
