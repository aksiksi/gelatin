package main

import (
	"flag"
	"io"
	"log"
	"net/http"

	"github.com/aksiksi/gelatin/api"
	"github.com/aksiksi/gelatin/api/emby"
	"github.com/aksiksi/gelatin/api/jellyfin"
)

var (
	jellyfinApiKey string
	embyApiKey     string
)
var ()

func verifyJellyfin() {
	httpClient := &http.Client{}
	client := jellyfin.NewJellyfinApiClient("http://192.168.0.99:8097", httpClient)
	apiKey := api.NewApiKey(jellyfinApiKey)

	if err := client.SystemPing(); err != nil {
		log.Panicf("failed to ping: %s", err)
	}

	systemInfo, err := client.SystemInfoPublic()
	if err != nil {
		log.Panicf("failed to get system info: %s", err)
	}

	log.Printf("Jellyfin info: %+v", systemInfo)

	logsInfo, err := client.SystemLogs(apiKey)
	if err != nil {
		log.Panicf("failed to get system logs: %s", err)
	}

	log.Printf("Jellyfin logs: %+v", logsInfo)

	logName, logSize := logsInfo[0].Name, logsInfo[0].Size
	data, err := client.SystemLogsName(apiKey, logName)
	if err != nil {
		log.Panicf("failed to get system log %s: %s", logName, err)
	}

	logData, _ := io.ReadAll(data)

	log.Printf("got log %s, size = %d, expected = %d", logName, len(logData), logSize)
}

func verifyEmby() {
	httpClient := &http.Client{}
	client := emby.NewEmbyApiClient("http://192.168.0.99:8096/emby", httpClient)
	apiKey := api.NewApiKey(embyApiKey)

	if err := client.SystemPing(); err != nil {
		log.Panicf("failed to ping: %s", err)
	}

	resp, err := client.SystemInfoPublic()
	if err != nil {
		log.Panicf("failed to get system info: %s", err)
	}

	log.Printf("Emby info: %+v", resp)

	logsInfo, err := client.SystemLogsQuery(apiKey)
	if err != nil {
		log.Panicf("failed to get system logs: %s", err)
	}

	log.Printf("Emby logs: %+v", logsInfo)

	logName, logSize := logsInfo.Items[0].Name, logsInfo.Items[0].Size
	data, err := client.SystemLogs(apiKey, logName)
	if err != nil {
		log.Panicf("failed to get system log %s: %s", logName, err)
	}

	logData, _ := io.ReadAll(data)

	log.Printf("got log %s, size = %d, expected = %d", logName, len(logData), logSize)
}

func main() {
	flag.StringVar(&jellyfinApiKey, "jellyfin-api-key", "", "Jellyfin API key")
	flag.StringVar(&embyApiKey, "emby-api-key", "", "Emby API key")
	flag.Parse()

	if jellyfinApiKey == "" {
		log.Fatal("Jellyfin API key must be specified")
	}
	if embyApiKey == "" {
		log.Fatal("Emby API key must be specified")
	}

	verifyEmby()
	verifyJellyfin()
}
