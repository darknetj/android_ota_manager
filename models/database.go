package models

import (
	"database/sql"
	"log"

	"github.com/copperhead/android_ota_manager/lib"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/gorp.v1"
)

var dbmap *gorp.DbMap
var BuildsPath string

func InitDb(dbPath string, builds string) *gorp.DbMap {
	log.Println("Connecting to database ", dbPath)
	db, err := sql.Open("sqlite3", dbPath)
	lib.CheckErr(err, "sql.Open failed")

	dbmap = &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	dbmap.AddTableWithName(Release{}, "releases").SetKeys(true, "Id")
	dbmap.AddTableWithName(User{}, "users").SetKeys(true, "Id")
	dbmap.AddTableWithName(File{}, "files").SetKeys(true, "Id")

	//dbmap.DropTables()
	//err = dbmap.TruncateTables()

	err = dbmap.CreateTablesIfNotExists()
	lib.CheckErr(err, "Create tables failed")

	BuildsPath = builds
	return dbmap
}
