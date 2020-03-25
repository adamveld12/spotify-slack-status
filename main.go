package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/nlopes/slack"
)

var (
	slackToken = flag.String("token", "", "slack token")
	playing    = flag.String("playing-emoji", "headphones", "the emoji to use when playing")
	paused     = flag.String("paused-emoji", "musical_note", "the emoji to use when paused")
)

func main() {
	flag.Parse()
	if *slackToken == "" {
		log.Fatal("-token must be defined")
	}

	sc, err := newSpotify()
	if err != nil {
		log.Fatalf("Failed to connect to session bus: %v", err)
	}

	rootCtx := context.Background()

	client := slack.New(*slackToken, slack.OptionDebug(true))

	timer := time.NewTimer(time.Millisecond)
	for {
		select {
		case <-rootCtx.Done():
			break
		case <-timer.C:
			info, err := sc.GetPlaybackInfo()
			if err != nil {
				log.Fatal(err)
			}

			statusEmoji := fmt.Sprintf(":%s:", *paused)
			if info.IsPlaying {
				statusEmoji = fmt.Sprintf(":%s:", *playing)
			}

			var album string
			if info.Album != "" {
				album = fmt.Sprintf(" on \"%s\"", info.Album)
			}

			status := fmt.Sprintf("\"%s\" by \"%s\"%s", info.Title, info.Artist, album)
			if err := client.SetUserCustomStatusContext(rootCtx, status, statusEmoji, info.Duration.Milliseconds()); err != nil {
				log.Printf("could not set user status to '%s': %+v", status, err)
			}

			timer.Reset(time.Second * 10)
		}
	}
}
