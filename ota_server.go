// Copperhead OTA Server

package main

import (
  "log"
  "flag"
  "sync"
  "net/http"
  "html/template"
  "path/filepath"
  "gopkg.in/gorp.v1"
  _ "github.com/mattn/go-sqlite3"
  "github.com/olebedev/config"
  "github.com/gorilla/mux"
  "github.com/copperhead-security/android_ota_server/models"
  "github.com/copperhead-security/android_ota_server/lib"
  // "html/template"
)

var (
  cfg *config.Config
  dbmap *gorp.DbMap
  // cookieStore sessions.Store
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
  dbmap := lib.InitDb(databasePath)
  defer dbmap.Db.Close()

  // Start server
  port,_ := cfg.String("port")
  log.Println("--- Started Copperhead OTA Server on port", port, "---")

  router()
}

func router() {
  r := mux.NewRouter()
  r.HandleFunc("/", models.ReleasesJSONHandler)
  r.HandleFunc("/releases.json", models.ReleasesJSONHandler)
  r.HandleFunc("/files", models.FilesHandler)
  http.Handle("/", r)
}

// Cached templates
var cachedTemplates = map[string]*template.Template{}
var cachedMutex sync.Mutex
var funcs = template.FuncMap{}

func T(name string) *template.Template {
  cachedMutex.Lock()
  defer cachedMutex.Unlock()

  if t, ok := cachedTemplates[name]; ok {
      return t
  }

  t := template.New("_base.html").Funcs(funcs)

  t = template.Must(t.ParseFiles(
      "views/layout.html",
      filepath.Join("templates", name),
  ))
  cachedTemplates[name] = t

  return t
}
