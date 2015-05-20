package models

import (
    "fmt"
    "time"
    "github.com/copperhead-security/android_ota_server/lib"
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

func (r *Release) Url() string {
    url := fmt.Sprintf("https://builds.copperhead.co/builds/%s", r.Filename)
    return url
}

func (r *Release) ChangelogUrl() string {
    url := fmt.Sprintf("https://builds.copperhead.co/changelog/%s.txt", r.VersionNo)
    return url
}

func LastVersionNo() string {
    var release Release
    _ = dbmap.SelectOne(&release, "SELECT * FROM releases ORDER BY release_id DESC LIMIT 1")
    return release.VersionNo
}

func NewRelease() Release {
    release := Release{
        Id: 1,
        Created: time.Now().UnixNano(),
        BuildDate: time.Now().UnixNano(),
        VersionNo: LastVersionNo(),
        ApiLevel: "21",
        Channel: "nightly",
        Filename: "copperhead-2015-05-13-NIGHTLY.zip",
        Md5sum: "4e0a335b378035d12cb6626b6623072b",
        Changelog: "Release Notes for Copperhead 0.0.1\n---\n- Add notes here...",
    }
    return release
}

func CreateRelease(release Release) {
    // Insert into db
    err := dbmap.Insert(&release)
    lib.CheckErr(err, "Insert release failed")
}

func UpdateRelease(release Release) {
    _, err := dbmap.Update(&release)
    lib.CheckErr(err, "Update failed")
}

func DeleteRelease(release Release) {
    _, err := dbmap.Delete(&release)
    lib.CheckErr(err, "Delete failed")
}

func FindRelease(id int64) Release {
    var release Release
    err := dbmap.SelectOne(&release, "select * from releases where release_id=?", id)
    lib.CheckErr(err, "Find release failed")
    return release
}


func ReleasesIndex() []Release {
    var releases []Release
    _, err := dbmap.Select(&releases, "select * from releases order by release_id DESC")
    lib.CheckErr(err, "Select all releases failed")
    return releases
}

func ReleasesListJSON() []map[string]string {
    var releasesJSON []map[string]string
    releases := ReleasesIndex()
    for _, r := range releases {
        newRelease := map[string]string{
            "filename": r.Filename,
            "url": r.Url(),
            "changes": r.ChangelogUrl(),
            "md5sum": r.Md5sum,
            "api_level": r.ApiLevel,
            "timestamp": string(r.BuildDate),
            "channel": r.Channel,
            "incremental": r.VersionNo,
        }
        releasesJSON = append(releasesJSON, newRelease)
    }
    return releasesJSON
}

