package main

import (
	"os"
	"spotify-monthly/internal/auth"
	scheduler "spotify-monthly/internal/cron"
	"spotify-monthly/internal/http"
	"spotify-monthly/internal/playlist"

	"github.com/zmb3/spotify/v2"
)

func main() {

	var url = os.Getenv("BASEURL")

	auth.Setup(url+":1000/callback", "the1")

	clientChannel := make(chan *spotify.Client)

	http.ConfigureServer()

	go auth.GetClient(clientChannel)

	client := <-clientChannel
	playlist.ClientChannel <- client

	internalCron := os.Getenv("internalCron")
	if internalCron == "true" {
		// Schedule playlist creation
		scheduler.SchedulePlaylistCreation()
	}

	select {}
}
