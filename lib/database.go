package lib

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
   "github.com/copperhead-security/android_ota_server/models"
    "gopkg.in/gorp.v1"
)

func InitDb(dbPath string) *gorp.DbMap {
  db, err := sql.Open("sqlite3", dbPath)
  CheckErr(err, "sql.Open failed")

  dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
  dbmap.AddTableWithName(models.Release{}, "releases").SetKeys(true, "Id")
  dbmap.AddTableWithName(models.User{}, "users").SetKeys(true, "Id")

  // dbmap.DropTables()
  // err = dbmap.TruncateTables()

  err = dbmap.CreateTablesIfNotExists()
  CheckErr(err, "Create tables failed")

  return dbmap
}
