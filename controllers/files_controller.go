package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/copperhead/android_ota_manager/lib"
	"github.com/copperhead/android_ota_manager/models"
	"github.com/gorilla/mux"
)

// GET /files
func Files(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{"files": models.FilesIndex()}
	R.HTML(w, http.StatusOK, "files", data)
}

// GET /builds/{name}
func DownloadFiles(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := fmt.Sprintf("%s/builds/published/%s", os.Getenv("OPENSHIFT_DATA_DIR"), vars["name"])
	log.Println("User downloading: ", path)
	w.Header().Set("Cache-Control", "public, max-age=31536000, no-transform")
	http.ServeFile(w, r, path)
}

// GET /changelog/{incremental}.txt
func ChangelogFiles(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	file, err := models.FindFileByIncremental(vars["incremental"])
	lib.CheckErr(err, "Find by incremental failed")
	release := models.FindReleaseByFile(file)
	changelog := strings.Join([]string{"Release notes for Copperhead OS #", file.Incremental, "\n---\n", release.Changelog}, "")
	fmt.Fprintf(w, changelog)
}

// GET /files/show/:id
func ShowFiles(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)
	file := models.FindFile(id)
	data := map[string]interface{}{"file": file}
	R.HTML(w, http.StatusOK, "files_show", data)
}

// GET /files/refresh
func RefreshFiles(w http.ResponseWriter, r *http.Request) {
	models.RefreshBuilds()
	http.Redirect(w, r, "/admin/files", http.StatusFound)
}

// POST /files/delete
func DeleteFiles(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id, _ := strconv.ParseInt(r.FormValue("Id"), 10, 64)
	file := models.FindFile(id)

	// Delete from DB
	models.DeleteFile(file)

	http.Redirect(w, r, "/admin/files", http.StatusFound)
}
