package scheduler

import (
	"fmt"
	"log"
	"time"

	"spotify-monthly/internal/playlist"

	"github.com/robfig/cron/v3"
)

func SchedulePlaylistCreation() {

	c := cron.New()

	_, err := c.AddFunc("0 0 1 * *", func() {
		fmt.Printf("Creating playlist at %s", time.Now().Format(time.RFC1123))
		playlist.CreatePlaylist()
	})

	if err != nil {
		log.Fatalf("Error scheduling monthly job: %v", err)
	}

	c.Start()
}
