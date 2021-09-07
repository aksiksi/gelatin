package main

import (
	"log"
	"net/http"

	"github.com/aksiksi/gelatin/api/emby"
	"github.com/aksiksi/gelatin/api/jellyfin"
)

func verifyJellyfin() {
	httpClient := &http.Client{}
	client := jellyfin.NewJellyfinApiClient("http://192.168.0.99:8097", httpClient)

	if err := client.SystemPing(); err != nil {
		log.Panicf("failed to ping: %s", err)
	}

	resp, err := client.SystemInfoPublic()
	if err != nil {
		log.Panicf("failed to get system info: %s", err)
	}

	log.Printf("Jellyfin info: %+v", resp)
}

func verifyEmby() {
	httpClient := &http.Client{}
	client := emby.NewEmbyApiClient("http://192.168.0.99:8096/emby", httpClient)

	if err := client.SystemPing(); err != nil {
		log.Panicf("failed to ping: %s", err)
	}

	resp, err := client.SystemInfoPublic()
	if err != nil {
		log.Panicf("failed to get system info: %s", err)
	}

	log.Printf("Emby info: %+v", resp)
}

func main() {
	verifyEmby()
	verifyJellyfin()
}
