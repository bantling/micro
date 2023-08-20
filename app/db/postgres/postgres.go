package db

// SPDX-License-Identifier: Apache-2.0

import (
  "database/sql"

  "github.com/bantling/micro/app"
  "github.com/bantling/micro/funcs"
)

const (
  queryDatabaseExists = "select 1 from pg_databases where name = $1"
)

// postgresImpl is the implementation of app.DB
type postgresImpl struct {}

// LoadSchema loads the actual schema from Postgres. Panics on any errors.
// The Postgres Dockerfile will ensure the database is created using CLI tool.
// The connection provided will go to the database.
func (impl postgresImpl) LoadSchema(conn *sql.DB) (cfg app.Configuration) {
  stmt := funcs.MustValue(conn.Prepare(queryDatabaseExists))
  defer stmt.Close

  rows := funcs.MustValue(stmt.Query(""))
  defer rows.Close

  return
}
