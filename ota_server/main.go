// Copperhead OTA Server

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/copperhead-security/android_ota_manager/database"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/gorp.v1"
)

var (
	db  *gorp.DbMap
	env string
	host string
)

func main() {
	// Parse CLI arguments
	port := flag.String("port", "8080", "Server port")
	host = *flag.String("host", "http://localhost", "Server host")
	env = *flag.String("env", "development", "Run in development or production mode")
	dbPath := flag.String("db", "./ota.sql", "Path to sqlite db file")
	staticDir := flag.String("static", "./admin_interface/static/", "Path to templates")
	buildsDir := flag.String("builds", "./builds", "Path to directory containing build images")
	userFlag := flag.Bool("add_user", false, "Run CLI for adding user to database")
	testFlag := flag.Bool("test", false, "Run test script to simulate client")
	flag.Parse()

	// Connect to database
	db := database.InitDb(*dbPath, *buildsDir)
	go database.RefreshBuilds()
	defer db.Db.Close()

	if *testFlag {
		// TestServer("http://localhost:8080")
	} else {
		if *userFlag {
			// Start CLI to create new user account
			addUser()
		} else {
			// Start server
			server(*port, *staticDir)
		}
	}
}

func isDevelopment() bool {
	return strings.Contains(env, "development")
}

func cat(s ...string) string {
	return strings.Join(s, "")
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

func CheckErr(err error, msg string) {
	if err != nil {
		log.Println(msg, err)
	}
}
