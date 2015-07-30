package main

import (
	"strconv"
	"strings"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/copperhead-security/android_ota_manager/database"
)

// GET /files
func Files(w rest.ResponseWriter, r *rest.Request) {
	files := fileApi(database.FilesIndex())
	w.WriteJson(apiListResponse(r, files))
}

// GET /files/show/:id
func File(w rest.ResponseWriter, r *rest.Request) {
	id, _ := strconv.ParseInt(r.PathParam("id"), 10, 64)
	file := fileResource(database.FindFile(id))
	w.WriteJson(apiResponse(r, file))
}

// POST /files
func FileCreate(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson("")
}

// GET /builds/{name}
func DownloadFiles(w rest.ResponseWriter, r *rest.Request) {
	// vars := mux.Vars(r)
	// file, _ := database.FindFileByName(vars["name"])
	// path := file.DownloadPath()
	// log.Println("User downloading: ", path)
	// http.ServeFile(w, r, path)
	w.WriteJson("")
}

// GET /changelog/{incremental}.txt
func FileChangelog(w rest.ResponseWriter, r *rest.Request) {
	incremental := r.PathParam("incremental")
	file, err := database.FindFileByIncremental(incremental)
	CheckErr(err, "Find by incremental failed")
	release := database.FindReleaseByFile(file)
	changelog := strings.Join([]string{"Release notes for Copperhead OS #", file.Incremental, "n---n", release.Changelog}, "")
	w.WriteJson(changelog)
}

// POST /files/delete
func FileDelete(w rest.ResponseWriter, r *rest.Request) {
	r.ParseForm()
	id, _ := strconv.ParseInt(r.FormValue("Id"), 10, 64)
	file := database.FindFile(id)

	// Delete from DB
	database.DeleteFile(file)
	w.WriteJson("")
}

// Generate list of files using ApiResource format
func fileApi(items []database.File) []ApiResource {
	resources := make([]ApiResource, 0, len(items))
	for _, item := range items {
		resources = append(resources, fileResource(item))
	}
	return resources
}

// Return a file using ApiResource format
func fileResource(item database.File) ApiResource {
	resource := Resource(item.Id, "File", "/files")
	resource.Attributes = item
	return resource
}
