package models

import (
  "net/http"
)

type Release struct {
    Id          int64 `db:"release_id"`
    Created     int64
    BuildDate   int64
    VersionNo   string
    ApiLevel    string
    Channel     string
    Filename    string
    Md5sum      string
    Changelog   string
}

func ReleasesJSONHandler(w http.ResponseWriter, r *http.Request) {}
