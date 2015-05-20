// Copperhead OTA Server

package main

import (
  "log"
  "flag"
  "net/http"
  "gopkg.in/gorp.v1"
  _ "github.com/mattn/go-sqlite3"
  "github.com/olebedev/config"
  "github.com/gorilla/mux"
  "github.com/copperhead-security/android_ota_server/models"
  "github.com/copperhead-security/android_ota_server/controllers"
  "github.com/copperhead-security/android_ota_server/lib"
)

var (
  cfg *config.Config
  db *gorp.DbMap
)

func main() {
  // Parse CLI arguments
  configPath := flag.String("config", "./config.yml", "Path to config file")
  env := flag.String("env", "development", "Run in development or production mode")
  flag.Parse()

  // Parse config file
  cfg,err := config.ParseYamlFile(*configPath)
  cfg,err = cfg.Get(*env)
  lib.CheckErr(err, "Config parsing failed")

  // Connect to database
  databasePath,_ := cfg.String("database")

  db := models.InitDb(databasePath)
  defer db.Db.Close()

  // Start server
  port,_ := cfg.String("port")
  log.Println("--- Started Copperhead OTA Server on port", port, "---")

  router()
}

func router() {
  r := mux.NewRouter()

  // Releases
  r.HandleFunc("/", controllers.Releases)
  r.HandleFunc("/releases", controllers.Releases)
  r.HandleFunc("/releases/{id}", controllers.ShowReleases)
  r.HandleFunc("/releases/edit{id}", controllers.EditReleases)
  r.HandleFunc("/releases/update", controllers.UpdateReleases)
  r.HandleFunc("/releases/new", controllers.NewReleases)
  r.HandleFunc("/releases/create", controllers.CreateReleases)
  r.HandleFunc("/releases/delete", controllers.DeleteReleases)

  // Files
  r.HandleFunc("/files", controllers.Files)
  r.HandleFunc("/files/delete", controllers.DeleteFiles)

  // Users
  r.HandleFunc("/users", controllers.Users)
  r.HandleFunc("/login", controllers.Login)
  r.HandleFunc("/authenticate", controllers.Authenticate)
  r.HandleFunc("/logout", controllers.Logout)

  // Server
  http.Handle("/", r)
  http.ListenAndServe(":8080", nil)
}
