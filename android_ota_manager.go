// Copperhead OTA Server

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/codegangsta/negroni"
	"github.com/copperhead/android_ota_manager/controllers"
	"github.com/copperhead/android_ota_manager/models"
	"github.com/copperhead/android_ota_manager/tests"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"
	"github.com/unrolled/secure"
	"gopkg.in/gorp.v1"
)

var (
	db *gorp.DbMap
)

func dynamicPublicCaching(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	w.Header().Set("Cache-Control", "public, no-cache")
	next(w, r)
}

func main() {
	userFlag := flag.Bool("add_user", false, "Run CLI for adding user to database")
	testFlag := flag.Bool("test", false, "Run test script to simulate client")
	flag.Parse()

	// Connect to database
	data := os.Getenv("OPENSHIFT_DATA_DIR")
	db := models.InitDb(data+"ota.sql", data+"builds")
	defer db.Db.Close()

	go models.RefreshBuilds()

	if *testFlag {
		tests.TestServer("http://localhost:8080")
	} else {
		if *userFlag {
			// Start CLI to create new user account
			addUser()
		} else {
			// Start server
			templates := "./views"
			controllers.InitMiddleware(templates)
			server(templates)
		}
	}
}

func server(templates string) {
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

	// Authentication
	r.HandleFunc("/login", controllers.Login)
	r.HandleFunc("/logout", controllers.Logout)
	r.HandleFunc("/authenticate", controllers.Authenticate)

	// Static content
	data := os.Getenv("OPENSHIFT_DATA_DIR")
	challenge := http.FileServer(http.Dir(data + "acme-challenge"))
	r.PathPrefix("/.well-known/acme-challenge").Handler(http.StripPrefix("/.well-known/acme-challenge", challenge))

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
	admin.HandleFunc("/admin/users/delete", controllers.DeleteUsers)

	// Negroni
	secureMiddleware := secure.New(secure.Options{
		AllowedHosts:          []string{"127.0.0.1", "localhost", "builds.copperhead.co", "builds-copperheadsec.rhcloud.com"},
		SSLRedirect:           true,
		SSLProxyHeaders:       map[string]string{"X-Forwarded-Proto": "https"},
		STSSeconds:            15552000,
		STSIncludeSubdomains:  true,
		STSPreload:            true,
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "default-src 'self'; style-src 'self' https://maxcdn.bootstrapcdn.com; font-src https://maxcdn.bootstrapcdn.com;",
		IsDevelopment:         false,
	})

	// Create a new negroni for the admin middleware
	r.PathPrefix("/admin").Handler(negroni.New(
		controllers.AuthMiddleware(),
		negroni.Wrap(admin),
	))

	n := negroni.Classic()
	n.Use(negroni.HandlerFunc(dynamicPublicCaching))
	n.Use(negroni.HandlerFunc(secureMiddleware.HandlerFuncWithNext))
	n.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET"},
		AllowedHeaders:   []string{"Origin"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: false,
	}))
	n.UseHandler(r)
	bind := fmt.Sprintf("%s:%s", os.Getenv("OPENSHIFT_GO_IP"), os.Getenv("OPENSHIFT_GO_PORT"))
	n.Run(bind)
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

	// Create user from creds
	models.CreateUser(strings.TrimSpace(username), strings.TrimSpace(password))

	// Exit
	log.Println("Done. Welcome", username)
}
