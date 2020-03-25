package main

import (
	"fmt"
	"io"
	"strings"
	"time"

	dbus "github.com/guelfey/go.dbus"
)

type SpotifyPlaybackInfo struct {
	Time      time.Duration
	Duration  time.Duration
	Artist    string
	Title     string
	Album     string
	IsPlaying bool
}

type SpotifyClient interface {
	io.Closer
	GetPlaybackInfo() (SpotifyPlaybackInfo, error)
}

func newSpotify() (SpotifyClient, error) {
	conn, err := dbus.SessionBus()
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to session bus: %w", err)
	}

	busObject := conn.Object("org.mpris.MediaPlayer2.spotify", "/org/mpris/MediaPlayer2")

	return &spotifyDBusInterface{
		Conn:      conn,
		BusObject: busObject,
	}, nil
}

type spotifyDBusInterface struct {
	Conn      *dbus.Conn
	BusObject *dbus.Object
}

func (sdi *spotifyDBusInterface) GetPlaybackInfo() (SpotifyPlaybackInfo, error) {
	var spi SpotifyPlaybackInfo

	positionProp, err := sdi.BusObject.GetProperty("org.mpris.MediaPlayer2.Player.Position")
	if err != nil {
		return spi, fmt.Errorf("could not retrieve playback position status: %w", err)
	}

	if value, ok := positionProp.Value().(int64); ok {
		spi.Time = time.Duration(value) / time.Nanosecond
	}

	playbackStatusProp, err := sdi.BusObject.GetProperty("org.mpris.MediaPlayer2.Player.PlaybackStatus")
	if err != nil {
		return spi, fmt.Errorf("could not retrieve playback status: %w", err)
	}

	if value, ok := playbackStatusProp.Value().(string); ok && value == "Playing" {
		spi.IsPlaying = true
	}

	metadataProp, err := sdi.BusObject.GetProperty("org.mpris.MediaPlayer2.Player.Metadata")
	if err != nil {
		return spi, fmt.Errorf("could not retrieve metadata: %w", err)
	}

	if metadata, ok := metadataProp.Value().(map[string]dbus.Variant); ok {
		if v := getProperty("xesam:title", metadata); v != nil {
			spi.Title = v.(string)
		}

		if v := getProperty("xesam:album", metadata); v != nil {

			spi.Album = v.(string)
		}

		if v := getProperty("xesam:albumArtist", metadata); v != nil {
			spi.Artist = strings.Join(v.([]string), ", ")
		}

		if spi.Album == spi.Title {
			spi.Album = ""
		}

		if v := getProperty("mpris:length", metadata); v != nil {
			spi.Duration = time.Duration(v.(uint64)) / time.Nanosecond
		}
	}

	return spi, nil
}

func (sdi *spotifyDBusInterface) Close() error {
	return sdi.Conn.Close()
}

func getProperty(key string, metadata map[string]dbus.Variant) interface{} {
	if v, ok := metadata[key]; ok {
		return v.Value()
	}

	return nil
}
