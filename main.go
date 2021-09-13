package main

import (
	"flag"
	"io"
	"log"

	"github.com/aksiksi/gelatin/emby"
	"github.com/aksiksi/gelatin/jellyfin"
)

var (
	jellyfinAdminUser string
	jellyfinAdminPass string
	embyAdminUser     string
	embyAdminPass     string
)

func verifyJellyfin() {
	client := jellyfin.NewJellyfinApiClient("http://192.168.0.99:8097")

	if err := client.System().Ping(); err != nil {
		log.Panicf("failed to ping: %s", err)
	}

	adminKey, err := client.User().Authenticate(jellyfinAdminUser, jellyfinAdminPass)
	if err != nil {
		log.Panicf("failed to authenticate")
	}

	systemInfo, err := client.System().Info(nil, true)
	if err != nil {
		log.Panicf("failed to get system info: %s", err)
	}

	log.Printf("Jellyfin info: %+v", systemInfo)

	logsInfo, err := client.System().GetLogs(adminKey)
	if err != nil {
		log.Panicf("failed to get system logs: %s", err)
	}

	log.Printf("Jellyfin logs: %+v", logsInfo)

	logName, logSize := logsInfo[0].Name, logsInfo[0].Size
	data, err := client.System().GetLogFile(adminKey, logName)
	if err != nil {
		log.Panicf("failed to get system log %s: %s", logName, err)
	}

	logData, _ := io.ReadAll(data)

	log.Printf("got log %s, size = %d, expected = %d", logName, len(logData), logSize)

	// Query public users
	_, err = client.User().GetUsers(nil, true)
	if err != nil {
		log.Panicf("failed to query users: %s", err)
	}

	// Query available users
	users, err := client.User().GetUsers(adminKey, false)
	if err != nil {
		log.Panicf("failed to query users: %s", err)
	}

	log.Printf("Users count: %d", len(users))

	// Create a new user
	user, err := client.User().NewUser(adminKey, "test123")
	if err != nil {
		log.Panicf("failed to create new user: %s", err)
	}

	log.Printf("User: %v", user)

	// Set user password
	err = client.User().UpdatePassword(adminKey, user.Id, "", "abcd1234", false)
	log.Printf("%v", err)

	// Make the user an admin
	user.Policy.IsAdministrator = true
	err = client.User().UpdatePolicy(adminKey, user.Id, &user.Policy)
	if err != nil {
		log.Panic(err)
	}

	// Delete the user
	err = client.User().DeleteUser(adminKey, user.Id)
	if err != nil {
		log.Panic(err)
	}
}

func verifyEmby() {
	client := emby.NewEmbyApiClient("http://192.168.0.99:8096")

	if err := client.System().Ping(); err != nil {
		log.Panicf("failed to ping: %s", err)
	}

	adminKey, err := client.User().Authenticate(embyAdminUser, embyAdminPass)
	if err != nil {
		log.Panicf("failed to authenticate")
	}

	resp, err := client.System().Info(nil, true)
	if err != nil {
		log.Panicf("failed to get system info: %s", err)
	}

	log.Printf("Emby info: %+v", resp)

	logsInfo, err := client.System().GetLogs(adminKey)
	if err != nil {
		log.Panicf("failed to get system logs: %s", err)
	}

	log.Printf("Emby logs: %+v", logsInfo)

	logName, logSize := logsInfo[0].Name, logsInfo[0].Size
	data, err := client.System().GetLogFile(adminKey, logName)
	if err != nil {
		log.Panicf("failed to get system log %s: %s", logName, err)
	}

	logData, _ := io.ReadAll(data)

	log.Printf("got log %s, size = %d, expected = %d", logName, len(logData), logSize)

	// Query public users
	_, err = client.User().GetUsers(nil, true)
	if err != nil {
		log.Panicf("failed to query users: %s", err)
	}

	// Query available users
	users, err := client.User().GetUsers(adminKey, false)
	if err != nil {
		log.Panicf("failed to query users: %s", err)
	}

	log.Printf("Users count: %d", len(users))

	// Create a new user
	user, err := client.User().NewUser(adminKey, "test123")
	if err != nil {
		log.Panicf("failed to create new user: %s", err)
	}

	log.Printf("User: %v", user)

	// Set user password
	err = client.User().UpdatePassword(adminKey, user.Id, "", "abcd1234", false)
	log.Printf("%v", err)

	// Make the user an admin
	user.Policy.IsAdministrator = true
	err = client.User().UpdatePolicy(adminKey, user.Id, &user.Policy)
	if err != nil {
		log.Panic(err)
	}

	// Delete the user
	err = client.User().DeleteUser(adminKey, user.Id)
	if err != nil {
		log.Panic(err)
	}
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

	verifyEmby()
	verifyJellyfin()
}
