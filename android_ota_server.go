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
	"github.com/copperhead/android_ota_manager/controllers"
	"github.com/copperhead/android_ota_manager/lib"
	"github.com/copperhead/android_ota_manager/models"
	"github.com/copperhead/android_ota_manager/tests"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/olebedev/config"
	"github.com/unrolled/secure"
	"github.com/yageek/cors"
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
	db := models.InitDb(databasePath, builds)
	go models.RefreshBuilds()
	defer db.Db.Close()

	if *testFlag {
		tests.TestServer("http://localhost:8080")
	} else {
		if *userFlag {
			// Start CLI to create new user account
			addUser()
		} else {
			// Start server
			controllers.InitMiddleware(templates)
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
	r.HandleFunc("/", controllers.ReleasesJSON).Methods("GET")
	r.HandleFunc("/", controllers.PostReleasesJSON).Methods("POST")
	r.HandleFunc("/releases.json", controllers.ReleasesJSON).Methods("GET")
	r.HandleFunc("/changelog/{incremental}.txt", controllers.ChangelogFiles).Methods("GET")
	r.HandleFunc("/builds/{name}", controllers.DownloadFiles).Methods("GET")
	r.HandleFunc("/v1/build/get_delta", controllers.GetDeltaReleases)
	//r.PathPrefix("/static").Handler(http.FileServer(http.Dir("/var/lib/static/")))

	// Authentication
	r.HandleFunc("/login", controllers.Login)
	r.HandleFunc("/logout", controllers.Logout)
	r.HandleFunc("/authenticate", controllers.Authenticate)

	// Releases
	admin.HandleFunc("/admin/releases", controllers.Releases)
	admin.HandleFunc("/admin/releases/edit/{id}", controllers.EditReleases)
	admin.HandleFunc("/admin/releases/update", controllers.UpdateReleases)
	admin.HandleFunc("/admin/releases/new", controllers.NewReleases)
	admin.HandleFunc("/admin/releases/create", controllers.CreateReleases)
	admin.HandleFunc("/admin/releases/delete", controllers.DeleteReleases)
	admin.HandleFunc("/admin/releases/changelog/{id}", controllers.ChangelogReleases)

	// Files
	admin.HandleFunc("/admin/files", controllers.Files)
	admin.HandleFunc("/admin/files/show/{id}", controllers.ShowFiles)
	admin.HandleFunc("/admin/files/delete", controllers.DeleteFiles)
	admin.HandleFunc("/admin/files/refresh", controllers.RefreshFiles)

	// Users
	admin.HandleFunc("/admin/users", controllers.Users)

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
		controllers.AuthMiddleware(),
		negroni.Wrap(admin),
	))

	n := negroni.Classic()
	n.Use(negroni.HandlerFunc(secureMiddleware.HandlerFuncWithNext))
	n.Use(cors.Allow(&cors.Options{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
	}))
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
	models.CreateUser(strings.TrimSpace(username), strings.TrimSpace(password))

	// Exit
	log.Println("Done. Welcome", username)
	log.Println("Exiting")
}
