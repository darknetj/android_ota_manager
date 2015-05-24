// Copperhead OTA Server

package main

import (
  "log"
  "flag"
  "bufio"
  "fmt"
  "os"
  "strings"
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
  userFlag := flag.Bool("add_user", false, "Run CLI for adding user to database")
  flag.Parse()

  // Parse config file
  cfg,err := config.ParseYamlFile(*configPath)
  cfg,err = cfg.Get(*env)
  port,_ := cfg.String("port")
  lib.CheckErr(err, "Config parsing failed")

  // Connect to database
  databasePath,_ := cfg.String("database")
  db := models.InitDb(databasePath)
  defer db.Db.Close()

  if *userFlag {
    // Start CLI to create new user account
    addUser()
  } else {
    // Start server
    server(port)
  }
}

func server(port string) {
  log.Println("--- Started Copperhead OTA Server on port", port, "---")
  // Create auth cookie store
  controllers.InitMiddleware()

  // Create router
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

func addUser() {
  // Add user CLI workflow
  reader := bufio.NewReader(os.Stdin)
  fmt.Print("---\nCopperhead OTA App\n---\n\n")
  fmt.Print("Add a user...\n\n")
  fmt.Print("Enter new Username: ")
  username, _ := reader.ReadString('\n')
  fmt.Print("Enter password: ")
  password,_ := reader.ReadString('\n')
  log.Println("\nSaving...\n")

  // Create user from creds
  models.CreateUser(strings.TrimSpace(username),strings.TrimSpace(password))

  // Exit
  log.Println("Done. Welcome", username)
  log.Println("Exiting")
}
