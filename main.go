package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/aksiksi/gelatin/api/emby"
	"github.com/aksiksi/gelatin/api/jellyfin"
)

var (
	jellyfinAdminUser string
	jellyfinAdminPass string
	embyAdminUser     string
	embyAdminPass     string
)

func verifyJellyfin() {
	httpClient := &http.Client{Timeout: 10 * time.Second}
	client := jellyfin.NewJellyfinApiClient("http://192.168.0.99:8097", httpClient)

	if err := client.SystemPing(); err != nil {
		log.Panicf("failed to ping: %s", err)
	}

	adminKey, err := client.UserAuth(jellyfinAdminUser, jellyfinAdminPass)
	if err != nil {
		log.Panicf("failed to authenticate")
	}

	systemInfo, err := client.SystemInfoPublic()
	if err != nil {
		log.Panicf("failed to get system info: %s", err)
	}

	log.Printf("Jellyfin info: %+v", systemInfo)

	logsInfo, err := client.SystemLogs(adminKey)
	if err != nil {
		log.Panicf("failed to get system logs: %s", err)
	}

	log.Printf("Jellyfin logs: %+v", logsInfo)

	logName, logSize := logsInfo[0].Name, logsInfo[0].Size
	data, err := client.SystemLogsName(adminKey, logName)
	if err != nil {
		log.Panicf("failed to get system log %s: %s", logName, err)
	}

	logData, _ := io.ReadAll(data)

	log.Printf("got log %s, size = %d, expected = %d", logName, len(logData), logSize)

	// Query available users
	users, err := client.UserQuery(adminKey)
	if err != nil {
		log.Panicf("failed to query users: %s", err)
	}

	log.Printf("Users count: %d", len(users))

	// Create a new user
	user, err := client.UserNew(adminKey, "test123")
	if err != nil {
		log.Panicf("failed to create new user: %s", err)
	}

	log.Printf("User: %v", user)

	// Set user password
	err = client.UserPassword(adminKey, user.Id, "", "abcd1234", false)
	log.Printf("%v", err)

	// Make the user an admin
	user, err = client.UserGet(adminKey, user.Id)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("%+v", user.Policy)

	user.Policy.IsAdministrator = true

	err = client.UserPolicy(adminKey, user.Id, &user.Policy)
	if err != nil {
		log.Panic(err)
	}

	err = client.UserDelete(adminKey, user.Id)
	if err != nil {
		log.Panic(err)
	}
}

func verifyEmby() {
	httpClient := &http.Client{Timeout: 10 * time.Second}
	client := emby.NewEmbyApiClient("http://192.168.0.99:8096", httpClient)

	if err := client.SystemPing(); err != nil {
		log.Panicf("failed to ping: %s", err)
	}

	adminKey, err := client.UserAuth(embyAdminUser, embyAdminPass)
	if err != nil {
		log.Panicf("failed to authenticate")
	}

	resp, err := client.SystemInfoPublic()
	if err != nil {
		log.Panicf("failed to get system info: %s", err)
	}

	log.Printf("Emby info: %+v", resp)

	logsInfo, err := client.SystemLogsQuery(adminKey)
	if err != nil {
		log.Panicf("failed to get system logs: %s", err)
	}

	log.Printf("Emby logs: %+v", logsInfo)

	logName, logSize := logsInfo.Items[0].Name, logsInfo.Items[0].Size
	data, err := client.SystemLogs(adminKey, logName)
	if err != nil {
		log.Panicf("failed to get system log %s: %s", logName, err)
	}

	logData, _ := io.ReadAll(data)

	log.Printf("got log %s, size = %d, expected = %d", logName, len(logData), logSize)

	// Query available users
	users, err := client.UserQuery(adminKey)
	if err != nil {
		log.Panicf("failed to query users: %s", err)
	}

	log.Printf("Users count: %d", len(users.Items))

	// // Create a new user
	// user, err := client.UserNew(adminKey, "test123")
	// if err != nil {
	// 	log.Panicf("failed to create new user: %s", err)
	// }

	// log.Printf("User: %v", user)

	// // Set user password
	// err = client.UserPassword(adminKey, user.Id, "", "abcd1234", true)
	// log.Printf("%v", err)
}

func main() {
	flag.StringVar(&jellyfinAdminUser, "jellyfin-admin-user", "", "Jellyfin admin username")
	flag.StringVar(&jellyfinAdminPass, "jellyfin-admin-pass", "", "Jellyfin admin password")
	flag.StringVar(&embyAdminUser, "emby-admin-user", "", "Emby admin username")
	flag.StringVar(&embyAdminPass, "emby-admin-pass", "", "Emby admin password")
	flag.Parse()

	if jellyfinAdminUser == "" || jellyfinAdminPass == "" {
		log.Fatal("Jellyfin admin info must be specified")
	}

	if embyAdminUser == "" || embyAdminPass == "" {
		log.Fatal("Emby admin info must be specified")
	}

	// verifyEmby()
	verifyJellyfin()
}
