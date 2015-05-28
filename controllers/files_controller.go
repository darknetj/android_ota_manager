package controllers

import (
    "strconv"
    "net/http"
    "github.com/gorilla/mux"
    "github.com/copperhead-security/android_ota_server/models"
)

// GET /files
func Files(w http.ResponseWriter, r *http.Request) {
    data := map[string]interface{} {"files": models.FilesIndex()}
    R.HTML(w, http.StatusOK, "files", data)
}

// GET /files/show/:id
func ShowFiles(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id,_ := strconv.ParseInt(vars["id"],10,64)
    file := models.FindFile(id)
    data := map[string]interface{} {"file": file}
    R.HTML(w, http.StatusOK, "files_show", data)
}

// GET /files/refresh
func RefreshFiles(w http.ResponseWriter, r *http.Request) {
    models.RefreshBuilds()
    http.Redirect(w, r, "/admin/files", http.StatusFound)
}

// POST /files/delete
func DeleteFiles(w http.ResponseWriter, r *http.Request) {
    r.ParseForm();
    id,_ := strconv.ParseInt(r.FormValue("Id"),10,64)
    file := models.FindFile(id)

    // Delete from DB
    models.DeleteFile(file)
    // TODO: mv file to /builds/deleted directory

    http.Redirect(w, r, "/admin/files", http.StatusFound)
}
