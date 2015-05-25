// Copperhead OTA Server

package main

import (
  "log"
  "flag"
  "bufio"
  "fmt"
  "os"
  "strings"
//  "net/http"
  "gopkg.in/gorp.v1"
  _ "github.com/mattn/go-sqlite3"
  "github.com/olebedev/config"
  "github.com/gorilla/mux"
  "github.com/codegangsta/negroni"
  "github.com/unrolled/secure"
  "github.com/copperhead-security/android_ota_server/models"
  "github.com/copperhead-security/android_ota_server/controllers"
  "github.com/copperhead-security/android_ota_server/lib"
)

var (
  cfg *config.Config
  db *gorp.DbMap
  development bool
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
  development = strings.Contains(*env, "development")
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
  admin := mux.NewRouter()

  // Releases API
  r.HandleFunc("/", controllers.Releases)
  r.HandleFunc("/releases.json", controllers.ReleasesJSON)

  // Authentication
  r.HandleFunc("/login", controllers.Login)
  r.HandleFunc("/logout", controllers.Logout)
  r.HandleFunc("/authenticate", controllers.Authenticate)

  // Releases
  admin.HandleFunc("/admin/releases", controllers.Releases)
  admin.HandleFunc("/admin/releases/{id}", controllers.ShowReleases)
  admin.HandleFunc("/admin/releases/edit{id}", controllers.EditReleases)
  admin.HandleFunc("/admin/releases/update", controllers.UpdateReleases)
  admin.HandleFunc("/admin/releases/new", controllers.NewReleases)
  admin.HandleFunc("/admin/releases/create", controllers.CreateReleases)
  admin.HandleFunc("/admin/releases/delete", controllers.DeleteReleases)

  // Files
  admin.HandleFunc("/admin/files", controllers.Files)
  admin.HandleFunc("/admin/files/delete", controllers.DeleteFiles)

  // Users
  admin.HandleFunc("/admin/users", controllers.Users)

  // Negroni
  secureMiddleware := secure.New(secure.Options{
    AllowedHosts: []string{"builds.copperhead.co"},
    SSLRedirect: true,
    STSSeconds: 315360000,
    STSPreload: true,
    FrameDeny: true,
    ContentTypeNosniff: true,
    BrowserXssFilter: true,
    ContentSecurityPolicy: "default-src 'self'; style-src 'self' https://maxcdn.bootstrapcdn.com; font-src https://maxcdn.bootstrapcdn.com;",
    IsDevelopment: development,
  })

  // Create a new negroni for the admin middleware
  r.PathPrefix("/admin").Handler(negroni.New(
    controllers.AuthMiddleware(),
    negroni.Wrap(admin),
  ))

  n := negroni.Classic()
  n.Use(negroni.HandlerFunc(secureMiddleware.HandlerFuncWithNext))
  n.UseHandler(r)
  n.Run(":8080")
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
