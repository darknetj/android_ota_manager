package main

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/bitly/go-simplejson"
	"github.com/copperhead-security/android_ota_manager/database"
)

// GET releases.json
func Releases(w rest.ResponseWriter, r *rest.Request) {
	releases := releaseApi(database.ReleasesIndex())
	w.WriteJson(apiListResponse(r, releases))
}

// GET releases/:id
func Release(w rest.ResponseWriter, r *rest.Request) {
	id, _ := strconv.ParseInt(r.PathParam("id"), 10, 64)
	release := releaseResource(database.FindRelease(id))
	w.WriteJson(apiResponse(r, release))
}

// POST /releases.json
func ReleasesPOST(w rest.ResponseWriter, r *rest.Request) {
	// Parse params from JSON
	val := new(struct {
		Method string           `json:"method"`
		Params *simplejson.Json `json:"params"`
	})
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&val)
	if err != nil {
		log.Println(err)
	}

	// Extract filters to apply to releases
	device, _ := val.Params.Get("device").String()
	channels, _ := val.Params.Get("channels").StringArray()

	// Prep JSON data
	data := map[string]interface{}{
		"id":     nil,
		"result": database.ReleasesListJSON(device, channels),
		"error":  nil,
	}
	js, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	w.WriteJson(js)
}

// POST /releases/create
func ReleaseCreate(w rest.ResponseWriter, r *rest.Request) {
	r.ParseForm()
	fileId, _ := strconv.ParseInt(r.FormValue("FileId"), 10, 64)
	file := database.FindFile(fileId)

	// Generate release
	release := database.Release{
		Created:   time.Now().UnixNano(),
		Changelog: r.FormValue("Changelog"),
		Channel:   r.FormValue("Channel"),
		FileId:    file.Id,
		FileName:  file.Name,
	}

	database.CreateRelease(release)
	database.PublishFile(file)
	go database.RefreshBuilds()
	w.WriteJson("")
}

// POST /releases/update
func ReleaseUpdate(w rest.ResponseWriter, r *rest.Request) {
	r.ParseForm()

	// Parse form and append to struct
	id, _ := strconv.ParseInt(r.FormValue("Id"), 10, 64)
	fileId, _ := strconv.ParseInt(r.FormValue("FileId"), 10, 64)
	release := database.FindRelease(id)
	file := database.FindFile(fileId)
	release.FileId = file.Id
	release.FileName = file.Name
	release.Channel = r.FormValue("Channel")
	release.Changelog = r.FormValue("Changelog")

	// Append to db
	database.UpdateRelease(release)

	// Redirect
	w.WriteJson("")
}

// POST /releases/delete
func ReleaseDelete(w rest.ResponseWriter, r *rest.Request) {
	r.ParseForm()
	id, _ := strconv.ParseInt(r.FormValue("Id"), 10, 64)
	release := database.FindRelease(id)
	file := database.FindFile(release.FileId)

	// Delete from DB
	database.DeleteRelease(release)
	database.UnpublishFile(file)
	go database.RefreshBuilds()

	// http.Redirect(w, r, "/admin/releases", http.StatusFound)
	w.WriteJson("")
}

// POST /releases/changelog/{id}
func ReleaseChangelog(w rest.ResponseWriter, r *rest.Request) {
	// id, _ := strconv.ParseInt(r.PathParam("id"), 10, 64)
	// release := database.FindRelease(id)
	// file := database.FindFile(release.FileId)
	// url := fmt.Sprintf("/changelog/%s.txt", file.Incremental)
	// http.Redirect(w, r, url, http.StatusFound)
	w.WriteJson("")
}

// POST /v1/builds/get_delta
func GetDeltaReleases(w rest.ResponseWriter, r *rest.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteJson("{ errors: [ { message: 'Unable to find delta' } ] }")
}

// Generate list of releases using ApiResource format
func releaseApi(items []database.Release) []ApiResource {
	resources := make([]ApiResource, 0)
	for _, item := range items {
		resources = append(resources, releaseResource(item))
	}
	return resources
}

// Return a release using ApiResource format
func releaseResource(item database.Release) ApiResource {
	resource := Resource(item.Id, "Release", "/releases")
	resource.Attributes = item
	return resource
}
