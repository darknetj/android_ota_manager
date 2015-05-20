package controllers

import (
    "fmt"
    "net/http"
    "github.com/copperhead-security/android_ota_server/models"
    "github.com/copperhead-security/android_ota_server/lib"
)

// GET /files
func Files(w http.ResponseWriter, r *http.Request) {
    data := map[string]interface{} {"files": models.Files()}
    lib.T("files.html").Execute(w, data)
}

// POST /files/delete
func DeleteFiles(w http.ResponseWriter, r *http.Request) {
    r.ParseForm();
    buildName := r.FormValue("buildName")

    // TODO: mv file to /builds/deleted directory

    url := fmt.Sprintf("/files?%s", buildName)
    http.Redirect(w, r, url, http.StatusFound)
}
