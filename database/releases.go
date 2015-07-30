package database

import (
	"time"
)

type Release struct {
	Id        int64 `db:"release_id"`
	Created   int64
	FileId    int64
	FileName  string
	Changelog string
	Channel   string
}

func (r Release) GetId() int64 {
    return r.Id
}

func NewRelease() Release {
	release := Release{
		Id:        0,
		Created:   time.Now().UnixNano(),
		Changelog: "- Add bullet point here\n- Another bullet...",
	}
	return release
}

func CreateRelease(release Release) {
	err := dbmap.Insert(&release)
	CheckErr(err, "Insert release failed")
}

func UpdateRelease(release Release) {
	_, err := dbmap.Update(&release)
	CheckErr(err, "Update failed")
}

func DeleteRelease(release Release) {
	_, err := dbmap.Delete(&release)
	CheckErr(err, "Delete failed")
}

func (r Release) ChannelNightly() bool {
	return r.Channel == "NIGHTLY"
}

func (r Release) ChannelSnapshot() bool {
	return r.Channel == "SNAPSHOT"
}

func FindRelease(id int64) Release {
	var release Release
	err := dbmap.SelectOne(&release, "select * from releases where release_id=? LIMIT 1", id)
	CheckErr(err, "Find release failed")
	return release
}

func FindReleaseByFile(file File) Release {
	var release Release
	err := dbmap.SelectOne(&release, "select * from releases where FileId=? LIMIT 1", file.Id)
	CheckErr(err, "Find release by file failed")
	return release
}

func ReleasesIndex() []Release {
	var releases []Release
	_, err := dbmap.Select(&releases, "select * from releases order by release_id DESC")
	CheckErr(err, "Select all releases failed")
	return releases
}

func ReleasesIndexJSON() []map[string]string {
	var releasesJSON []map[string]string
	releases := ReleasesIndex()
	for _, r := range releases {
		f := FindFile(r.FileId)
		newRelease := map[string]string{
			"channel":     r.Channel,
			"filename":    f.Name,
			"url":         f.DownloadUrl(),
			"changes":     f.ChangelogUrl(),
			"md5sum":      f.Md5,
			"api_level":   f.ApiLevel,
			"timestamp":   string(f.BuildDate),
			"incremental": f.Incremental,
		}
		releasesJSON = append(releasesJSON, newRelease)
	}
	return releasesJSON
}

func ReleasesListJSON(device string, channels []string) []map[string]string {
	var releasesJSON []map[string]string
	releases := ReleasesIndex()
	for _, r := range releases {
		f := FindFile(r.FileId)
		if f.Device == device && StringInSlice(r.Channel, channels) {
			newRelease := map[string]string{
				"channel":     r.Channel,
				"filename":    f.Name,
				"url":         f.DownloadUrl(),
				"changes":     f.ChangelogUrl(),
				"md5sum":      f.Md5,
				"api_level":   f.ApiLevel,
				"timestamp":   string(f.BuildDate),
				"incremental": f.Incremental,
			}
			releasesJSON = append(releasesJSON, newRelease)
		}
	}
	return releasesJSON
}
