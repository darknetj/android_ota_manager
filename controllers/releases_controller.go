package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/copperhead-security/android_ota_server/models"
	"github.com/gorilla/mux"
)

// GET releases.json
func ReleasesJSON(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"id":     nil,
		"result": models.ReleasesIndexJSON(),
		"error":  nil,
	}
	js, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// POST /releases.json
func PostReleasesJSON(w http.ResponseWriter, r *http.Request) {
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
		"result": models.ReleasesListJSON(device, channels),
		"error":  nil,
	}
	js, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// GET /releases
func Releases(w http.ResponseWriter, r *http.Request) {
	session, _ := CookieStore.Get(r, "auth")
	data := map[string]interface{}{
		"releases": models.ReleasesIndex(),
		"flashes":  session.Flashes(),
	}
	session.Save(r, w)
	R.HTML(w, http.StatusOK, "releases", data)
}

// GET /releases/new
func NewReleases(w http.ResponseWriter, r *http.Request) {
	files := models.Files()
	if len(files) > 0 {
		data := map[string]interface{}{
			"release":  models.NewRelease(),
			"files":    files,
			"title":    "New Release",
			"endpoint": "/admin/releases/create",
		}
		R.HTML(w, http.StatusOK, "releases_form", data)
	} else {
		session, _ := CookieStore.Get(r, "auth")
		session.AddFlash("No files available to create a new release. Upload a build image first.")
		session.Save(r, w)
		http.Redirect(w, r, "/admin/releases", http.StatusFound)
	}
}

// POST /releases/create
func CreateReleases(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fileId, _ := strconv.ParseInt(r.FormValue("FileId"), 10, 64)
	file := models.FindFile(fileId)

	// Generate release
	release := models.Release{
		Created:   time.Now().UnixNano(),
		Changelog: r.FormValue("Changelog"),
		Channel:   r.FormValue("Channel"),
		FileId:    file.Id,
		FileName:  file.Name,
	}

	models.CreateRelease(release)
	models.PublishFile(file)
	go models.RefreshBuilds()

	http.Redirect(w, r, "/admin/releases", http.StatusFound)
}

// POST /releases/edit/:id
func EditReleases(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)
	release := models.FindRelease(id)
	data := map[string]interface{}{
		"release":  release,
		"files":    models.FilesIndex(),
		"title":    "Edit Release",
		"endpoint": "/admin/releases/update",
	}
	R.HTML(w, http.StatusOK, "releases_form", data)
}

// POST /releases/update
func UpdateReleases(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	// Parse form and append to struct
	id, _ := strconv.ParseInt(r.FormValue("Id"), 10, 64)
	fileId, _ := strconv.ParseInt(r.FormValue("FileId"), 10, 64)
	release := models.FindRelease(id)
	file := models.FindFile(fileId)
	release.FileId = file.Id
	release.FileName = file.Name
	release.Channel = r.FormValue("Channel")
	release.Changelog = r.FormValue("Changelog")

	// Append to db
	models.UpdateRelease(release)

	// Redirect
	http.Redirect(w, r, "/admin/releases", http.StatusFound)
}

// POST /releases/delete
func DeleteReleases(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id, _ := strconv.ParseInt(r.FormValue("Id"), 10, 64)
	release := models.FindRelease(id)
	file := models.FindFile(release.FileId)

	// Delete from DB
	models.DeleteRelease(release)
	models.UnpublishFile(file)
	go models.RefreshBuilds()

	http.Redirect(w, r, "/admin/releases", http.StatusFound)
}

// POST /releases/changelog/{id}
func ChangelogReleases(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)
	release := models.FindRelease(id)
	file := models.FindFile(release.FileId)
	url := fmt.Sprintf("/changelog/%s.txt", file.Incremental)
	http.Redirect(w, r, url, http.StatusFound)
}

// POST /v1/builds/get_delta
func GetDeltaReleases(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "{ errors: [ { message: 'Unable to find delta' } ] }")
}
