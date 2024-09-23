package main

import (
	"os"
	"spotify-monthly/internal/auth"
	"spotify-monthly/internal/http"

	"github.com/zmb3/spotify/v2"
)

func main() {

	auth.Setup(os.Getenv("BASEURL")+"/callback", "the1")

	clientChannel := make(chan *spotify.Client)

	http.ConfigureServer()

	go auth.GetClient(clientChannel)

	select {}
}
