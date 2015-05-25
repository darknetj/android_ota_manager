package controllers

import (
    "fmt"
    "log"
    "strconv"
    "time"
    "net/http"
    "encoding/json"
    "github.com/copperhead-security/android_ota_server/models"
    "github.com/gorilla/mux"
)

// GET releases.json
func ReleasesJSON(w http.ResponseWriter, r *http.Request) {
    data := map[string]interface{} {
        "id": nil,
        "result": models.ReleasesListJSON(),
        "error": nil,
    }
    js, _ := json.Marshal(data)
    w.Header().Set("Content-Type", "application/json")
    w.Write(js)
}

// POST /releases.json
func PostReleasesJSON(w http.ResponseWriter, r *http.Request) {
    r.ParseForm();

    // Print form values
    fmt.Printf("%+v\n", r.Form)
    for key, values := range r.Form {
       for _, value := range values {
            fmt.Println(key, value)
       }
    }

    // Prep JSON data
    data := map[string]interface{} {
        "id": nil,
        "result": models.ReleasesListJSON(),
        "error": nil,
    }
    js,_ := json.Marshal(data)
    w.Header().Set("Content-Type", "application/json")
    w.Write(js)
}

// GET /releases
func Releases(w http.ResponseWriter, r *http.Request) {
    data := map[string]interface{} {
      "releases": models.ReleasesIndex(),
    }
    R.HTML(w, http.StatusOK, "releases", data)
}

// GET /releases/new
func NewReleases(w http.ResponseWriter, r *http.Request) {
    data := map[string]interface{} {
         "release": models.NewRelease(),
         "files": models.Files(),
         "title": "New Release",
         "endpoint": "/releases/create",
    }
    R.HTML(w, http.StatusOK, "releases_form", data)
}

// POST /releases/create
func CreateReleases(w http.ResponseWriter, r *http.Request) {
    r.ParseForm();

    // Generate release
    release := models.Release{
        Created: time.Now().UnixNano(),
        VersionNo: r.FormValue("VersionNo"),
        ApiLevel: r.FormValue("ApiLevel"),
        Channel: r.FormValue("Channel"),
        Filename: r.FormValue("Filename"),
        Md5sum: "x3u3j3j3j",
        Changelog: r.FormValue("Changelog"),
    }

    models.CreateRelease(release)

    url := fmt.Sprintf("/admin/releases")
    http.Redirect(w, r, url, http.StatusFound)
}

// GET /releases/show/:id
func ShowReleases(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id,_ := strconv.ParseInt(vars["id"],10,64)
    release := models.FindRelease(id)
    data := map[string]interface{} {"release": release}
    R.HTML(w, http.StatusOK, "releases_show", data)
}

// POST /releases/edit/:id
func EditReleases(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id,_ := strconv.ParseInt(vars["id"],10,64)
    release := models.FindRelease(id)
    data := map[string]interface{} {
            "release": release,
            "files": models.Files(),
            "title": "Edit Release",
            "endpoint": "/admin/releases/update",
    }
    R.HTML(w, http.StatusOK, "releases_form", data)
}

// POST /releases/update
func UpdateReleases(w http.ResponseWriter, r *http.Request) {
    r.ParseForm();

    // Parse form and append to struct
    id,_ := strconv.ParseInt(r.FormValue("Id"),10,64)
    log.Println(id)
    release := models.FindRelease(id)
    release.VersionNo = r.FormValue("VersionNo")
    release.ApiLevel = r.FormValue("ApiLevel")
    release.Channel = r.FormValue("Channel")
    release.Filename = r.FormValue("Filename")

    // Append to db
    models.UpdateRelease(release)

    // Redirect
    url := fmt.Sprintf("/admin/releases/show/%i", id)
    http.Redirect(w, r, url, http.StatusFound)
}

// POST /releases/delete
func DeleteReleases(w http.ResponseWriter, r *http.Request) {
    r.ParseForm();
    id,_ := strconv.ParseInt(r.FormValue("Id"),10,64)
    release := models.FindRelease(id)

    // Delete from DB
    models.DeleteRelease(release)

    url := fmt.Sprintf("/admin/releases")
    http.Redirect(w, r, url, http.StatusFound)
}
