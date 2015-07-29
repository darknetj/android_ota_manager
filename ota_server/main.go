// Copperhead OTA Server

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/codegangsta/negroni"
	"github.com/copperhead-security/android_ota_manager/database"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/olebedev/config"
	"github.com/unrolled/secure"
	"gopkg.in/gorp.v1"
)

var (
	cfg         *config.Config
	db          *gorp.DbMap
	development bool
)

func main() {
	// Parse CLI arguments
	configPath := flag.String("config", "./config.yml", "Path to config file")
	env := flag.String("env", "development", "Run in development or production mode")
	userFlag := flag.Bool("add_user", false, "Run CLI for adding user to database")
	testFlag := flag.Bool("test", false, "Run test script to simulate client")
	flag.Parse()

	// Parse config file
	cfg, err := config.ParseYamlFile(*configPath)
	cfg, err = cfg.Get(*env)
	port, _ := cfg.String("port")
	development = strings.Contains(*env, "development")
	templates, _ := cfg.String("templates")
	builds, _ := cfg.String("builds")
	lib.CheckErr(err, "Config parsing failed")

	// Connect to database
	databasePath, _ := cfg.String("database")
	db := database.InitDb(databasePath, builds)
	go database.RefreshBuilds()
	defer db.Db.Close()

	if *testFlag {
		tests.TestServer("http://localhost:8080")
	} else {
		if *userFlag {
			// Start CLI to create new user account
			addUser()
		} else {
			// Start server
			InitMiddleware(templates)
			server(port, templates)
		}
	}
}

func server(port string, templates string) {
	log.Println("--- Started Copperhead OTA Server on port", port, "---")

	// Create router
	r := mux.NewRouter()
	admin := mux.NewRouter()

	// Public routes
	r.HandleFunc("/", ReleasesJSON).Methods("GET")
	r.HandleFunc("/", PostReleasesJSON).Methods("POST")
	r.HandleFunc("/releases.json", ReleasesJSON).Methods("GET")
	r.HandleFunc("/changelog/{incremental}.txt", ChangelogFiles).Methods("GET")
	r.HandleFunc("/builds/{name}", DownloadFiles).Methods("GET")
	r.HandleFunc("/v1/build/get_delta", GetDeltaReleases)
	//r.PathPrefix("/static").Handler(http.FileServer(http.Dir("/var/lib/static/")))

	// Authentication
	r.HandleFunc("/login", Login)
	r.HandleFunc("/logout", Logout)
	r.HandleFunc("/authenticate", Authenticate)

	// Releases
	admin.HandleFunc("/admin/releases", Releases)
	admin.HandleFunc("/admin/releases/edit/{id}", EditReleases)
	admin.HandleFunc("/admin/releases/update", UpdateReleases)
	admin.HandleFunc("/admin/releases/new", NewReleases)
	admin.HandleFunc("/admin/releases/create", CreateReleases)
	admin.HandleFunc("/admin/releases/delete", DeleteReleases)
	admin.HandleFunc("/admin/releases/changelog/{id}", ChangelogReleases)

	// Files
	admin.HandleFunc("/admin/files", Files)
	admin.HandleFunc("/admin/files/show/{id}", ShowFiles)
	admin.HandleFunc("/admin/files/delete", DeleteFiles)
	admin.HandleFunc("/admin/files/refresh", RefreshFiles)

	// Users
	admin.HandleFunc("/admin/users", Users)

	// Negroni
	secureMiddleware := secure.New(secure.Options{
		AllowedHosts:          []string{"127.0.0.1", "localhost", "builds.copperhead.co"},
		SSLRedirect:           false,
		STSPreload:            false,
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "default-src 'self'; style-src 'self' https://maxcdn.bootstrapcdn.com; font-src https://maxcdn.bootstrapcdn.com;",
		//IsDevelopment: development,
		IsDevelopment: true,
	})

	// Create a new negroni for the admin middleware
	r.PathPrefix("/admin").Handler(negroni.New(
		AuthMiddleware(),
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
	password, _ := reader.ReadString('\n')
	log.Println("\nSaving...\n")

	// Create user from creds
	database.CreateUser(strings.TrimSpace(username), strings.TrimSpace(password))

	// Exit
	log.Println("Done. Welcome", username)
	log.Println("Exiting")
}
